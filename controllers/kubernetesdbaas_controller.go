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
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dbaasv1alpha1 "github.com/bedag/kubernetes-dbaas/api/v1alpha1"
)

// KubernetesDbaasReconciler reconciles a KubernetesDbaas object
type KubernetesDbaasReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dbaas.bedag.ch,resources=kubernetesdbaas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.bedag.ch,resources=kubernetesdbaas/status,verbs=get;update;patch

func (r *KubernetesDbaasReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("kubernetesdbaas", req.NamespacedName)
	dbaasResource := &dbaasv1alpha1.KubernetesDbaas{}
	err := r.Get(ctx, req.NamespacedName, dbaasResource)

	if err != nil {
		// Delete
		if errors.IsNotFound(err) {
			// TODO: Encapsulate
			logger.Info("Deleting " + req.String() + "...")
			dbConn, err := database.New(dbaasResource.Spec.DbmsType)
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			if err != nil {
				logger.Error(err, "Failed to establish a DBMS connection")
				return ctrl.Result{}, err
			}

			err = dbConn.DeleteDb()
			if err != nil {
				logger.Error(err, "Failed to delete DB from DBMS")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get KubernetesDbaas object")
		return ctrl.Result{}, err
	}

	// Create
	// TODO: Encapsulate
	logger.Info("Creating " + req.String() + "...")
	logger.Info("dbname: " + dbaasResource.Spec.DbName)
	logger.Info("dbstage: " + dbaasResource.Spec.DbStage)
	logger.Info("dbmstype: " + dbaasResource.Spec.DbmsType)

	if !validateDbmsType(dbaasResource.Spec.DbmsType) {
		return ctrl.Result{}, fmt.Errorf("the following DBMS type: \"%s\" is not supported", dbaasResource.Spec.DbmsType)
	}

	dbConn, err := database.New(dbaasResource.Spec.DbmsType)
	if err != nil {
		logger.Error(err, "Failed to establish a DBMS connection")
		return ctrl.Result{}, err
	}
	outParams, err := dbConn.CreateDb(dbaasResource.Spec.DbName, dbaasResource.Spec.DbStage)
	if err != nil {
		logger.Error(err, "Failed to create DB instance")
		return ctrl.Result{}, err
	}
	// TODO: Refactor out params
	username := outParams[0]
	password := outParams[1]

	logger.Info("Creating secret...")

	err = r.createSecret(username, password, req.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *KubernetesDbaasReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbaasv1alpha1.KubernetesDbaas{}).
		Complete(r)
}

func validateDbmsType(s string) bool {
	for _, supportedDb := range database.GetSupportedDbms() {
		if supportedDb == s {
			return true
		}
	}
	return false
}

func (r *KubernetesDbaasReconciler) createSecret(username, password, namespace string) error {
	err := r.Client.Create(context.Background(), &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: namespace, // TODO: Set namespace where the CR is created
		},
		StringData: map[string]string{
			"username": username,
			"password": password,
		},
	})

	return err
}
