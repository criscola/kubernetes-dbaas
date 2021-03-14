// Package cmd contains the CLI command definitions along with their initialization and startup code
package cmd

import (
	"fmt"
	dbaasv1 "github.com/bedag/kubernetes-dbaas/api/v1"
	"github.com/bedag/kubernetes-dbaas/controllers"
	"github.com/bedag/kubernetes-dbaas/pkg/config"
	"github.com/bedag/kubernetes-dbaas/pkg/pool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	metricsAddr          string
	enableLeaderElection bool
	cfgFile              string

	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const (
	DefaultConfigType	  = "yaml"
	DefaultConfigFilename = "config"
	DefaultConfigFilepath = "/var/kubedbaas"
	ConfigLoadError       = "problem loading configuration file"
	ConfigParseError      = "problem parsing configuration file"
	DbmsConnOpenError     = "problem opening a connection to DBMS endpoint"
)

// rootCmd represents the root 'kubedbaas' command
var rootCmd = &cobra.Command{
	Use:   "kubedbaas",
	Short: "kubedbaas is a Kubernetes Operator written in Go used to provision databases on external infrastructure",
	Long: `A Kubernetes Operator able to trigger stored procedures in external DBMS which in turn provision new database instances.
				Users are able to create new database instances by writing an API Object configuration using Custom Resources.
				The Operator watches for new API Objects and tells the target DBMS to trigger a certain stored procedure based on the content of the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup Logger
		ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
		// Load configuration
		setupLog.Info("loading config...")
		LoadConfig()
		setupLog.Info("config loaded: " + viper.ConfigFileUsed())
		setupLog.Info("registering endpoints...")
		RegisterEndpoints()
		setupLog.Info("endpoints registered")

		// Finally start the Operator
		StartOperator()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// StartOperator starts the operator by creating a new manager which in turn starts the operator controller.
func StartOperator() {
	var metricsAddr string
	var enableLeaderElection bool

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "bfa62c96.bedag.ch",
	})

	if err != nil {
		fatalError(err, "unable to start manager")
	}

	if err = (&controllers.KubernetesDbaasReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("KubernetesDbaas"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		fatalErrorWithValues(err, "unable to create controller", "controller", "KubernetesDbaas")
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fatalError(err, "problem running manager")
	}
}

// LoadConfig attempts to load the operator configuration.
//
// See config.ReadOperatorConfig for details.
func LoadConfig() {
	// If CLI flag was set
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		setupLog.Info("setting config file location from CLI flag: " + cfgFile)
	} else {
		viper.SetConfigType(DefaultConfigType)
		viper.SetConfigName(DefaultConfigFilename)
		viper.AddConfigPath(".") // search for config file in the same location as the executable file
		viper.AddConfigPath(DefaultConfigFilepath)
	}

	if err := viper.ReadInConfig(); err != nil {
		fatalError(err, ConfigLoadError)
	}

	// Parse config file
	err := config.ReadOperatorConfig(viper.GetViper())
	if err != nil {
		fatalError(err, ConfigParseError)
	}
}

// RegisterEndpoints attempts to register the endpoints specified in the operator configuration loaded from LoadConfig.
//
// See pool.Register for details.
func RegisterEndpoints() {
	for _, dbmsConfigEntry := range config.GetDbmsConfig() {
		if err := pool.Register(dbmsConfigEntry); err != nil {
			fatalError(err, DbmsConnOpenError)
		}
	}
}

func init() {
	// TODO: Metrics
	//rootCmd.Flags().StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	//_ = viper.BindPFlag("port", rootCmd.Flags().Lookup("metrics-addr"))

	rootCmd.Flags().StringVar(&cfgFile, "load-config", "", "Loads the config file from path")

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(dbaasv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func fatalError(err error, msg string) {
	setupLog.Error(err, msg)
	os.Exit(1)
}

func fatalErrorWithValues(err error, msg string, values ...string) {
	setupLog.Error(err, msg, values)
	os.Exit(1)
}
