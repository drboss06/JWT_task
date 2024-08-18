package JWTServiceObjects

import "time"

type Session struct {
	Guid         string    `bson:"guid"`
	RefreshToken []byte    `bson:"refresh_token"`
	LiveTime     time.Time `bson:"live_time"`
	ClientIp     string    `bson:"client_ip"`
}
