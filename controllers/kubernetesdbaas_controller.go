/*


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

// Package controllers contain the Kubernetes controllers responsible for their Custom Resource.
package controllers

import (
	"context"
	"fmt"
	"github.com/bedag/kubernetes-dbaas/pkg/config"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	"github.com/go-logr/logr"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/bedag/kubernetes-dbaas/api/v1"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	dbaasResourceFinalizer = "finalizer.database.bedag.ch"
	DateTimeLayout         = time.UnixDate
)

// Reconcile tries to reconcile the state of the cluster with the state of Database resources.
// +kubebuilder:rbac:groups=dbaas.bedag.ch,resources=database,verbs=list;watch;update
// +kubebuilder:rbac:groups=dbaas.bedag.ch,resources=database/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=list;watch;create;update;delete
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("database", req.NamespacedName)
	logger.Info("Reconcile called.")
	dbaasResource := &Database{}
	err := r.Get(ctx, req.NamespacedName, dbaasResource)

	if err != nil {
		// Fetch the Database instance
		dbaasResource := &Database{}
		err := r.Get(ctx, req.NamespacedName, dbaasResource)
		if err != nil {
			if k8sError.IsNotFound(err) {
				// Request object not found, could have been deleted after reconcile request.
				// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
				// Return and don't requeue
				logger.Info("Database resource not found. Ignoring since object must be deleted.")
				// Not managed through ManageSuccess as the resource has been deleted and it will create problems if
				// the controller tries to manipulate its status
				return reconcile.Result{}, nil
			}
			// Error reading the object - requeue the request.
			logger.Error(err, "Failed to get Database resource.")
			return r.ManageError(dbaasResource, err)
		}
	}

	// Check if the Database instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isDbaasResourceMarkedToBeDeleted := dbaasResource.GetDeletionTimestamp() != nil
	if isDbaasResourceMarkedToBeDeleted {
		if contains(dbaasResource.GetFinalizers(), dbaasResourceFinalizer) {
			// Run finalization logic for DatabaseFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			logger.Info("Finalizing resource...")
			if err := r.finalizeDbaasResource(logger, dbaasResource); err != nil {
				logger.Error(err, "Failed to get Database resource.")
				return r.ManageError(dbaasResource, err)
			}

			// Remove databaseFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			logger.Info("Removing finalizer...")
			controllerutil.RemoveFinalizer(dbaasResource, dbaasResourceFinalizer)
			if err := r.Update(ctx, dbaasResource); err != nil {
				logger.Error(err, "Error updating resource after removing finalizer...")
				return r.ManageError(dbaasResource, err)
			}
		}
		logger.Info("Resource finalized with success.")
		return r.ManageSuccess(dbaasResource)
	}

	// Add finalizer for this CR
	if !contains(dbaasResource.GetFinalizers(), dbaasResourceFinalizer) {
		logger.Info("Adding finalizer...")
		if err := r.addFinalizer(dbaasResource); err != nil {
			logger.Error(err, "Error adding finalizer")
			return r.ManageError(dbaasResource, err)
		}
	}

	// Create.
	logger.Info("Calling Create operation...")
	if err = r.createDb(dbaasResource); err != nil {
		logger.Error(err, "Failed to create resource")
		return r.ManageError(dbaasResource, err)
	}

	logger.Info("Reached end of Reconcile")
	return r.ManageSuccess(dbaasResource)
}

// SetupWithManager creates the controller responsible for Database resources by means of a ctrl.Manager.
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&Database{}).
		Owns(&v1.Secret{}).
		Complete(r)
}

// ManageSuccess manages a successful reconciliation.
func (r *DatabaseReconciler) ManageSuccess(obj *Database) (reconcile.Result, error) {
	r.Log.Info("ManageSuccess called.")

	// If the object is nil,
	if obj == nil {
		r.Log.Info("ManageSuccess called on nil resource, ignoring.")
		return reconcile.Result{}, nil
	}

	// If the object is marked unrecoverable, ignore
	if obj.Status.Unrecoverable {
		return reconcile.Result{}, nil
	}

	obj.Status.LastError = ""
	obj.Status.LastUpdate = metav1.Now().Format(DateTimeLayout)
	obj.Status.LastErrorUpdateCount = 0

	err := r.Status().Update(context.Background(), obj)
	if err != nil {
		// TODO: Implement conditions pre-update checks
		if k8sError.IsConflict(err) {
			return reconcile.Result{}, nil
		}
		return r.ManageError(obj, err)
	}
	return reconcile.Result{}, nil
}

// ManageSuccess manages a failed reconciliation.
func (r *DatabaseReconciler) ManageError(obj *Database, issue error) (reconcile.Result, error) {
	r.Log.Info("ManageError called.")
	if issue == nil || obj == nil {
		return r.ManageSuccess(obj)
	}

	// If error is unrecoverable, ignore
	if obj.Status.Unrecoverable {
		return reconcile.Result{}, nil
	}

	// If 3 seconds haven't passed since last requeue, don't do anything
	t, _ := time.Parse(DateTimeLayout, obj.Status.LastUpdate)
	timeSinceLastRequeue := time.Now().Sub(t)
	if timeSinceLastRequeue <= 3*time.Second {
		return reconcile.Result{}, nil
	}

	// Set/Update error fields
	obj.Status.LastError = issue.Error()
	obj.Status.LastUpdate = metav1.Now().Format(DateTimeLayout)
	obj.Status.LastErrorUpdateCount++

	if err := r.Status().Update(context.Background(), obj); err != nil {
		return reconcile.Result{}, err
	}

	if obj.Status.LastErrorUpdateCount < 14 {
		return reconcile.Result{RequeueAfter: 3 * time.Second}, nil
	}

	// The controller couldn't fix the problem itself
	obj.Status.Unrecoverable = true
	r.Log.Error(issue, "resource is in unrecoverable state")
	if err := r.Status().Update(context.Background(), obj); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

// finalizeDbaasResource cleans up resources not owned by dbaasResource.
func (r *DatabaseReconciler) finalizeDbaasResource(logger logr.Logger, dbaasResource *Database) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	logger.Info("Debug", "resource UID: %s", string(dbaasResource.UID))

	err := r.deleteDb(dbaasResource)
	if err != nil {
		return err
	}
	logger.Info("Successfully finalized dbaas resource")
	return nil
}

// addFinalizer adds a finalizer to a Database resource.
func (r *DatabaseReconciler) addFinalizer(dbaasResource *Database) error {
	r.Log.Info("Adding Finalizer for the Database resource")
	controllerutil.AddFinalizer(dbaasResource, dbaasResourceFinalizer)

	// Update CR
	err := r.Update(context.TODO(), dbaasResource)
	if err != nil {
		r.Log.Error(err, "Failed to update Database resource with finalizer")
		return err
	}
	return nil
}

// createDb creates a new database instance on the external provisioner based on the dbaasResource data.
func (r *DatabaseReconciler) createDb(dbaasResource *Database) error {
	r.Log.Info(fmt.Sprintf("Creating database instance for: %s", dbaasResource.UID))
	conn, err := pool.GetConnByDriverAndEndpointName(dbaasResource.Spec.Provisioner, dbaasResource.Spec.Endpoint)
	if err != nil {
		return err
	}

	dbms, err := config.
		GetDbmsConfig().
		GetByDriverAndEndpoint(dbaasResource.Spec.Provisioner, dbaasResource.Spec.Endpoint)
	if err != nil {
		return fmt.Errorf("could not get dbms entry: %s", err)
	}

	opValues, err := newOpValuesFromResource(dbaasResource)
	if err != nil {
		return fmt.Errorf("could not get generate operation values from resource: %s", err)
	}

	createOp, err := dbms.RenderOperation(database.CreateMapKey, opValues)
	if err != nil {
		return fmt.Errorf("could not render create operation values: %s", err)
	}

	output := conn.CreateDb(createOp)
	if output.Err != nil {
		return fmt.Errorf("could not create database: %s", output.Err)
	}

	// Create Secret
	err = r.createSecret(dbaasResource, output)
	if err != nil {
		return fmt.Errorf("could not create secret: %s", err)
	}

	return nil
}

// deleteDb deletes the database instance on the external provisioner based on the dbaasResource data.
func (r *DatabaseReconciler) deleteDb(dbaasResource *Database) error {
	r.Log.Info(fmt.Sprintf("Deleting database instance for: %s", dbaasResource.UID))
	conn, err := pool.GetConnByDriverAndEndpointName(dbaasResource.Spec.Provisioner, dbaasResource.Spec.Endpoint)
	if err != nil {
		return err
	}

	dbms, err := config.
		GetDbmsConfig().
		GetByDriverAndEndpoint(dbaasResource.Spec.Provisioner, dbaasResource.Spec.Endpoint)
	if err != nil {
		return fmt.Errorf("could not get dbms entry: %s", err)
	}

	opValues, err := newOpValuesFromResource(dbaasResource)
	if err != nil {
		return fmt.Errorf("could not get generate operation values from resource: %s", err)
	}

	deleteOp, err := dbms.RenderOperation(database.DeleteMapKey, opValues)
	if err != nil {
		return fmt.Errorf("could not render create operation values: %s", err)
	}

	output := conn.DeleteDb(deleteOp)
	if output.Err != nil {
		return fmt.Errorf("could not delete database: %s", output.Err)
	}

	return nil
}

// createSecret creates a new K8s secret owned by owner with the data contained in output and dsn.
func (r *DatabaseReconciler) createSecret(owner *Database, output database.OpOutput) error {
	var ownerRefs []metav1.OwnerReference

	ownerRefs = append(ownerRefs, metav1.OwnerReference{
		APIVersion: owner.APIVersion,
		Kind:       owner.Kind,
		Name:       owner.Name,
		UID:        owner.UID,
		Controller: &[]bool{true}[0], // sets this controller as owner
	})

	newSecret := &v1.Secret{
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
			"dsn":      database.NewDsn(owner.Spec.Provisioner, output.Out[0], output.Out[1], output.Out[3], output.Out[4], output.Out[2]).String(),
		},
	}

	oldSecret := v1.Secret{}
	key := client.ObjectKey{
		Namespace: owner.Namespace,
		Name:      owner.Name + "-credentials",
	}
	err := r.Client.Get(context.TODO(), key, &oldSecret)
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
func newOpValuesFromResource(resource *Database) (database.OpValues, error) {
	metaIn := resource.ObjectMeta
	var metadata map[string]interface{}
	temp, _ := json.Marshal(metaIn)
	err := json.Unmarshal(temp, &metadata)
	if err != nil {
		return database.OpValues{}, err
	}
	specIn := resource.Spec.Params
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
