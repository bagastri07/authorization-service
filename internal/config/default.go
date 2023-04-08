package config

import "time"

const (
	DefaultDatabaseMaxIdleConns    = 3
	DefaultDatabaseMaxOpenConns    = 5
	DefaultDatabaseConnMaxLifetime = 1 * time.Hour
	DefaultDatabasePingInterval    = 1 * time.Second
	DefaultDatabaseRetryAttempts   = 3

	DefaultGracefulShutdownTimeOut = 30 * time.Second

	DefaultMetricsReadTimeOut  = 5 * time.Second
	DefaultMetricsWriteTimeout = 10 * time.Second

	DefaultRedisDialTimeout  = 5 * time.Second
	DefaultRedisWriteTimeout = 2 * time.Second
	DefaultRedisReadTimeout  = 2 * time.Second
	DefaultRedisCacheTTL     = 15 * time.Minute
)
