package infrastructure

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"

	"github.com/bagastri07/authorization-service/internal/config"
	"github.com/jpillora/backoff"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// PostgreSQL represents gorm DB
	PostgreSQL *gorm.DB

	// StopTickerCh signal for closing ticker channel
	StopTickerCh chan bool
)

func InitializePostgresConn() {
	conn, err := openPostgresConn(config.DatabaseDSN())
	if err != nil {
		log.WithField("databaseDSN", config.DatabaseDSN()).Fatal("failed to connect postgresql database: ", err)
	}

	PostgreSQL = conn
	StopTickerCh = make(chan bool)

	go checkConnection(time.NewTicker(config.DatabasePingInterval()))

	PostgreSQL.Logger = newCustomLogger()

	switch config.LogLevel() {
	case "error":
		PostgreSQL.Logger = PostgreSQL.Logger.LogMode(gormLogger.Error)
	case "warn":
		PostgreSQL.Logger = PostgreSQL.Logger.LogMode(gormLogger.Warn)
	default:
		PostgreSQL.Logger = PostgreSQL.Logger.LogMode(gormLogger.Info)

	}

	log.Info("Connection to PostgreSQL Server success...")
}

func checkConnection(ticker *time.Ticker) {
	for {
		select {
		case <-StopTickerCh:
			ticker.Stop()
			return
		case <-ticker.C:
			if _, err := PostgreSQL.DB(); err != nil {
				reconnectPostgresConn()
			}
		}
	}
}

func reconnectPostgresConn() {
	b := backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    100 * time.Millisecond,
		Max:    1 * time.Second,
	}

	postgresRetryAttempts := config.DatabaseRetryAttempts()

	for b.Attempt() < postgresRetryAttempts {
		conn, err := openPostgresConn(config.DatabaseDSN())
		if err != nil {
			log.WithField("databaseDSN", config.DatabaseDSN()).Error("failed to connect postgresql database: ", err)
		}

		if conn != nil {
			PostgreSQL = conn
			break
		}
		time.Sleep(b.Duration())
	}

	if b.Attempt() >= postgresRetryAttempts {
		log.Fatal("maximum retry to connect database")
	}
	b.Reset()
}

func openPostgresConn(dsn string) (*gorm.DB, error) {
	psqlDialector := postgres.Open(dsn)
	db, err := gorm.Open(psqlDialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMaxIdleConns(config.DatabaseMaxIdleConns())
	conn.SetMaxOpenConns(config.DatabaseMaxOpenConns())
	conn.SetConnMaxLifetime(config.DatabaseConnMaxLifetime())

	return db, nil
}

type CustomLogger struct {
	gormLogger.Interface
}

func (c CustomLogger) Error(ctx context.Context, msg string, v ...interface{}) {
	if len(v) == 1 {
		if err, ok := v[0].(error); ok {
			if err == gorm.ErrRecordNotFound {
				// log a warning for "record not found" errors
				c.Interface.Warn(ctx, msg, v...)
				return
			}
		}
	}
	// log an error for other errors
	c.Interface.Error(ctx, msg, v...)
}

func (c CustomLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return CustomLogger{c.Interface.LogMode(level)}
}

func newCustomLogger() CustomLogger {
	return CustomLogger{gormLogger.Default.LogMode(gormLogger.Info)}
}
