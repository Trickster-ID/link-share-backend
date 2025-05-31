package constants

import "time"

const (
	ACCESS_TOKEN_EXPIRED  = time.Hour * 12
	REFRESH_TOKEN_EXPIRED = (time.Hour * 24) * 7
)
