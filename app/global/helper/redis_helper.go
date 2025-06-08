package helper

import (
	"fmt"
	"linkshare/app/constants"
)

func GenRedisKeyAccessTokenSessionByUserID(userId int64) string {
	return fmt.Sprintf("%s:%d", constants.ACCESS_TOKEN_SESSIONS_COL, userId)
}

func GenRedisKeyRefreshTokenSessionByUserID(userId int64) string {
	return fmt.Sprintf("%s:%d", constants.REFRESH_TOKEN_SESSIONS_COL, userId)
}
