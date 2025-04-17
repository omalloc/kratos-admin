package data

import (
	"context"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/omalloc/contrib/kratos/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/omalloc/kratos-admin/internal/biz"
	"github.com/omalloc/kratos-admin/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	wire.Bind(new(orm.DataSourceManager), new(*Data)),
	orm.NewTransactionManager,

	// rbac modules.
	NewUserRepo,
	NewRoleRepo,
	NewPermissionRepo,
)

var (
	emptyCallback = func() {}
)

type DriverDialector func(dsn string) gorm.Dialector

var driverSupported = map[string]DriverDialector{
	"sqlite": sqlite.Open,
	"mysql":  mysql.Open,
}

// Data .
type Data struct {
	db *gorm.DB
	// wrapped database client with transaction manager
	txm orm.Transaction
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	log.Infof("begin connection database with %s", c.Database.Driver)

	driver, ok := driverSupported[c.Database.Driver]
	if !ok {
		driver = sqlite.Open
	}

	db, err := orm.New(
		orm.WithDriver(driver(c.Database.Source)),
		orm.WithTracingOpts(orm.WithDatabaseName("kratos_cp")),
		orm.WithLogger(
			// orm.WithDebug(),
			orm.WIthSlowThreshold(time.Second*2),
			orm.WithSkipCallerLookup(false),
			orm.WithSkipErrRecordNotFound(true),
			orm.WithLogHelper(logger),
		),
	)

	if err != nil {
		return nil, emptyCallback, err
	}

	_ = db.Session(&gorm.Session{SkipHooks: true}).
		AutoMigrate(
			&biz.User{},
			&biz.Role{},
			&biz.Permission{},
			&biz.RolePermission{},
			&biz.UserNamespace{},
			&biz.UserRole{},
		)

	cleanup := func() {
		log.Info("closing the data resources")
	}

	data := &Data{
		db: db,
	}

	return data, cleanup, nil
}

func (d *Data) Check(ctx context.Context) error {
	// check database connection
	// but skip gorm hooks (tracing and more...)
	tx := d.db.Session(&gorm.Session{SkipHooks: true})
	if c, err := tx.DB(); err == nil {
		return c.PingContext(ctx)
	}
	return nil
}

func (d *Data) GetDataSource() *gorm.DB {
	return d.db
}
