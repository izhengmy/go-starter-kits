package gormx

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Config struct {
	Sources []DataSourceConfig `mapstructure:"data-sources"`
	Log     LogConfig          `mapstructure:"log"`
}

type DataSourceConfig struct {
	ID              string   `mapstructure:"id"`
	Driver          string   `mapstructure:"driver"`
	DSN             string   `mapstructure:"dsn"`
	MaxIdleConns    int      `mapstructure:"max-idle-conns"`
	MaxOpenConns    int      `mapstructure:"max-open-conns"`
	ConnMaxLifetime int      `mapstructure:"conn-max-lifetime"`
	Slaves          []string `mapstructure:"slaves"`
}

type DataSources map[string]*gorm.DB

func InitDataSources(zapLogger *zap.Logger, config Config) (DataSources, func(), error) {
	dataSources := DataSources{}
	cleanups := make([]func(), 0, len(config.Sources))

	for _, cfg := range config.Sources {
		logger := newLogger(zapLogger, config.Log, cfg.ID)
		db, cleanup, err := CreateDB(logger, cfg)
		if err != nil {
			return nil, nil, err
		}
		dataSources[cfg.ID] = db
		cleanups = append(cleanups, cleanup)
	}

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
	}

	return dataSources, cleanup, nil
}

func CreateDB(logger *logger, config DataSourceConfig) (*gorm.DB, func(), error) {
	driver := config.Driver
	dsn := config.DSN

	dialector, err := createDialector(driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create dialector for [%s]: %w", config.ID, err)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger,
	})

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve database connection for [%s]: %w", config.ID, err)
	}

	if len(config.Slaves) > 0 {
		replicas := make([]gorm.Dialector, 0, len(config.Slaves))

		for i := range config.Slaves {
			dsn := config.Slaves[i]
			dialector, err := createDialector(driver, dsn)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create replica dialector for [%s]: %w", config.ID, err)
			}
			replicas = append(replicas, dialector)
		}

		resolverConfig := dbresolver.Config{
			Replicas:          replicas,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}

		err = db.Use(dbresolver.Register(resolverConfig))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to use dbresolver plugin for [%s]: %w", config.ID, err)
		}
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Minute)

	cleanup := func() {
		_ = sqlDB.Close()
	}

	return db, cleanup, nil
}

func createDialector(driver string, dsn string) (gorm.Dialector, error) {
	switch driver {
	case "mysql":
		return mysql.Open(dsn), nil
	case "postgres":
	case "pgsql":
		return postgres.Open(dsn), nil
	case "sqlite":
		return sqlite.Open(dsn), nil
	}
	return nil, fmt.Errorf("unsupported database driver %s", driver)
}
