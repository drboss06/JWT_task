package JWTServiceObjects

import "time"

type Session struct {
	Guid         string    `db:"guid"`
	RefreshToken []byte    `db:"refresh_token"`
	LiveTime     time.Time `db:"live_time"`
	ClientIp     string    `db:"client_ip"`
}
