package connection

import (
	"errors"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"coresense/pkg/common/config"
	"coresense/pkg/common/utils/database"
)

type lowCaseNamingStrategy struct {
	schema.NamingStrategy
}

func GetConnection(databaseConfig config.Database, logger zerolog.Logger) (*gorm.DB, error) {
	if !databaseConfig.HasURL() {
		return nil, errors.New("database connection string is empty")
	}

	db, err := gorm.Open(newDialector(databaseConfig.URL), &gorm.Config{
		Logger:         database.NewLogger(logger, databaseConfig.Logger),
		NamingStrategy: lowCaseNamingStrategy{},
		PrepareStmt:    true,
		TranslateError: true,
	})

	if err != nil {
		return nil, errors.Join(errors.New("failed to connect to database"), err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Join(errors.New("failed to get database connection"), err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, errors.Join(errors.New("failed to ping database"), err)
	}

	sqlDB.SetMaxOpenConns(databaseConfig.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(databaseConfig.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(databaseConfig.ConnectionMaxLifeTime)

	return db, nil
}
