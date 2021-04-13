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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	"github.com/go-logr/logr"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	databasev1 "github.com/bedag/kubernetes-dbaas/apis/database/v1"
	databaseclassv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	EventRecorder record.EventRecorder
}

const (
	DatabaseControllerName = "database-controller"
	DatabaseClass	       = "databaseclass"
	EndpointName           = "endpoint-name"
	databaseFinalizer      = "finalizer.database.bedag.ch"
)

var logger logr.Logger

// +kubebuilder:rbac:groups=database.dbaas.bedag.ch,resources=databases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.dbaas.bedag.ch,resources=databases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=database.dbaas.bedag.ch,resources=databases/finalizers,verbs=update
// +kubebuilder:rbac:groups=database.dbaas.bedag.ch,resources=databases/events,verbs=update
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
	logger = r.Log.V(0).WithValues("database", req.NamespacedName)
	logger.Info("Reconcile called")

	obj := &databasev1.Database{}
	err := r.Get(ctx, req.NamespacedName, obj)

	if err != nil {
		if k8sError.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info(MsgDbDeleted)
			// Not managed through ManageSuccess as the resource has been deleted and it will create problems if
			// the controller tries to manipulate its status
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, r.defaultErrHandling(obj, err, RsnDbGetFail, MsgDbGetFail)
	}

	// Set reason to unknown to indicate the resource was correctly received by a controller but no action was resolved yet
	meta.SetStatusCondition(&obj.Status.Conditions, metav1.Condition{
		Type:               TypeReady,
		Status:             metav1.ConditionUnknown,
		Reason:             "",
		Message:            "",
	})

	// Check if the Database instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if obj.GetDeletionTimestamp() != nil {
		if contains(obj.GetFinalizers(), databaseFinalizer) {
			// Run finalization logic for DatabaseFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			logger.Info("Finalizing database resource")
			if err := r.deleteDb(obj); err != nil {
				// Errors handled internally in r.deleteDb
				return reconcile.Result{}, err
			}

			// Remove databaseFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			logger.Info("Removing finalizer")
			controllerutil.RemoveFinalizer(obj, databaseFinalizer)
			if err := r.Update(ctx, obj); err != nil {
				logger.Error(err, MsgDbUpdateFail)
				r.EventRecorder.Event(obj, Warning, RsnDbUpdateFail, MsgDbUpdateFail)
				meta.SetStatusCondition(&obj.Status.Conditions, metav1.Condition{
					Type:    TypeReady,
					Status:  metav1.ConditionFalse,
					Reason:  RsnDbUpdateFail,
					Message: MsgDbUpdateFail,
				})
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Create
	if err = r.createDb(obj); err != nil {
		logger.Error(err, MsgDbProvisionFail)
		return reconcile.Result{}, err
	}

	// If finalizer is not present, add finalizer to resource
	if !contains(obj.GetFinalizers(), databaseFinalizer) {
		logger.Info("Adding finalizer")
		if err := r.addFinalizer(obj); err != nil {
			logger.Error(err, "Error adding finalizer")
			return reconcile.Result{}, err
		}
	}

	meta.SetStatusCondition(&obj.Status.Conditions, metav1.Condition{
		Type:    TypeReady,
		Status:  metav1.ConditionTrue,
		Reason:  RsnCreate,
		Message: MsgDbProvisionSucc,
	})

	logger.Info("Reached end of reconcile")
	return reconcile.Result{}, nil
}

// addFinalizer adds a finalizer to a Database resource.
func (r *DatabaseReconciler) addFinalizer(obj *databasev1.Database) error {
	logger.Info("Adding finalizer to the database resource")
	controllerutil.AddFinalizer(obj, databaseFinalizer)

	// Update CR
	err := r.Update(context.Background(), obj)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbUpdateFail, MsgDbUpdateFail)
	}
	return nil
}

// createDb creates a new Database instance on the external provisioner based on the Database data.
func (r *DatabaseReconciler) createDb(obj *databasev1.Database) error {
	logger.Info(MsgDbCreateInProg)
	r.EventRecorder.Event(obj, Normal, RsnDbCreateInProg, MsgDbCreateInProg)

	// TODO: Extract in method getDbClass(obj databasev1.Database) databaseclassv1.DatabaseClass
	// Get DatabaseClass name from DBMS config
	dbmsList, err := config.GetDbmsList()
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbmsConfigGetFail, MsgDbmsConfigGetFail)
	}

	// Get DatabaseClass resource from api server
	dbClassName, err := dbmsList.GetDatabaseClassNameByEndpointName(obj.Spec.Endpoint)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbcConfigGetFail, MsgDbcConfigGetFail,
			DatabaseClass, dbClassName, EndpointName, obj.Spec.Endpoint)

	}
	// Add some logging values
	loggingKv := stringSliceToInterfaceSlice([]string{DatabaseClass, dbClassName, database.OperationsConfigKey, database.CreateMapKey})

	dbClass := databaseclassv1.DatabaseClass{}
	err = r.Client.Get(context.Background(), client.ObjectKey{Namespace: "", Name: dbClassName}, &dbClass)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbcGetFail, MsgDbcGetFail, loggingKv...)
	}

	// Render operation
	createOpTemplate, exists := dbClass.Spec.Operations[database.CreateMapKey]
	if !exists {
		return r.defaultErrHandling(obj, nil, RsnOpNotSupported, MsgOpNotSupported, loggingKv...)
	}
	opValues, err := newOpValuesFromResource(obj)
	if err != nil {
		return r.defaultErrHandling(obj, nil, RsnOpValuesCreateFail, MsgOpValuesCreateFail, loggingKv...)
	}
	createOp, err := createOpTemplate.RenderOperation(opValues)
	if err != nil {
		return r.defaultErrHandling(obj, nil, RsnOpRenderFail, MsgOpRenderFail, loggingKv...)
	}

	// Execute operation on DBMS
	loggingKv = append(loggingKv, EndpointName, obj.Spec.Endpoint)
	conn, err := pool.GetConnByEndpointName(obj.Spec.Endpoint)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbmsEndpointConnFail, MsgDbmsEndpointConnFail, loggingKv...)
	}
	output := conn.CreateDb(createOp)
	if output.Err != nil {
		return r.defaultErrHandling(obj, err, RsnDbCreateFail, MsgDbCreateFail, loggingKv...)
	}

	// Create Secret
	err = r.createSecret(obj, dbClass.Spec.Driver, output)
	if err != nil {
		return  r.defaultErrHandling(obj, err, RsnSecretCreateFail, MsgSecretCreateFail, loggingKv...)
	}

	return nil
}

// deleteDb deletes the database instance on the external provisioner.
func (r *DatabaseReconciler) deleteDb(obj *databasev1.Database) error {
	logger.Info(MsgDbDeleteInProg)
	r.EventRecorder.Event(obj, Normal, RsnDbDeleteInProg, MsgDbDeleteInProg)

	// TODO: Extract in method getDbClass(obj databasev1.Database) databaseclassv1.DatabaseClass
	// Get DatabaseClass resource
	dbmsList, err := config.GetDbmsList()
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbmsConfigGetFail, MsgDbmsConfigGetFail)
	}
	dbClassName, err := dbmsList.GetDatabaseClassNameByEndpointName(obj.Spec.Endpoint)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbcConfigGetFail, MsgDbcConfigGetFail,
			DatabaseClass, dbClassName, EndpointName, obj.Spec.Endpoint)
	}
	// Add some logging values
	loggingKv := stringSliceToInterfaceSlice([]string{DatabaseClass, dbClassName, database.OperationsConfigKey, database.DeleteMapKey})

	dbClass := databaseclassv1.DatabaseClass{}
	err = r.Client.Get(context.Background(), client.ObjectKey{Namespace: "", Name: dbClassName}, &dbClass)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbcGetFail, MsgDbcGetFail, loggingKv...)
	}

	// Render operation
	deleteOpTemplate, exists := dbClass.Spec.Operations[database.DeleteMapKey]
	if !exists {
		return r.defaultErrHandling(obj, nil, RsnOpNotSupported, MsgOpNotSupported, loggingKv...)
	}
	opValues, err := newOpValuesFromResource(obj)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnOpValuesCreateFail, MsgOpValuesCreateFail, loggingKv...)
	}
	deleteOp, err := deleteOpTemplate.RenderOperation(opValues)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnOpRenderFail, MsgOpRenderFail, loggingKv...)
	}

	// Execute operation on DBMS
	loggingKv = append(loggingKv, EndpointName, obj.Spec.Endpoint)
	conn, err := pool.GetConnByEndpointName(obj.Spec.Endpoint)
	if err != nil {
		return r.defaultErrHandling(obj, err, RsnDbmsEndpointConnFail, MsgDbmsEndpointConnFail, loggingKv...)
	}
	output := conn.DeleteDb(deleteOp)
	if output.Err != nil {
		return r.defaultErrHandling(obj, err, RsnDbDeleteFail, MsgDbDeleteFail, loggingKv...)
	}

	return nil
}

// defaultErrHandling sets the obj Conditions type Ready to false and sets the relative fields error and message,
// it records a Warning event with reason and message for the given obj and
// logs err (if present) and message to the global logger. It returns an error formatted as following: "<message>: <err>".
func (r *DatabaseReconciler) defaultErrHandling(obj *databasev1.Database, err error, reason, message string, keyAndValues ...interface{}) error {
	keyAndValuesLen := len(keyAndValues)
	if keyAndValuesLen % 2 != 0 {
		panic("odd number of keyAndValues provided!")
	}
	meta.SetStatusCondition(&obj.Status.Conditions, metav1.Condition{
		Type:    TypeReady,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})
	if keyAndValuesLen > 0 {
		eventMessage := fmt.Sprintf("%s: ", message)
		for i := 0; i < keyAndValuesLen; i+=2 {
			if i != keyAndValuesLen-2 {
				eventMessage += fmt.Sprintf(`{ "%s": "%v" }, `, keyAndValues[i], keyAndValues[i+1])
				continue
			}
			eventMessage += fmt.Sprintf(`{ "%s": "%v" }`, keyAndValues[i], keyAndValues[i+1])
		}
		r.EventRecorder.Event(obj, Warning, reason, eventMessage)
		logger.Error(err, message, keyAndValues...)
	} else {
		r.EventRecorder.Event(obj, Warning, reason, message)
		logger.Error(err, message)
	}
	if err != nil {
		return fmt.Errorf("%s: %s", message, err)
	}
	return fmt.Errorf(message)
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
func newOpValuesFromResource(obj *databasev1.Database) (database.OpValues, error) {
	metaIn := obj.ObjectMeta
	var metadata map[string]interface{}
	temp, _ := json.Marshal(metaIn)
	err := json.Unmarshal(temp, &metadata)
	if err != nil {
		return database.OpValues{}, err
	}
	specIn := obj.Spec.Params
	var spec map[string]string
	temp, err = json.Marshal(specIn)
	err = json.Unmarshal(temp, &spec)
	if err != nil {
		return database.OpValues{}, err
	}

	// Ensure meta.namespace and meta.name are set
	if metadata["namespace"] == "" {
		metadata["namespace"] = "default"
	}
	if metadata["name"] == "" {
		// Generate name randomly
		metadata["name"] = randSeq(16)
	}

	return database.OpValues{
		Metadata:   metadata,
		Parameters: spec,
	}, nil
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

// randSeq generates a random alphanumeric string of length n
func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())

	var alphabet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func stringSliceToInterfaceSlice(x []string) []interface{} {
	y := make([]interface{}, len(x))
	for i, v := range x {
		y[i] = v
	}
	return y
}