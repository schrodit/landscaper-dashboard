module github.com/schrodit/landscaper-dashboard

go 1.16

require (
	github.com/gardener/landscaper/apis v0.7.0
	github.com/gin-gonic/gin v1.7.1
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.2.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.15.0
	gopkg.in/olahol/melody.v1 v1.0.0-20170518105555-d52139073376
	k8s.io/api v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/yaml v1.2.0
)
