package gormx

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Config struct {
	DBConfig        `mapstructure:",squash"`
	Slaves          []DBConfig `mapstructure:"slaves"`
	MaxIdleConns    int        `mapstructure:"max-idle-conns"`
	MaxOpenConns    int        `mapstructure:"max-open-conns"`
	ConnMaxLifetime int        `mapstructure:"conn-max-lifetime"`
	Log             LogConfig  `mapstructure:"log"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

func NewDB(zapLogger *zap.Logger, config Config) (*gorm.DB, func(), error) {
	logger := newLogger(zapLogger, config.Log)
	dsn := buildDSN(config.Host, config.Port, config.Username, config.Password, config.Database, config.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve database connection: %v", err)
	}

	if len(config.Slaves) > 0 {
		replicas := make([]gorm.Dialector, 0, len(config.Slaves))
		for i := range config.Slaves {
			s := &config.Slaves[i]
			mergeMasterSlaveConfig(config, s)
			dsn := buildDSN(s.Host, s.Port, s.Username, s.Password, s.Database, s.Charset)
			replicas = append(replicas, mysql.Open(dsn))
		}

		resolverConfig := dbresolver.Config{
			Replicas:          replicas,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}

		err = db.Use(dbresolver.Register(resolverConfig))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to use dbresolver plugin: %v", err)
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

func buildDSN(host, port, username, password, database, charset string) string {
	dsn := "%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local"
	return fmt.Sprintf(dsn, username, password, host, port, database, charset)
}

func mergeMasterSlaveConfig(masterConfig Config, slaveConfig *DBConfig) {
	if slaveConfig.Port == "" {
		slaveConfig.Port = masterConfig.Port
	}
	if slaveConfig.Username == "" {
		slaveConfig.Username = masterConfig.Username
	}
	if slaveConfig.Password == "" {
		slaveConfig.Password = masterConfig.Password
	}
	if slaveConfig.Database == "" {
		slaveConfig.Database = masterConfig.Database
	}
	if slaveConfig.Charset == "" {
		slaveConfig.Charset = masterConfig.Charset
	}
}
