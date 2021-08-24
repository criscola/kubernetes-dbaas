module github.com/bedag/kubernetes-dbaas

go 1.16

require (
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jackc/pgx/v4 v4.11.0
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/onsi/ginkgo v1.16.1
	github.com/onsi/gomega v1.11.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/xo/dburl v0.7.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/ratelimit v0.2.0
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20210608053332-aa57babbf139 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.1.2 // indirect
	golang.stackrox.io/kube-linter v0.0.0-20210707220328-19fa6db01f27 // indirect
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/apiserver v0.20.2
	k8s.io/client-go v0.20.4
	sigs.k8s.io/controller-runtime v0.8.3
)
