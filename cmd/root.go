/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"context"
	"github.com/bedag/kubernetes-dbaas/internal/logging"
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	"github.com/go-logr/logr"

	//"context"
	"fmt"
	operatorconfigv1 "github.com/bedag/kubernetes-dbaas/apis/config/v1"
	databasev1 "github.com/bedag/kubernetes-dbaas/apis/database/v1"
	databaseclassv1 "github.com/bedag/kubernetes-dbaas/apis/databaseclass/v1"
	controllers "github.com/bedag/kubernetes-dbaas/controllers/database"

	//"github.com/bedag/kubernetes-dbaas/pkg/pool"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	//+kubebuilder:scaffold:imports

	"github.com/spf13/viper"
)

const (
	LoadConfigKey    = "load-config"
	DebugKey         = "debug"
	WebhookEnableKey = "enable-webhooks"
	ZapLogLevelKey   = "log-level"
	DisableStackTrace = "disable-stacktrace"
	RpsKey 			  = "rps"

	// Flag overrides for flags specified in OperatorConfig
	MetricsBindAddressKey     = "metrics.bindAddress"
	HealthProbeBindAddressKey = "health.healthProbeBindAddress"
	LeaderElectEnableKey      = "leaderElection.leaderElect"
	LeaderElectResName        = "leaderElection.resourceName"
	WebhookPortKey            = "webhook.port"
)

var (
	dbmsPool pool.DbmsPool
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "kubedbaas",
	Short: "kubedbaas is a Kubernetes Operator written in Go used to provision databases on external infrastructure",
	Long: `A Kubernetes Operator able to trigger stored procedures in external DBMS which in turn provision new database instances.
				Users are able to create new database instances by writing an API Object configuration using Custom Resources.
				The Operator watches for new API Objects and tells the target DBMS to trigger a certain stored procedure based on the content of the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLog.Info("registering endpoints...")
		registerEndpoints()
		setupLog.Info("endpoints registered")

		setupLog.Info("starting operator")
		loadOperator()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfigFile, initLogger)

	initFlags()

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(operatorconfigv1.AddToScheme(scheme))
	utilruntime.Must(databasev1.AddToScheme(scheme))
	utilruntime.Must(databaseclassv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func initFlags() {
	rootCmd.PersistentFlags().String(LoadConfigKey, "", "The location of the Operator's config file")
	rootCmd.PersistentFlags().Bool(DebugKey, false, "Enable debug mode for development purposes. If set, --log-level defaults to 1")
	rootCmd.PersistentFlags().Bool(WebhookEnableKey, true, "Enable webhooks servers.")
	rootCmd.PersistentFlags().String(MetricsBindAddressKey, ":8080", "The address the metric endpoint binds to")
	rootCmd.PersistentFlags().String(HealthProbeBindAddressKey, ":8081", "The address the probe endpoint binds to")
	rootCmd.PersistentFlags().Bool(LeaderElectEnableKey, true, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager")
	rootCmd.PersistentFlags().String(LeaderElectResName, "bfa62c96.dbaas.bedag.ch", "The resource name to lock during election cycles")
	rootCmd.PersistentFlags().Int(WebhookPortKey, 9443, "The port the webhook server binds to")
	rootCmd.PersistentFlags().Int(ZapLogLevelKey, 0, "The verbosity of the logging output. Can be one out of: 0 info, 1 debug, 2 trace. If debug mode is on, defaults to 1")
	rootCmd.PersistentFlags().Bool(DisableStackTrace, false, "Disable stacktrace printing in logger errors")
	rootCmd.PersistentFlags().Int(RpsKey, 0, "The number of operation executed per second per endpoint. If set to 0, operations won't be rate-limited.")

	// Bind all flags to Viper
	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(flag.Name, flag)
	})
}

// initConfig reads in the Operator's configuration.
func initConfigFile() {
	if cfgFile := viper.GetString(LoadConfigKey); cfgFile != "" {
		// Use config file set from the flag.
		viper.SetConfigFile(viper.GetString("load-config"))
	} else {
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/kubedbaas")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// A config file must be specified.
		_, _ = fmt.Fprintf(os.Stderr, "error reading config file (%s): %s\n", viper.ConfigFileUsed(), err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintln(os.Stdout, "config file loaded:", viper.ConfigFileUsed())
}

// initLogger initializes the Operator's logger.
func initLogger() {
	// Distinguish between 'debug' and 'production' setting.
	level := viper.GetInt(ZapLogLevelKey)
	if level > 0 {
		level -= 2 * level
	}

	var logger logr.Logger
	if viper.GetBool(DebugKey) {
		fmt.Println("setting up logger in development mode...")
		// Check if default is set
		if !viper.IsSet(ZapLogLevelKey) {
			level = -1
		}
		logger = logging.GetDevelopmentLogger(level, viper.GetBool(DisableStackTrace))
	} else {
		fmt.Println("setting up logger in production mode...")
		logger = logging.GetProductionLogger(level, viper.GetBool(DisableStackTrace))
	}
	ctrl.SetLogger(logger)
}

// loadOperator registers all the Manager's controllers, webhooks and starts them.
func loadOperator() {
	// Load the Operator configuration
	// Set CLI flags given by user or set by default
	var err error
	ctrlConfig := operatorconfigv1.OperatorConfig{}
	options := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     viper.GetString(MetricsBindAddressKey),
		Port:                   viper.GetInt(WebhookPortKey),
		HealthProbeBindAddress: viper.GetString(HealthProbeBindAddressKey),
		LeaderElection:         viper.GetBool(LeaderElectEnableKey),
		LeaderElectionID:       viper.GetString(LeaderElectResName),
	}

	// Build and pass the configuration file to the controller.
	if cfgFile := viper.GetString(LoadConfigKey); cfgFile != "" {
		options, err = options.AndFrom(ctrl.ConfigFile().AtPath(cfgFile).OfKind(&ctrlConfig))

		if err != nil {
			fatalError(err, "unable to load the config file into the controller")
		}
	}

	// TODO: Check status of https://github.com/kubernetes-sigs/controller-runtime/issues/1463
	options.LeaderElection = viper.GetBool(LeaderElectEnableKey)

	// Initialize manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		fatalError(err, "unable to initialize manager")
	}

	// Setup controllers
	dbmsList, err := getDbmsList()
	if err != nil {
		fatalError(err, "unable to get dbms list")
	}

	if err = (&controllers.DatabaseReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("Database"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor(controllers.DatabaseControllerName),
		DbmsList: dbmsList,
		Pool: dbmsPool,
	}).SetupWithManager(mgr); err != nil {
		fatalError(err, "unable to create controller", "controller", "Database")
	}

	// Setup webhooks
	if viper.GetBool(WebhookEnableKey) {
		if err = (&databasev1.Database{}).SetupWebhookWithManager(mgr); err != nil {
			fatalError(err, "unable to create webhook", "webhook", "Database")
		}
	}

	//+kubebuilder:scaffold:builder

	// Setup probes
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		fatalError(err, "unable to set up health check")
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		fatalError(err, "unable to set up ready check")
	}

	// Finally start controllers and webhooks
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fatalError(err, "problem running manager")
	}
}

// getDbmsList returns the dbms endpoint list of the operator as specified in the operator's configuration file stored
// in Viper.
func getDbmsList() (database.DbmsList, error) {
	dbms := database.DbmsList{}
	if err := viper.UnmarshalKey(database.DbmsConfigKey, &dbms); err != nil {
		return nil, err
	}

	return dbms, nil
}

// RegisterEndpoints attempts to register the endpoints specified in the operator configuration loaded from LoadConfig.
//
// See pool.Register for details.
func registerEndpoints() {
	c, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		setupLog.Error(err, "unable to create client instance")
		os.Exit(1)
	}
	dbmsList, err := getDbmsList()
	if err != nil {
		fatalError(err, "unable to retrieve dbms configuration")
	}
	dbmsPool = pool.NewDbmsPool(viper.GetInt(RpsKey))
	for _, dbms := range dbmsList {
		dbClass := databaseclassv1.DatabaseClass{}
		err = c.Get(context.Background(), client.ObjectKey{Namespace: "", Name: dbms.DatabaseClassName}, &dbClass)
		if err != nil {
			fatalError(err, "problem getting databaseclass from api server", "databaseClassName",
				dbms.DatabaseClassName)
		}

		if err := dbmsPool.RegisterDbms(dbms, dbClass.Spec.Driver); err != nil {
			fatalError(err, "problem registering dbms endpoint", "databaseClassName", dbClass.Name)
		}
	}
}

func fatalError(err error, msg string, values ...interface{}) {
	setupLog.Error(err, msg, values...)
	os.Exit(1)
}
