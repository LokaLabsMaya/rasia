package bootstrap

import (
	"github.com/sirupsen/logrus"

	"github.com/redhajuanda/komon/logger"
	"github.com/lokalabsmaya/rasia/configs"
	"github.com/lokalabsmaya/rasia/internal/adapter/inbound/http"
	httpHandler "github.com/lokalabsmaya/rasia/internal/adapter/inbound/http/handler"
	"github.com/lokalabsmaya/rasia/internal/adapter/outbound/sqlite"
	"github.com/lokalabsmaya/rasia/internal/core/port/outbound"
	"github.com/lokalabsmaya/rasia/internal/core/service/export"
	"github.com/lokalabsmaya/rasia/internal/core/service/filecontent"
	"github.com/lokalabsmaya/rasia/internal/core/service/namespace"
	"github.com/lokalabsmaya/rasia/internal/core/service/secret"
	"github.com/lokalabsmaya/rasia/internal/core/service/secretfile"
)

// There are 4 types of resources:
// - Resource[T] is a generic resource that can be initialized and retrieved
// - ResourceRunnable[T] is a resource that can be run, this resource should implement `OnStart(ctx context.Context) error` and `OnStop(ctx context.Context) error` methods
// - ResourceExecutable[T] is a resource that can be executed, this resource should implement `Execute(ctx context.Context) error` method
// - ResourceClosable[T] is a resource that can be closed, this resource should implement `Close() error` method
type Dependency struct {
	cfgFile string

	cfg                Resource[*configs.Config]
	log                Resource[logger.Logger]
	secretsRepo        Resource[outbound.Repository]
	serviceNamespace   Resource[*namespace.Service]
	serviceSecretFile  Resource[*secretfile.Service]
	serviceSecret      Resource[*secret.Service]
	serviceFileContent Resource[*filecontent.Service]
	serviceExport      Resource[*export.Service]
	httpHandlers       Resource[[]http.Handler]

	sqliteDB ResourceClosable[*sqlite.DB]

	httpRunner ResourceRunnable[*http.HTTP]
}

// NewDependency creates a new dependency instance
func NewDependency(cfgFile string) *Dependency {
	return &Dependency{
		cfgFile: cfgFile,
	}
}

// GetConfig resolves and returns the config dependency
func (d *Dependency) GetConfig() *configs.Config {
	return d.cfg.Resolve(func() *configs.Config {
		return configs.LoadConfig(d.cfgFile)
	})
}

// GetLogger resolves and returns the logger dependency
func (d *Dependency) GetLogger() logger.Logger {
	return d.log.Resolve(func() logger.Logger {
		cfg := d.GetConfig()
		log := logger.New(cfg.App.Name, logger.Options{
			RedactedFields: cfg.Log.RedactedFields,
		})
		logger.SetLevel(logrus.Level(cfg.Log.Level))
		if cfg.Log.Format == "json" {
			logger.SetFormatter(&logrus.JSONFormatter{
				PrettyPrint: false,
			})
		} else {
			logger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp: true,
				ForceColors:   true,
			})
		}
		return log.WithParam("service", cfg.App.Name)
	})
}

// GetSQLiteDB resolves and returns the SQLite database dependency
func (d *Dependency) GetSQLiteDB() *sqlite.DB {
	return d.sqliteDB.Resolve(func() *sqlite.DB {
		cfg := d.GetConfig()
		return sqlite.NewDB(cfg.Secrets.DBPath)
	})
}

// GetSecretsRepository resolves and returns the secrets SQLite repository dependency
func (d *Dependency) GetSecretsRepository() outbound.Repository {
	return d.secretsRepo.Resolve(func() outbound.Repository {
		return sqlite.NewRepository(d.GetSQLiteDB())
	})
}

// GetServiceNamespace resolves and returns the namespace service dependency
func (d *Dependency) GetServiceNamespace(repo outbound.Repository) *namespace.Service {
	return d.serviceNamespace.Resolve(func() *namespace.Service {
		return namespace.NewService(d.GetConfig(), d.GetLogger(), repo)
	})
}

// GetServiceSecretFile resolves and returns the secret file service dependency
func (d *Dependency) GetServiceSecretFile(repo outbound.Repository) *secretfile.Service {
	return d.serviceSecretFile.Resolve(func() *secretfile.Service {
		return secretfile.NewService(d.GetConfig(), d.GetLogger(), repo)
	})
}

// GetServiceSecret resolves and returns the secret service dependency
func (d *Dependency) GetServiceSecret(repo outbound.Repository) *secret.Service {
	return d.serviceSecret.Resolve(func() *secret.Service {
		return secret.NewService(d.GetConfig(), d.GetLogger(), repo)
	})
}

// GetServiceFileContent resolves and returns the file content service dependency
func (d *Dependency) GetServiceFileContent(repo outbound.Repository) *filecontent.Service {
	return d.serviceFileContent.Resolve(func() *filecontent.Service {
		return filecontent.NewService(d.GetConfig(), d.GetLogger(), repo)
	})
}

// GetServiceExport resolves and returns the export service dependency
func (d *Dependency) GetServiceExport(repo outbound.Repository) *export.Service {
	return d.serviceExport.Resolve(func() *export.Service {
		return export.NewService(d.GetConfig(), d.GetLogger(), repo)
	})
}

// GetHTTPHandlers resolves and returns the http handlers dependency
func (d *Dependency) GetHTTPHandlers() []http.Handler {
	return d.httpHandlers.Resolve(func() []http.Handler {
		repo := d.GetSecretsRepository()
		return []http.Handler{
			httpHandler.NewNamespaceHandler(d.GetConfig(), d.GetLogger(), d.GetServiceNamespace(repo)),
			httpHandler.NewSecretFileHandler(d.GetConfig(), d.GetLogger(), d.GetServiceSecretFile(repo)),
			httpHandler.NewSecretHandler(d.GetConfig(), d.GetLogger(), d.GetServiceSecret(repo)),
			httpHandler.NewFileContentHandler(d.GetConfig(), d.GetLogger(), d.GetServiceFileContent(repo), d.GetServiceExport(repo)),
		}
	})
}

// GetHTTP resolves and returns the http dependency
func (d *Dependency) GetHTTP() *http.HTTP {
	return d.httpRunner.Resolve(func() *http.HTTP {
		return http.NewHTTP(d.GetConfig(), d.GetLogger(), d.GetHTTPHandlers())
	})
}
