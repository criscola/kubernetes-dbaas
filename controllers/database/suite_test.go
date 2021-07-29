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

package controllers_test

import (
	"context"
	operatorconfigv1 "github.com/bedag/kubernetes-dbaas/apis/config/v1"
	databasev1 "github.com/bedag/kubernetes-dbaas/apis/database/v1"
	databaseclassv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
	. "github.com/bedag/kubernetes-dbaas/controllers/database"
	"github.com/bedag/kubernetes-dbaas/internal/logging"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"os"
	"path"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"testing"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var testEnv *envtest.Environment
var dbmsPool pool.DbmsPool
var rootPath = path.Join("..", "..")
var cfgFilepath = path.Join(rootPath, "config.yaml")
var testdataFilepath = path.Join(rootPath, "testdata")
var ctrlConfig = operatorconfigv1.OperatorConfig{}

const (
	DbcPostgresFilename  = "dbclass-postgres.yaml"
	DbcSqlserverFilename = "dbclass-sqlserver.yaml"
	DbcMariadbFilename   = "dbclass-mariadb.yaml"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t,"Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(logging.GetDevelopmentLogger(-2, false))
	Context("when setting up the test infrastructure", func() {
		var cfg *rest.Config
		var err error
		By("bootstrapping test environment", func() {
			if os.Getenv("TEST_USE_EXISTING_CLUSTER") == "true" {
				testEnv = &envtest.Environment{
					UseExistingCluster: &[]bool{true}[0],
				}
			} else {
				testEnv = &envtest.Environment{
					CRDDirectoryPaths:     []string{filepath.Join(rootPath, "config", "crd", "bases")},
					ErrorIfCRDPathMissing: true,
				}
			}
			if testConfigPath := os.Getenv("TEST_CONFIG_PATH"); testConfigPath != "" {
				cfgFilepath = path.Join(rootPath, testConfigPath)
			}
			cfg, err = testEnv.Start()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
		})
		By("adding CRD schemes", func() {
			err = databasev1.AddToScheme(scheme.Scheme)
			Expect(err).NotTo(HaveOccurred())
			err = databaseclassv1.AddToScheme(scheme.Scheme)
			Expect(err).NotTo(HaveOccurred())
			err = operatorconfigv1.AddToScheme(scheme.Scheme)
			Expect(err).NotTo(HaveOccurred())
			//+kubebuilder:scaffold:scheme
		})

		var options manager.Options
		var k8sManager ctrl.Manager
		By("loading the operator config from " + cfgFilepath, func() {
			options = ctrl.Options{Scheme: scheme.Scheme}
			dat, err := ioutil.ReadFile(cfgFilepath)
			Expect(err).NotTo(HaveOccurred())
			err = yaml.Unmarshal(dat, &ctrlConfig)
			Expect(err).NotTo(HaveOccurred())
			options, err = options.AndFrom(&ctrlConfig)
			Expect(err).NotTo(HaveOccurred())
		})
		By("creating Manager and Client", func() {
			k8sManager, err = ctrl.NewManager(cfg, options)
			Expect(err).ToNot(HaveOccurred())
			k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
			Expect(k8sClient).NotTo(BeNil())
		})
		By("registering DatabaseClasses", func() {
			dbcPostgres, err := getPostgresDbc()
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Create(context.Background(), &dbcPostgres)
			Expect(err).ToNot(HaveOccurred())

			dbcMariadb, err := getMariadbDbc()
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Create(context.Background(), &dbcMariadb)
			Expect(err).ToNot(HaveOccurred())

			dbcSqlserver, err := getSqlserverDbc()
			Expect(err).ToNot(HaveOccurred())
			err = k8sClient.Create(context.Background(), &dbcSqlserver)
			Expect(err).ToNot(HaveOccurred())
		})
		By("registering the pool of connections", func() {
			dbmsPool = pool.NewDbmsPool(0)
			for _, dbms := range ctrlConfig.DbmsList {
				dbc := databaseclassv1.DatabaseClass{}
				err = k8sClient.Get(context.Background(), client.ObjectKey{Name: dbms.DatabaseClassName}, &dbc)
				Expect(err).ToNot(HaveOccurred())

				err = dbmsPool.RegisterDbms(dbms, dbc.Spec.Driver)
				Expect(err).ToNot(HaveOccurred())
			}
		})
		By("starting the DatabaseReconciler instance", func() {
			err = (&DatabaseReconciler{
				Client:        k8sManager.GetClient(),
				Scheme:        k8sManager.GetScheme(),
				Log:           ctrl.Log.WithName("controllers").WithName("database"),
				EventRecorder: k8sManager.GetEventRecorderFor(DatabaseControllerName),
				DbmsList:      ctrlConfig.DbmsList,
				Pool:          dbmsPool,
			}).SetupWithManager(k8sManager)
			Expect(err).ToNot(HaveOccurred())

			go func() {
				err = k8sManager.Start(ctrl.SetupSignalHandler())
				Expect(err).ToNot(HaveOccurred())
			}()
		})

	})
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment", func() {
		err := testEnv.Stop()
		Expect(err).NotTo(HaveOccurred())
	})
})

func getSqlserverDbc() (databaseclassv1.DatabaseClass, error) {
	return readDbcYaml(DbcSqlserverFilename)
}

func getMariadbDbc() (databaseclassv1.DatabaseClass, error) {
	return readDbcYaml(DbcMariadbFilename)
}

func getPostgresDbc() (databaseclassv1.DatabaseClass, error) {
	return readDbcYaml(DbcPostgresFilename)
}

func readDbcYaml(filename string) (databaseclassv1.DatabaseClass, error) {
	dbcFilepath := path.Join(testdataFilepath, filename)
	dbc := databaseclassv1.DatabaseClass{}
	dat, err := ioutil.ReadFile(dbcFilepath)
	if err != nil {
		return databaseclassv1.DatabaseClass{}, err
	}
	err = yaml.Unmarshal(dat, &dbc)
	if err != nil {
		return databaseclassv1.DatabaseClass{}, err
	}
	return dbc, nil
}
