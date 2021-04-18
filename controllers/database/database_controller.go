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
	"github.com/bedag/kubernetes-dbaas/pkg/config"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	. "github.com/bedag/kubernetes-dbaas/pkg/typeutil"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	databasev1 "github.com/bedag/kubernetes-dbaas/apis/database/v1"
	databaseclassv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
)

const (
	DatabaseControllerName = "database-controller"
	DatabaseClass          = "databaseclass"
	EndpointName           = "endpoint-name"
	databaseFinalizer      = "finalizer.database.bedag.ch"
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
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// https://pastebin.com/FXLmBa13
	obj := &databasev1.Database{}
	_ = r.Get(ctx, req.NamespacedName, obj)
	meta.SetStatusCondition(&obj.Status.Conditions, metav1.Condition{
		Type:    TypeReady,
		Status:  metav1.ConditionTrue,
		Reason:  "Test",
		Message: "this is a test",
	})
	err := r.Status().Update(ctx, obj)
	if err != nil {
		fmt.Println(err)
	}
	newObj := &databasev1.Database{}
	_ = r.Get(ctx, req.NamespacedName, newObj)
	return ctrl.Result{}, nil
}

// addFinalizer adds a finalizer to a Database resource.
func (r *DatabaseReconciler) addFinalizer(obj *databasev1.Database) error {
	controllerutil.AddFinalizer(obj, databaseFinalizer)

	// Update CR
	return r.Update(context.Background(), obj)
}

// createDb creates a new Database instance on the external provisioner based on the Database data.
func (r *DatabaseReconciler) createDb(obj *databasev1.Database) ReconcileError {
	logger.Info(MsgDbCreateInProg)
	r.EventRecorder.Event(obj, Normal, RsnDbCreateInProg, MsgDbCreateInProg)

	dbClass, err := r.getDbmsClassFromDb(obj)
	if err.IsNotEmpty() {
		return err
	}
	loggingKv := stringsToInterfaceSlice(DatabaseClass, dbClass.Name, database.OperationsConfigKey, database.CreateMapKey)

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
	conn, simpleErr := pool.GetConnByEndpointName(obj.Spec.Endpoint)
	if simpleErr != nil {
		return ReconcileError{
			Reason:         RsnDbmsEndpointConnFail,
			Message:        MsgDbmsEndpointConnFail,
			Err:            simpleErr,
			AdditionalInfo: loggingKv,
		}
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

	// Create Secret
	simpleErr = r.createSecret(obj, dbClass.Spec.Driver, output)
	if simpleErr != nil {
		return ReconcileError{
			Reason:         RsnSecretCreateFail,
			Message:        MsgSecretCreateFail,
			Err:            simpleErr,
			AdditionalInfo: loggingKv,
		}
	}

	return ReconcileError{}
}

// deleteDb deletes the database instance on the external provisioner.
func (r *DatabaseReconciler) deleteDb(obj *databasev1.Database) ReconcileError {
	logger.Info(MsgDbDeleteInProg)
	r.EventRecorder.Event(obj, Normal, RsnDbDeleteInProg, MsgDbDeleteInProg)

	dbClass, err := r.getDbmsClassFromDb(obj)
	if err.IsNotEmpty() {
		return err
	}
	loggingKv := stringsToInterfaceSlice(DatabaseClass, dbClass.Name, database.OperationsConfigKey, database.DeleteMapKey)

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
	opValues, err := newOpValuesFromResource(obj)
	if err.IsNotEmpty() {
		return err.With(loggingKv)
	}
	deleteOp, simpleErr := deleteOpTemplate.RenderOperation(opValues)
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
	conn, simpleErr := pool.GetConnByEndpointName(obj.Spec.Endpoint)
	if simpleErr != nil {
		return ReconcileError{
			Reason:         RsnDbmsEndpointConnFail,
			Message:        MsgDbmsEndpointConnFail,
			Err:            simpleErr,
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

func (r *DatabaseReconciler) getDbmsClassFromDb(obj *databasev1.Database) (databaseclassv1.DatabaseClass, ReconcileError) {
	// Get DatabaseClass name from DBMS config
	dbmsList, err := config.GetDbmsList()
	if err != nil {
		return databaseclassv1.DatabaseClass{}, ReconcileError{
			Reason:  RsnDbmsConfigGetFail,
			Message: MsgDbmsConfigGetFail,
			Err:     err,
		}
	}
	// Get DatabaseClass resource from api server
	dbClassName, err := dbmsList.GetDatabaseClassNameByEndpointName(obj.Spec.Endpoint)
	if err != nil {
		return databaseclassv1.DatabaseClass{}, ReconcileError{
			Reason:         RsnDbcConfigGetFail,
			Message:        MsgDbcConfigGetFail,
			Err:            err,
			AdditionalInfo: stringsToInterfaceSlice(DatabaseClass, dbClassName, EndpointName, obj.Spec.Endpoint),
		}
	}

	dbClass := databaseclassv1.DatabaseClass{}
	err = r.Client.Get(context.Background(), client.ObjectKey{Namespace: "", Name: dbClassName}, &dbClass)
	if err != nil {
		return databaseclassv1.DatabaseClass{}, ReconcileError{
			Reason:         RsnDbcGetFail,
			Message:        MsgDbcGetFail,
			Err:            err,
			AdditionalInfo: stringsToInterfaceSlice(DatabaseClass, dbClassName),
		}
	}
	return dbClass, ReconcileError{}
}

// defaultErrHandling sets the obj Conditions type Ready to false and sets the relative fields error and message,
// it records a Warning event with reason and message for the given obj and logs err (if present) and message to the
// global logger.
func (r *DatabaseReconciler) defaultErrHandling(obj *databasev1.Database, err ReconcileError) {
	keyAndValuesLen := len(err.AdditionalInfo)
	if keyAndValuesLen%2 != 0 {
		panic("odd number of keyAndValues provided!")
	}
	if keyAndValuesLen > 0 {
		eventMessage := fmt.Sprintf("%s: ", err.Message)
		for i := 0; i < keyAndValuesLen; i += 2 {
			if i != keyAndValuesLen-2 {
				eventMessage += fmt.Sprintf(`{"%s": "%v"}, `, err.AdditionalInfo[i], err.AdditionalInfo[i+1])
				continue
			}
			eventMessage += fmt.Sprintf(`{"%s": "%v"}`, err.AdditionalInfo[i], err.AdditionalInfo[i+1])
		}
		r.EventRecorder.Event(obj, Warning, err.Reason, eventMessage)
		logger.Error(err.Err, err.Message, err.AdditionalInfo...)
	} else {
		r.EventRecorder.Event(obj, Warning, err.Reason, err.Message)
		logger.Error(err.Err, err.Message)
	}
	if updateErr := r.updateReadyCondition(obj, metav1.ConditionFalse, err.Reason, err.Message); updateErr != nil {
		logger.Error(err.Err, MsgDbUpdateFail)
	}
}

// createSecret creates a new K8s secret owned by owner with the data contained in output and dsn.
func (r *DatabaseReconciler) createSecret(owner *databasev1.Database, driver string, output database.OpOutput) error {
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
			Name:            owner.Name + "-credentials",
			Namespace:       owner.Namespace,
			OwnerReferences: ownerRefs,
		},
		StringData: map[string]string{
			"username": output.Out[0],
			"password": output.Out[1],
			"host":     output.Out[3],
			"port":     output.Out[4],
			"dbName":   output.Out[2],
			"dsn":      database.NewDsn(driver, output.Out[0], output.Out[1], output.Out[3], output.Out[4], output.Out[2]).String(),
		},
	}

	oldSecret := corev1.Secret{}
	key := client.ObjectKey{
		Namespace: owner.Namespace,
		Name:      owner.Name + "-credentials",
	}
	err := r.Client.Get(context.Background(), key, &oldSecret)
	if err != nil {
		if k8sError.IsNotFound(err) {
			return r.Client.Create(context.Background(), newSecret)
		}
		return err
	}

	// Secret exists, overwrite
	return r.Client.Update(context.Background(), newSecret)
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

// updateReadyCondition updates the Ready Condition status of obj.
func (r *DatabaseReconciler) updateReadyCondition(obj *databasev1.Database, status metav1.ConditionStatus, reason, message string) error {
	meta.SetStatusCondition(&obj.Status.Conditions, metav1.Condition{
		Type:    TypeReady,
		Status:  status,
		Reason:  reason,
		Message: message,
	})

	// Update condition field
	return r.Client.Update(context.Background(), obj)
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

func stringsToInterfaceSlice(values ...string) []interface{} {
	y := make([]interface{}, len(values))
	for i, v := range values {
		y[i] = v
	}
	return y
}

// IsNotEmpty checks if r is not empty using reflect.DeepEqual. Needed because field AdditionalInfo is not comparable.
func (r ReconcileError) IsNotEmpty() bool {
	return !reflect.DeepEqual(r, ReconcileError{})
}

// With creates a copy of the receiver and appends values to its AdditionalInfo field.
func (r ReconcileError) With(values ...interface{}) ReconcileError {
	return ReconcileError{
		Reason:         r.Reason,
		Message:        r.Message,
		Err:            r.Err,
		AdditionalInfo: append(r.AdditionalInfo, values...),
	}
}
