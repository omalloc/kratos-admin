package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/omalloc/contrib/kratos/health"
	trace "github.com/omalloc/contrib/kratos/tracing"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/omalloc/kratos-admin/internal/conf"
	"github.com/omalloc/kratos-admin/internal/server"
	_ "github.com/omalloc/kratos-admin/pkg/gorm-schema"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// To render a whole-file example, we need a package-level declaration.
	_ = ""
	// Name is the name of the compiled software.
	Name string = "kratos-admin"
	// Version is the version of the compiled software.
	Version string = "v0.0.0"
	// GitHash is the git-hash of the compiled software.
	GitHash string = "-"
	// Built is build-time the compiled software.
	Built string = "0"
	// flagconf is the config flag.
	flagconf    string
	flagverbose bool

	id, _ = os.Hostname()
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(nil))

	json.MarshalOptions.UseProtoNames = true
	json.MarshalOptions.UseEnumNumbers = true

	rootCmd.PersistentFlags().StringVar(&flagconf, "conf", "../../configs", "config path")
	rootCmd.PersistentFlags().BoolVarP(&flagverbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(versionCmd)
}

func newApp(logger log.Logger, embedEtcd *server.EmbedEtcdServer, registrar registry.Registrar, gs *grpc.Server, hs *http.Server, hh *health.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{
			"hang": "true",
		}),
		kratos.Logger(logger),
		kratos.Registrar(registrar),
		kratos.Server(
			gs,
			hs,
			hh,
			embedEtcd,
		),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

	logger := log.With(log.NewFilter(log.DefaultLogger, log.FilterLevel(log.LevelDebug), log.FilterFunc(filter)),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	log.SetLogger(logger)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	_ = c.Watch("logger.level", func(s string, value config.Value) {
		lvl := log.ParseLevel(value.Load().(string))
		log.Infof("log level has changed %v", lvl)
		//l := log.With(log.NewFilter(logger, log.FilterLevel(lvl)))
	})

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	trace.InitTracer(
		trace.WithServiceName(Name), // service-name registered in jaeger-service
		trace.WithEndpoint(bc.Tracing.GetEndpoint()),
	)

	app, cleanup, err := wireApp(&bc, bc.Server, bc.Data, bc.Passport, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

const fuzzyStr = "***"

var re = regexp.MustCompile(`password:"([^"]+)"`)

func filter(level log.Level, keyvals ...any) bool {
	if len(keyvals) == 0 {
		return false
	}

	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == "args" && strings.Contains(keyvals[i+1].(string), "password") {
			keyvals[i+1] = re.ReplaceAllString(keyvals[i+1].(string), fuzzyStr)
		}
	}
	return false
}
