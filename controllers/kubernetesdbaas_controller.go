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

package controllers

import (
	"context"
	"fmt"
	. "github.com/bedag/kubernetes-dbaas/api/v1alpha1"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	k8sError "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"math"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// KubernetesDbaasReconciler reconciles a KubernetesDbaas object
type KubernetesDbaasReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	dbaasResourceFinalizer = "finalizer.kubernetesdbaas.bedag.ch"
	// maxRequeueAfterDuration specifies the max wait time between two error retries
	maxRequeueAfterDuration = 6 * time.Hour
)

// +kubebuilder:rbac:groups=dbaas.bedag.ch,resources=kubernetesdbaas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.bedag.ch,resources=kubernetesdbaas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dbaas.bedag.ch,resources=kubernetesdbaas/finalizers,verbs=update
// +kubebuilder:rbac:groups="dbaas.bedag.ch",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
func (r *KubernetesDbaasReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("kubernetesdbaas", req.NamespacedName)
	logger.Info("Reconcile called.")
	dbaasResource := &KubernetesDbaas{}

	err := r.Get(ctx, req.NamespacedName, dbaasResource)

	if err != nil {
		// Fetch the KubernetesDbaas instance
		dbaasResource := &KubernetesDbaas{}
		err := r.Get(ctx, req.NamespacedName, dbaasResource)
		if err != nil {
			if k8sError.IsNotFound(err) {
				// Request object not found, could have been deleted after reconcile request.
				// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
				// Return and don't requeue
				logger.Info("KubernetesDbaas resource not found. Ignoring since object must be deleted.")
				// Not managed through ManageSuccess as the resource has been deleted and it will create problems if
				// the controller tries to manipulate its status
				return reconcile.Result{}, nil
			}
			// Error reading the object - requeue the request.
			logger.Error(err, "Failed to get KubernetesDbaas resource.")
			return r.ManageError(dbaasResource, err)
		}
	}

	// Check if the Memcached instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isDbaasResourceMarkedToBeDeleted := dbaasResource.GetDeletionTimestamp() != nil
	if isDbaasResourceMarkedToBeDeleted {
		if contains(dbaasResource.GetFinalizers(), dbaasResourceFinalizer) {
			// Run finalization logic for KubernetesDbaasFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			logger.Info("Finalizing resource...")
			if err := r.finalizeDbaasResource(logger, dbaasResource); err != nil {
				logger.Error(err, "Failed to get KubernetesDbaas resource.")
				return r.ManageError(dbaasResource, err)
			}

			// Remove KubernetesDbaasFinalizer. Once all finalizers have been
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
	return reconcile.Result{Requeue: false}, nil
}

// SetupWithManager creates the controller responsible for KubernetesDbaas resources by means of a ctrl.Manager.
func (r *KubernetesDbaasReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&KubernetesDbaas{}).
		Owns(&v1.Secret{}).
		Complete(r)
}

// ManageSuccess manages a successful reconciliation.
func (r *KubernetesDbaasReconciler) ManageSuccess(obj *KubernetesDbaas) (reconcile.Result, error) {
	obj.Status.LastError = ""
	obj.Status.LastUpdate = metav1.Now()
	obj.Status.LastErrorUpdateCount = 0

	err := r.Client.Status().Update(context.Background(), obj)
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
func (r *KubernetesDbaasReconciler) ManageError(obj *KubernetesDbaas, issue error) (reconcile.Result, error) {
	r.Log.Info("ManageError called.")
	if issue == nil {
		return r.ManageSuccess(obj)
	}

	// Set/Update error fields
	obj.Status.LastError = issue.Error()
	obj.Status.LastUpdate = metav1.Now()
	obj.Status.LastErrorUpdateCount++

	if obj.Status.LastErrorUpdateCount == 1 {
		// First iteration: wait just one second
		return reconcile.Result{RequeueAfter: time.Second}, issue
	} else if obj.Status.LastErrorUpdateCount > 14 {
		// 14 is calculated as 2^n = 21600[s] so at the 14th retry, the wait time will be approximately 6 hours
		// The controller was unable to fix the error by itself
		obj.Status.Unrecoverable = true
		return reconcile.Result{}, nil
	}

	// If the number of retries is at its maximum, the error is deemed unrecoverable
	tSinceLastRequeue := metav1.Now().Sub(obj.Status.LastUpdate.Time).Round(time.Second)

	// Double the wait time
	return reconcile.Result{
		// We just make sure the wait time between retries doesn't get too large by setting
		// maxRequeueAfterDuration as max wait time for a retry wait
		RequeueAfter: time.Duration(math.Min(float64(tSinceLastRequeue)*2, float64(maxRequeueAfterDuration))),
	}, issue
}

// finalizeDbaasResource cleans up resources not owned by dbaasResource.
func (r *KubernetesDbaasReconciler) finalizeDbaasResource(logger logr.Logger, dbaasResource *KubernetesDbaas) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	logger.Info("debug", "resource UID: %s", string(dbaasResource.UID))

	err := r.deleteDb(dbaasResource)
	if err != nil {
		return err
	}
	logger.Info("Successfully finalized dbaas resource")
	return nil
}

// addFinalizer adds a finalizer to a KubernetesDbaas resource.
func (r *KubernetesDbaasReconciler) addFinalizer(dbaasResource *KubernetesDbaas) error {
	r.Log.Info("Adding Finalizer for the KubernetesDbaas resource")
	controllerutil.AddFinalizer(dbaasResource, dbaasResourceFinalizer)

	// Update CR
	err := r.Update(context.TODO(), dbaasResource)
	if err != nil {
		r.Log.Error(err, "Failed to update KubernetesDbaas resource with finalizer")
		return err
	}
	return nil
}

// createDb creates a new database instance on the external provisioner based on the dbaasResource data.
func (r *KubernetesDbaasReconciler) createDb(dbaasResource *KubernetesDbaas) error {
	r.Log.Info("Creating database instance for: %s", dbaasResource.UID)
	conn, err := pool.GetConnByDriverAndEndpointName(dbaasResource.Spec.Provisioner, dbaasResource.Spec.Endpoint)
	if err != nil {
		return err
	}

	output := conn.CreateDb(string(dbaasResource.UID))
	if output.Err != nil {
		return fmt.Errorf("could not create database: %s", output.Err)
	}

	dsn, err := pool.GetDsnByDriverAndEndpointName(dbaasResource.Spec.Provisioner, dbaasResource.Spec.Endpoint)
	if err != nil {
		return err
	}

	// Create Secret
	err = r.createSecret(dbaasResource, output, dsn)
	if err != nil {
		return fmt.Errorf("could not create secret: %s", err)
	}

	return nil
}

// deleteDb deletes the database instance on the external provisioner based on the dbaasResource data.
func (r *KubernetesDbaasReconciler) deleteDb(dbaasResource *KubernetesDbaas) error {
	r.Log.Info("Deleting database instance for: %s", dbaasResource.UID)
	conn, err := pool.GetConnByDriverAndEndpointName(dbaasResource.Spec.Provisioner, dbaasResource.Spec.Endpoint)
	if err != nil {
		return err
	}

	if err = conn.DeleteDb(string(dbaasResource.UID)).Err; err != nil {
		return err
	}
	return nil
}

// createSecret creates a new K8s secret owned by owner with the data contained in output and dsn.
func (r *KubernetesDbaasReconciler) createSecret(owner *KubernetesDbaas, output database.QueryOutput, dsn database.Dsn) error {
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
			"dsn":      dsn.WithTable(output.Out[2]).String(),
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

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
