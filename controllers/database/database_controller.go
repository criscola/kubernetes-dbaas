/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bedag/kubernetes-dbaas/internal/logging"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	. "github.com/bedag/kubernetes-dbaas/pkg/typeutil"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/client-go/tools/record"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"

	databasev1 "github.com/bedag/kubernetes-dbaas/apis/database/v1"
	databaseclassv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
)

const (
	InfoLevel  = logging.InfoLevel
	DebugLevel = logging.ZapDebugLevel
	TraceLevel = logging.ZapTraceLevel

	DatabaseControllerName = "database-controller"
	DatabaseClass          = "databaseclass"
	EndpointName           = "endpoint-name"
	SecretName             = "secret-name"
	databaseFinalizer      = "finalizer.database.bedag.ch"
	rotateAnnotationKey    = "dbaas.bedag.ch/rotate"
)

type ReconcileError struct {
	Reason         string
	Message        string
	Err            error
	AdditionalInfo []interface{}
}

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
	DbmsList      database.DbmsList
	Pool          pool.Pool
}

var logger logr.Logger

// +kubebuilder:rbac:groups=database.dbaas.bedag.ch,resources=databases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.dbaas.bedag.ch,resources=databases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=database.dbaas.bedag.ch,resources=databases/finalizers,verbs=update
// +kubebuilder:rbac:groups=databaseclass.dbaas.bedag.ch,resources=databaseclasses,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=list;watch;create;update;delete
// SetupWithManager creates the controller responsible for Database resources by means of a ctrl.Manager.
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named(DatabaseControllerName).
		For(&databasev1.Database{}).
		Owns(&corev1.Secret{}).
		WithEventFilter(r.triggerReconciler()).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = r.Log.WithValues("database", req.NamespacedName)
	logger.V(TraceLevel).Info("Reconcile called")

	obj := &databasev1.Database{}
	err := r.Get(ctx, req.NamespacedName, obj)
	if err != nil {
		if k8sError.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.V(TraceLevel).Info(MsgDbDeleted)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		r.handleReconcileError(obj, ReconcileError{
			Reason:  RsnDbGetFail,
			Message: MsgDbGetFail,
			Err:     err,
		})
		return reconcile.Result{Requeue: true}, nil
	}

	// Set reason to unknown to indicate the resource was correctly received by a controller but no action was resolved yet
	// Update condition field
	if meta.FindStatusCondition(obj.Status.Conditions, TypeReady) == nil {
		logger.V(TraceLevel).Info("Updating ConditionStatus")
		if err = r.updateReadyCondition(obj, metav1.ConditionUnknown, RsnDbOpQueueSucc, MsgDbOpQueueSucc); err != nil {
			r.handleReconcileError(obj, ReconcileError{
				Reason:  RsnReadyCondUpdateFail,
				Message: MsgReadyCondUpdateFail,
				Err:     err,
			})

			return ctrl.Result{Requeue: true}, nil
		}
	}

	// Check if the Database instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if obj.GetDeletionTimestamp() != nil {
		if contains(obj.GetFinalizers(), databaseFinalizer) {
			// Run finalization logic for DatabaseFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			logger.V(TraceLevel).Info("Finalizing database resource")
			if err := r.deleteDb(obj); err.IsNotEmpty() {
				r.handleReconcileError(obj, err)
				return reconcile.Result{Requeue: true}, nil
			}

			// Remove databaseFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			logger.V(TraceLevel).Info("Removing finalizer")
			controllerutil.RemoveFinalizer(obj, databaseFinalizer)
			if err := r.Update(ctx, obj); err != nil {
				if !shouldIgnoreUpdateErr(err) {
					r.handleReconcileError(obj, ReconcileError{
						Reason:         RsnDbUpdateFail,
						Message:        MsgDbUpdateFail,
						Err:            err,
						AdditionalInfo: StringsToInterfaceSlice("finalizer", databaseFinalizer),
					})
				}

				return reconcile.Result{Requeue: true}, nil
			}
		}
		return reconcile.Result{}, nil
	}

	// If Database is ready
	if meta.IsStatusConditionTrue(obj.Status.Conditions, TypeReady) {
		logger.V(TraceLevel).Info("Database resource is in Ready state")
		// Check if Database credentials should be rotated
		shouldRotate, err := r.shouldRotate(obj)
		if err.IsNotEmpty() {
			r.handleReconcileError(obj, err)
			return ctrl.Result{Requeue: true}, nil
		}
		if shouldRotate {
			// Update Ready condition to false, Database credentials must be rotated
			if err := r.updateReadyCondition(obj, metav1.ConditionFalse, RsnDbRotateInProg, MsgDbRotateInProg); err != nil {
				r.handleReadyConditionError(obj, err)
				return ctrl.Result{Requeue: true}, nil
			}
			if err := r.rotate(obj); err.IsNotEmpty() {
				r.handleReconcileError(obj, err)
				return ctrl.Result{Requeue: true}, nil
			}
			// Update Ready condition to true
			if err := r.updateReadyCondition(obj, metav1.ConditionTrue, RsnDbRotateSucc, MsgDbRotateSucc); err != nil {
				r.handleReadyConditionError(obj, err)
				return ctrl.Result{Requeue: true}, nil
			}
			r.logInfoEvent(obj, RsnDbRotateSucc, MsgDbRotateSucc)
		} else {
			// Database is ready and credentials shouldn't be rotated, nothing else to do
			logger.V(TraceLevel).Info("Credentials should not be rotated, nothing left to do")
			return ctrl.Result{}, nil
		}
	} else {
		// Create
		if err := r.createDb(obj); err.IsNotEmpty() {
			r.handleReconcileError(obj, err)
			return ctrl.Result{Requeue: true}, nil
		}

		logger.V(TraceLevel).Info("Updating ConditionStatus")
		if err := r.updateReadyCondition(obj, metav1.ConditionTrue, RsnDbCreateSucc, MsgDbCreateSucc); err != nil {
			r.handleReadyConditionError(obj, err)
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// If finalizer is not present, add finalizer to resource
	if !contains(obj.GetFinalizers(), databaseFinalizer) {
		logger.V(TraceLevel).Info("Adding finalizer")
		if err := r.addFinalizer(obj); err != nil {
			r.handleReconcileError(obj, ReconcileError{
				Reason:         RsnDbUpdateFail,
				Message:        MsgDbUpdateFail,
				Err:            err,
				AdditionalInfo: StringsToInterfaceSlice("finalizer", databaseFinalizer),
			})
			return ctrl.Result{Requeue: true}, nil
		}
	}

	logger.V(TraceLevel).Info("Reached end of reconcile")
	return ctrl.Result{}, nil
}

// addFinalizer adds a finalizer to a Database resource.
func (r *DatabaseReconciler) addFinalizer(obj *databasev1.Database) error {
	controllerutil.AddFinalizer(obj, databaseFinalizer)
	return r.Update(context.Background(), obj)
}

// createDb creates a new Database instance on the external provisioner based on the Database data.
func (r *DatabaseReconciler) createDb(obj *databasev1.Database) ReconcileError {
	r.logInfoEvent(obj, RsnDbCreateInProg, MsgDbCreateInProg)

	dbClass, err := r.getDbmsClassFromDb(obj)
	if err.IsNotEmpty() {
		return err
	}
	loggingKv := StringsToInterfaceSlice(DatabaseClass, dbClass.Name, database.OperationsConfigKey, database.CreateMapKey)

	// Render operation
	createOpTemplate, exists := dbClass.Spec.Operations[database.CreateMapKey]
	if !exists {
		return ReconcileError{
			Reason:         RsnOpNotSupported,
			Message:        MsgOpNotSupported,
			Err:            nil,
			AdditionalInfo: loggingKv,
		}
	}
	opValues, err := newOpValuesFromResource(obj)
	if err.IsNotEmpty() {
		return err.With(loggingKv)
	}
	createOp, simpleErr := createOpTemplate.RenderOperation(opValues)
	if simpleErr != nil {
		return ReconcileError{
			Reason:         RsnOpRenderFail,
			Message:        MsgOpRenderFail,
			Err:            simpleErr,
			AdditionalInfo: loggingKv,
		}
	}
	loggingKv = append(loggingKv, EndpointName, obj.Spec.Endpoint)

	// Execute operation on DBMS
	// Check preconditions
	var conn database.Driver
	if conn, err = r.getDbmsConnectionByEndpointName(obj.Spec.Endpoint); err.IsNotEmpty() {
		return err.With(loggingKv)
	}
	output := conn.CreateDb(createOp)
	if output.Err != nil {
		return ReconcileError{
			Reason:         RsnDbCreateFail,
			Message:        MsgDbCreateFail,
			Err:            output.Err,
			AdditionalInfo: loggingKv,
		}
	}

	// Log success
	r.logInfoEvent(obj, RsnDbCreateSucc, MsgDbCreateSucc)
	logger.V(TraceLevel).Info(fmt.Sprint(dbClass.Spec.SecretFormat))
	logger.V(TraceLevel).Info(fmt.Sprint(output))
	// Create Secret
	err = r.createSecret(obj, dbClass.Spec.SecretFormat, output)
	if err.IsNotEmpty() {
		return err.With(loggingKv)
	}

	return ReconcileError{}
}

// deleteDb deletes the database instance on the external provisioner.
func (r *DatabaseReconciler) deleteDb(obj *databasev1.Database) ReconcileError {
	r.logInfoEvent(obj, RsnDbDeleteInProg, MsgDbDeleteInProg)

	dbClass, reconcileErr := r.getDbmsClassFromDb(obj)
	if reconcileErr.IsNotEmpty() {
		return reconcileErr
	}
	loggingKv := StringsToInterfaceSlice(DatabaseClass, dbClass.Name, database.OperationsConfigKey, database.DeleteMapKey)

	// Render operation
	deleteOpTemplate, exists := dbClass.Spec.Operations[database.DeleteMapKey]
	if !exists {
		return ReconcileError{
			Reason:         RsnOpNotSupported,
			Message:        MsgOpNotSupported,
			Err:            nil,
			AdditionalInfo: loggingKv,
		}
	}
	opValues, reconcileErr := newOpValuesFromResource(obj)
	if reconcileErr.IsNotEmpty() {
		return reconcileErr.With(loggingKv)
	}
	deleteOp, err := deleteOpTemplate.RenderOperation(opValues)
	if err != nil {
		return ReconcileError{
			Reason:         RsnOpRenderFail,
			Message:        MsgOpRenderFail,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	loggingKv = append(loggingKv, EndpointName, obj.Spec.Endpoint)

	// Execute operation on DBMS
	// Check preconditions
	var conn database.Driver
	if conn, reconcileErr = r.getDbmsConnectionByEndpointName(obj.Spec.Endpoint); reconcileErr.IsNotEmpty() {
		return reconcileErr.With(loggingKv)
	}
	conn = r.Pool.Get(obj.Spec.Endpoint)
	if conn == nil {
		return ReconcileError{
			Reason:         RsnDbmsEndpointNotFound,
			Message:        MsgDbmsEndpointNotFound,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	output := conn.DeleteDb(deleteOp)
	if output.Err != nil {
		return ReconcileError{
			Reason:         RsnDbDeleteFail,
			Message:        MsgDbDeleteFail,
			Err:            output.Err,
			AdditionalInfo: loggingKv,
		}
	}
	return ReconcileError{}
}

// rotate rotates the database credentials on the external provisioner.
func (r *DatabaseReconciler) rotate(obj *databasev1.Database) ReconcileError {
	r.logInfoEvent(obj, RsnDbRotateInProg, MsgDbRotateInProg)

	dbClass, reconcileErr := r.getDbmsClassFromDb(obj)
	if reconcileErr.IsNotEmpty() {
		return reconcileErr
	}
	loggingKv := StringsToInterfaceSlice(DatabaseClass, dbClass.Name, database.OperationsConfigKey, database.RotateMapKey)
	rotateOpTemplate, exists := dbClass.Spec.Operations[database.RotateMapKey]
	if !exists {
		return ReconcileError{
			Reason:         RsnOpNotSupported,
			Message:        MsgOpNotSupported,
			Err:            nil,
			AdditionalInfo: loggingKv,
		}
	}
	opValues, reconcileErr := newOpValuesFromResource(obj)
	if reconcileErr.IsNotEmpty() {
		return reconcileErr.With(loggingKv)
	}
	rotateOp, err := rotateOpTemplate.RenderOperation(opValues)
	if err != nil {
		return ReconcileError{
			Reason:         RsnOpRenderFail,
			Message:        MsgOpRenderFail,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	loggingKv = append(loggingKv, EndpointName, obj.Spec.Endpoint)

	// Execute operation on DBMS
	// Check preconditions
	var conn database.Driver
	if conn, reconcileErr = r.getDbmsConnectionByEndpointName(obj.Spec.Endpoint); reconcileErr.IsNotEmpty() {
		return reconcileErr.With(loggingKv)
	}
	conn = r.Pool.Get(obj.Spec.Endpoint)
	if conn == nil {
		return ReconcileError{
			Reason:         RsnDbmsEndpointNotFound,
			Message:        MsgDbmsEndpointNotFound,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	output := conn.Rotate(rotateOp)
	if output.Err != nil {
		return ReconcileError{
			Reason:         RsnDbRotateFail,
			Message:        MsgDbRotateFail,
			Err:            output.Err,
			AdditionalInfo: loggingKv,
		}
	}

	if isSecretPresent, err := r.isSecretPresent(obj); isSecretPresent {
		if err.IsNotEmpty() {
			return err.With(loggingKv)
		}
		// Secret is already present, update it
		err = r.updateSecret(obj, dbClass.Spec.SecretFormat, output)
		if err.IsNotEmpty() {
			return err.With(loggingKv)
		}
	} else {
		// Secret is not present, create it
		err = r.createSecret(obj, dbClass.Spec.SecretFormat, output)
		if err.IsNotEmpty() {
			return err.With(loggingKv)
		}
	}

	// Remove annotation if present
	if isRotateAnnotationTrue(obj) {
		logger.V(TraceLevel).Info("Removing rotate annotation")
		delete(obj.GetAnnotations(), rotateAnnotationKey)
		err := r.Client.Update(context.Background(), obj)
		if err != nil {
			return ReconcileError{
				Reason:         RsnDbUpdateFail,
				Message:        MsgDbUpdateFail,
				Err:            err,
				AdditionalInfo: loggingKv,
			}
		}
	}

	return ReconcileError{}
}

func (r *DatabaseReconciler) getDbmsClassFromDb(obj *databasev1.Database) (databaseclassv1.DatabaseClass, ReconcileError) {
	// Get DatabaseClass resource from api server
	dbClassName := r.DbmsList.GetDatabaseClassNameByEndpointName(obj.Spec.Endpoint)
	if dbClassName == "" {
		return databaseclassv1.DatabaseClass{}, ReconcileError{
			Reason:         RsnDbcConfigGetFail,
			Message:        MsgDbcConfigGetFail,
			Err:            fmt.Errorf("could not find any DatabaseClass for endpoint '%s'", obj.Spec.Endpoint),
			AdditionalInfo: StringsToInterfaceSlice(EndpointName, obj.Spec.Endpoint),
		}
	}

	dbClass := databaseclassv1.DatabaseClass{}
	err := r.Client.Get(context.Background(), client.ObjectKey{Namespace: "", Name: dbClassName}, &dbClass)
	if err != nil {
		return databaseclassv1.DatabaseClass{}, ReconcileError{
			Reason:         RsnDbcGetFail,
			Message:        MsgDbcGetFail,
			Err:            err,
			AdditionalInfo: StringsToInterfaceSlice(DatabaseClass, dbClassName),
		}
	}
	return dbClass, ReconcileError{}
}

func (r *DatabaseReconciler) getDbmsConnectionByEndpointName(endpointName string) (database.Driver, ReconcileError) {
	// Check if the endpoint is currently stored in the connection pool
	conn := r.Pool.Get(endpointName)
	if conn == nil {
		return nil, ReconcileError{
			Reason:  RsnDbmsEndpointNotFound,
			Message: MsgDbmsEndpointNotFound,
			Err:     nil,
		}
	}
	// Make a first check to acknowledge whether the connection looks alive
	if simpleErr := conn.Ping(); simpleErr != nil {
		return nil, ReconcileError{
			Reason:  RsnDbmsConnFail,
			Message: MsgDbmsConnFail,
			Err:     simpleErr,
		}
	}
	return conn, ReconcileError{}
}

// handleReconcileError sets the obj Conditions type Ready to false and sets the relative fields error and message,
// it records a Warning event with reason and message for the given obj and logs err (if present) and message to the
// global logger.
// It ignores optimistic locking error, see shouldIgnoreUpdateErr.
func (r *DatabaseReconciler) handleReconcileError(obj *databasev1.Database, err ReconcileError) {
	if shouldIgnoreUpdateErr(err.Err) {
		logger.V(TraceLevel).Info(err.Err.Error())
		return
	}
	keyAndValuesLen := len(err.AdditionalInfo)
	if keyAndValuesLen%2 != 0 {
		logger.Error(nil, "odd number of keyAndValues provided!", err.AdditionalInfo...)
		// Set length to 0 so additional values are ignored
		keyAndValuesLen = 0
	}
	if keyAndValuesLen > 0 {
		r.EventRecorder.Event(obj, Warning, err.Reason, formatEventMessage(err.Message, err.AdditionalInfo...))
		logger.Error(err.Err, err.Message, err.AdditionalInfo...)
	} else {
		r.EventRecorder.Event(obj, Warning, err.Reason, err.Message)
		logger.Error(err.Err, err.Message)
	}
	if updateErr := r.updateReadyCondition(obj, metav1.ConditionFalse, err.Reason, err.Message); updateErr != nil {
		logger.Error(err.Err, MsgDbUpdateFail)
	}
}

// handleReadyConditionError records an event of type Warning to obj using RsnReadyCondUpdateFail, MsgReadyCondUpdateFail
// and additionalInfo. additionalInfo is formatted as JSON and attached to the event message.
// An error log using message and additionalInfo is written using the global logger.
// It ignores optimistic locking error, see shouldIgnoreUpdateErr.
func (r *DatabaseReconciler) handleReadyConditionError(obj *databasev1.Database, err error, additionalInfo ...interface{}) {
	if shouldIgnoreUpdateErr(err) {
		logger.V(TraceLevel).Info(err.Error())
		return
	}
	// In the grim situation where the Ready condition cannot be updated, dump everything to the resource event stream
	// and logger
	eventMessage := formatEventMessage(MsgReadyCondUpdateFail, additionalInfo...)
	r.EventRecorder.Event(obj, Warning, RsnReadyCondUpdateFail, eventMessage)
	logger.Error(err, MsgReadyCondUpdateFail, additionalInfo...)
}

// logInfoEvent records an event of type Normal to obj using reason, message and additionalInfo. additionalInfo is formatted
// as JSON and attached to the event message. An info log using message and additionalInfo is written using the global logger
func (r *DatabaseReconciler) logInfoEvent(obj *databasev1.Database, reason, message string, additionalInfo ...interface{}) {
	eventMessage := formatEventMessage(message, additionalInfo...)
	r.EventRecorder.Event(obj, Normal, reason, eventMessage)
	logger.Info(message, additionalInfo...)
}

// createSecret creates a new K8s secret owned by owner with the data contained in output and dsn.
func (r *DatabaseReconciler) createSecret(owner *databasev1.Database, secretFormat database.SecretFormat, output database.OpOutput) ReconcileError {
	logger.V(DebugLevel).Info("Creating secret for database resource")

	// Init vars
	secretName := FormatSecretName(owner)
	loggingKv := StringsToInterfaceSlice("secret", secretName)
	secretData, err := secretFormat.RenderSecretFormat(output)
	if err != nil {
		return ReconcileError{
			Reason:         RsnSecretRenderFail,
			Message:        MsgSecretRenderFail,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	var ownerRefs []metav1.OwnerReference
	ownerRefs = append(ownerRefs, metav1.OwnerReference{
		APIVersion: owner.APIVersion,
		Kind:       owner.Kind,
		Name:       owner.Name,
		UID:        owner.UID,
		Controller: &[]bool{true}[0], // sets this controller as owner
	})
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            secretName,
			Namespace:       owner.Namespace,
			OwnerReferences: ownerRefs,
		},
		StringData: secretData,
	}
	key := client.ObjectKey{
		Namespace: owner.Namespace,
		Name:      secretName,
	}
	oldSecret := corev1.Secret{}
	// Get old Secret if present
	err = r.Client.Get(context.Background(), key, &oldSecret)
	if err != nil {
		// If Secret was not found, it must be created
		if k8sError.IsNotFound(err) {
			// Create new Secret
			if err := r.Client.Create(context.Background(), secret); err != nil {
				return ReconcileError{
					Reason:         MsgSecretCreateFail,
					Message:        MsgSecretCreateFail,
					Err:            err,
					AdditionalInfo: loggingKv,
				}
			}
			r.logInfoEvent(owner, RsnSecretCreateSucc, MsgSecretCreateSucc, loggingKv...)
			return ReconcileError{}
		}
		// Return error
		return ReconcileError{
			Reason:         RsnSecretGetFail,
			Message:        MsgSecretGetFail,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	} else {
		// Create was called on already existing Secret
		return ReconcileError{
			Reason:         RsnSecretExists,
			Message:        MsgSecretExists,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
}

// createSecret creates a new K8s secret owned by owner with the data contained in output and dsn.
func (r *DatabaseReconciler) updateSecret(owner *databasev1.Database, secretFormat database.SecretFormat, output database.OpOutput) ReconcileError {
	logger.V(DebugLevel).Info("Updating secret for database resource")

	// TODO: extract common behavior of Secret rendering into a method and put it in createSecret as well (factory method)?
	// Init vars
	secretName := FormatSecretName(owner)
	loggingKv := StringsToInterfaceSlice("secret", secretName)
	secretData, err := secretFormat.RenderSecretFormat(output)
	if err != nil {
		return ReconcileError{
			Reason:         RsnSecretRenderFail,
			Message:        MsgSecretRenderFail,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	var ownerRefs []metav1.OwnerReference
	ownerRefs = append(ownerRefs, metav1.OwnerReference{
		APIVersion: owner.APIVersion,
		Kind:       owner.Kind,
		Name:       owner.Name,
		UID:        owner.UID,
		Controller: &[]bool{true}[0], // sets this controller as owner
	})
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            secretName,
			Namespace:       owner.Namespace,
			OwnerReferences: ownerRefs,
		},
		StringData: secretData,
	}
	if err := r.Client.Update(context.Background(), newSecret); err != nil {
		return ReconcileError{
			Reason:         MsgSecretUpdateFail,
			Message:        MsgSecretUpdateFail,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	r.logInfoEvent(owner, RsnSecretUpdateSucc, MsgSecretUpdateSucc, loggingKv...)
	return ReconcileError{}
}

// updateReadyCondition updates the Ready Condition status of obj.
func (r *DatabaseReconciler) updateReadyCondition(obj *databasev1.Database, status metav1.ConditionStatus, reason, message string) error {
	meta.SetStatusCondition(&obj.Status.Conditions, metav1.Condition{
		Type:    TypeReady,
		Status:  status,
		Reason:  reason,
		Message: message,
	})

	// Update condition field
	return r.Client.Status().Update(context.Background(), obj)
}

// shouldRotate returns true if there isn't any Secret associated with the given Database object (secret deletion),
// or if the rotate annotation is present. It returns false otherwise, or if an error was generated during execution.
func (r *DatabaseReconciler) shouldRotate(obj *databasev1.Database) (bool, ReconcileError) {
	logger.V(TraceLevel).Info("Checking if credentials should be rotated")
	if isSecretPresent, err := r.isSecretPresent(obj); !isSecretPresent {
		if err.IsNotEmpty() {
			return false, ReconcileError{
				Reason:  RsnSecretGetFail,
				Message: MsgSecretGetFail,
				Err:     err.Err,
			}
		}
		return true, ReconcileError{}
	}
	// secret is present, check if rotate annotation is present, if yes, rotate, else, just keep going
	if isRotateAnnotationTrue(obj) {
		return true, ReconcileError{}
	}
	return false, ReconcileError{}
}

// isSecretPresent returns true if the Secret bound to obj is present. It returns false otherwise, or if an
// error was generated during execution.
func (r *DatabaseReconciler) isSecretPresent(obj *databasev1.Database) (bool, ReconcileError) {
	secretName := FormatSecretName(obj)
	loggingKv := StringsToInterfaceSlice(SecretName, secretName)
	logger.V(TraceLevel).Info("Checking if secret bound to Database resource is present")

	var secret corev1.Secret
	secretObjKey := client.ObjectKey{Namespace: obj.Namespace, Name: FormatSecretName(obj)}
	if err := r.Client.Get(context.Background(), secretObjKey, &secret); err != nil {
		if k8sError.IsNotFound(err) {
			// Secret for given object is not present
			return false, ReconcileError{}
		}
		// Another error was generated while getting the Secret, return it
		return false, ReconcileError{
			Reason:         RsnSecretGetFail,
			Message:        MsgSecretGetFail,
			Err:            err,
			AdditionalInfo: loggingKv,
		}
	}
	return true, ReconcileError{}
}

// IsNotEmpty checks if r is not empty using reflect.DeepEqual. Needed because field AdditionalInfo is not comparable.
func (r ReconcileError) IsNotEmpty() bool {
	return !reflect.DeepEqual(r, ReconcileError{})
}

// With creates a copy of the receiver and appends values to its AdditionalInfo field.
func (r ReconcileError) With(values []interface{}) ReconcileError {
	return ReconcileError{
		Reason:         r.Reason,
		Message:        r.Message,
		Err:            r.Err,
		AdditionalInfo: append(r.AdditionalInfo, values...),
	}
}

// FormatSecretName returns the name of a Database's Secret resource as it should appear in metadata.name.
func FormatSecretName(obj *databasev1.Database) string {
	return obj.Name + "-credentials"
}

// triggerReconciler checks whether reconciliation should be triggered or not to avoid useless reconciliations.
// See also predicate.Predicate.
func (r *DatabaseReconciler) triggerReconciler() predicate.Predicate {
	return predicate.Funcs{
		GenericFunc: func(e event.GenericEvent) bool {
			obj := e.Object.(*databasev1.Database)
			// If credentials are supposed to be rotated
			if shouldRotate, err := r.shouldRotate(obj); shouldRotate || err.IsNotEmpty() {
				return true
			}
			// If object is supposed to be deleted
			if obj.GetDeletionTimestamp() != nil && contains(e.Object.GetFinalizers(), databaseFinalizer) {
				return true
			}
			// If ready condition is false or unknown
			if !meta.IsStatusConditionTrue(obj.Status.Conditions, TypeReady) {
				return true
			}
			return false
		},
	}
}

// isRotateAnnotationTrue checks whether the rotate annotation is present and set to true or ""
func isRotateAnnotationTrue(obj client.Object) bool {
	if val, ok := obj.GetAnnotations()[rotateAnnotationKey]; ok && val == "" || val == "true" {
		return true
	}
	return false
}

// newOpValuesFromResource constructs a database.OpValues struct starting from a Database resource.
func newOpValuesFromResource(obj *databasev1.Database) (database.OpValues, ReconcileError) {
	metaIn := obj.ObjectMeta
	var metadata map[string]interface{}
	temp, _ := json.Marshal(metaIn)
	err := json.Unmarshal(temp, &metadata)
	if err != nil {
		return database.OpValues{}, ReconcileError{
			Reason:  RsnDbMetaParseFail,
			Message: MsgDbMetaParseFail,
			Err:     err,
		}
	}
	specIn := obj.Spec.Params
	var spec map[string]string
	temp, err = json.Marshal(specIn)
	err = json.Unmarshal(temp, &spec)
	if err != nil {
		return database.OpValues{}, ReconcileError{
			Reason:  RsnDbSpecParseFail,
			Message: MsgDbSpecParseFail,
			Err:     err,
		}
	}

	return database.OpValues{
		Metadata:   metadata,
		Parameters: spec,
	}, ReconcileError{}
}

// formatEventMessage formats an event message with key and values formatted as a json key-value structure. If keyAndValues
// is empty, it returns the message back.
func formatEventMessage(message string, keyAndValues ...interface{}) string {
	if len(keyAndValues) > 0 {
		extraValues := formatKeyAndValuesAsJson(keyAndValues)
		if extraValues != "" {
			return fmt.Sprintf("%s: %s", message, extraValues)
		}
	}
	return message
}

// formatKeyAndValuesAsJson converts a slice of interface{} into a json key-value string. Keys need to be strings by JSON's convention.
func formatKeyAndValuesAsJson(keyAndValues []interface{}) string {
	keyAndValuesLen := len(keyAndValues)
	if keyAndValuesLen%2 != 0 {
		logger.Error(fmt.Errorf("expected an even number of arguments, provided: %d", keyAndValuesLen),
			"odd number of keyAndValues provided!", keyAndValues...)
		// Set length to 0 so additional values are ignored
		keyAndValuesLen = 0
	}
	if keyAndValuesLen > 0 {
		keyAndValuesMap := make(map[string]interface{}, keyAndValuesLen/2)
		for i := 0; i < keyAndValuesLen; i += 2 {
			key, ok := keyAndValues[i].(string)
			if !ok {
				logger.Error(fmt.Errorf("expected key at position %d to be string, provided: %T", i, keyAndValues[i]),
					"keys must be strings!", keyAndValues...)
				// Set length to 0 so additional values are ignored
				keyAndValuesLen = 0
			}
			keyAndValuesMap[key] = keyAndValues[i+1]
		}
		str, err := json.Marshal(keyAndValuesMap)
		if err != nil {
			logger.Error(err, "cannot marshal keyAndValues to JSON")
			return ""
		}
		return string(str)
	}
	return ""
}

// contains is a very small utility function which returns true if s has been found in list.
func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// shouldIgnoreUpdateErr checks if an error message is generated due to the optimistic locking mechanism of Kubernetes API.
// This specific error is innocuous and should be generally ignored.
// See https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
func shouldIgnoreUpdateErr(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), genericregistry.OptimisticLockErrorMsg) {
		// do manual retry without error
		return true
	}
	return false
}
