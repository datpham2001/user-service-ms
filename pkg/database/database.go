package database

import (
	"fmt"
	"log"

	"github.com/datpham/user-service-ms/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := getDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.Postgres.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.Postgres.MaxOpenConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return &Database{DB: db}, nil
}

func getDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.DBName,
		cfg.Database.Postgres.SSLMode,
	)
}

func (d *Database) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %v", err)
		}

		return sqlDB.Close()
	}

	return nil
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	if d.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	return d.DB.AutoMigrate(models...)
}
