package main

import (
	"context"
	"fmt"
	"golang.org/x/xerrors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type (
	Repository interface {
		GenerateDSNLocal(name string, username string, password string, ip string, port int64) string
		DB() *gorm.DB
		Open(ctx context.Context, dsn string) (*gorm.DB, error)
	}

	repository struct {
		dsn string
		db  *gorm.DB
		ctx context.Context
	}
)

func NewRepository(ctx context.Context, dsn string) Repository {
	r := &repository{
		dsn: dsn,
		ctx: ctx,
	}

	//r.dsn = r.GenerateDSNLocal(
	//	name, userName, password, ip, port)

	db, err := r.Open(ctx, r.dsn)
	if err != nil {
		log.Fatalf(":%+v\n", err)
	}

	// Set DB
	r.db = db

	return r
}

func (r *repository) GenerateDSNLocal(name string, username string, password string, ip string, port int64) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, ip, port, name)
}

func (r *repository) DB() *gorm.DB {
	return r.db
}

func (r *repository) Open(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Caches Prepared Statement
		// https://gorm.io/docs/performance.html#Caches-Prepared-Statement
		PrepareStmt: false,
	})

	// Set Context
	db.WithContext(ctx)

	if err != nil {
		log.Printf("error connect to db %+v", xerrors.Errorf(": %w", err))
		return nil, xerrors.Errorf(": %w", err)
	}

	// GetDomains generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("error to fetch sqlDB %+v", xerrors.Errorf(": %w", err))
		return nil, xerrors.Errorf(": %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(100)

	// SetMaxOpenConns sets the maximum number of Open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
