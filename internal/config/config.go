package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// embed from build flag
	serviceName    string
	serviceVersion string
)

func parseDuration(in string, defaultDuration time.Duration) time.Duration {
	dur, err := time.ParseDuration(in)
	if err != nil {
		return defaultDuration
	}
	return dur
}

func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.SetConfigName("config")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Warningf("%v", err)
	}
}

func ServiceName() string {
	return serviceName
}

func ServiceVersion() string {
	return serviceVersion
}

func Env() string {
	return viper.GetString("env")
}

func LogLevel() string {
	return viper.GetString("log_level")
}

func GRPCPort() string {
	return viper.GetString("ports.grpc")
}

func MetricsPort() string {
	return viper.GetString("metrics.port")
}

func MetricsReadTimeout() time.Duration {
	mrt := viper.GetString("metrics.read_timeout")
	return parseDuration(mrt, DefaultMetricsReadTimeOut)
}

func MetricsWriteTimeout() time.Duration {
	mrt := viper.GetString("metrics.read_timeout")
	return parseDuration(mrt, DefaultMetricsWriteTimeout)
}

func DatabaseHost() string {
	return viper.GetString("postgres.host")
}

func DatabaseName() string {
	return viper.GetString("postgres.database")
}

func DatabaseUsername() string {
	return viper.GetString("postgres.username")
}

func DatabasePassword() string {
	return viper.GetString("postgres.password")
}

func DatabaseSSLMode() string {
	if viper.IsSet("postgres.sslmode") {
		return viper.GetString("postgres.sslmode")
	}
	return "disable"
}

func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		DatabaseUsername(),
		DatabasePassword(),
		DatabaseHost(),
		DatabaseName(),
		DatabaseSSLMode())
}

func DatabasePingInterval() time.Duration {
	if viper.GetInt("postgres.ping_interval") <= 0 {
		return DefaultDatabasePingInterval
	}
	return time.Duration(viper.GetInt("postgres.ping_interval")) * time.Millisecond
}

func DatabaseRetryAttempts() float64 {
	if viper.GetInt("postgres.retry_attempts") > 0 {
		return float64(viper.GetInt("postgres.retry_attempts"))
	}
	return DefaultDatabaseRetryAttempts
}

func DatabaseMaxIdleConns() int {
	if viper.GetInt("postgres.max_idle_conns") <= 0 {
		return DefaultDatabaseMaxIdleConns
	}
	return viper.GetInt("postgres.max_idle_conns")
}

func DatabaseMaxOpenConns() int {
	if viper.GetInt("postgres.max_open_conns") <= 0 {
		return DefaultDatabaseMaxOpenConns
	}
	return viper.GetInt("postgres.max_open_conns")
}

func DatabaseConnMaxLifetime() time.Duration {
	if !viper.IsSet("postgres.conn_max_lifetime") {
		return DefaultDatabaseConnMaxLifetime
	}
	return time.Duration(viper.GetInt("postgres.conn_max_lifetime")) * time.Millisecond
}

func GracefulShutdownTimeOut() time.Duration {
	cfg := viper.GetString("graceful_shutdown_timeout")
	return parseDuration(cfg, DefaultGracefulShutdownTimeOut)
}

func RedisCacheHost() string {
	return viper.GetString("redis.cache_host")
}

func RedisDialTimeout() time.Duration {
	cfg := viper.GetString("redis.dial_timeout")
	return parseDuration(cfg, DefaultRedisDialTimeout)
}

func RedisWriteTimeout() time.Duration {
	cfg := viper.GetString("redis.write_timeout")
	return parseDuration(cfg, DefaultRedisWriteTimeout)
}

func RedisReadTimeout() time.Duration {
	cfg := viper.GetString("redis.read_timeout")
	return parseDuration(cfg, DefaultRedisReadTimeout)
}

func RedisCacheTTL() time.Duration {
	cfg := viper.GetString("cache_ttl")
	return parseDuration(cfg, DefaultRedisCacheTTL)
}

func JwtAccessTokenSecret() []byte {
	return []byte(viper.GetString("jwt.access_token_secret"))
}

func JwtRefreshTokenSecret() []byte {
	return []byte(viper.GetString("jwt.refresh-token-secret"))
}

func CryptoCost() int {
	if !viper.IsSet("crypto.cost") {
		return DefaultCryptoCost
	}
	return viper.GetInt("cyrpto.cost")
}

func AccessTokenDuration() time.Duration {
	cfg := viper.GetString("jwt.access_token_duration")
	return parseDuration(cfg, DefaultAccessTokenDuration)
}

func RefreshTokenDuration() time.Duration {
	cfg := viper.GetString("jwt.refresh_token_duration")
	return parseDuration(cfg, DefaultRefreshTokenDuration)
}
