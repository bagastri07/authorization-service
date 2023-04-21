package constant

import "time"

const (
	JwtAccessTokenExpirationTime  = 10 * time.Minute
	JwtRefreshTokenExpirationTime = 6 * time.Hour
)
