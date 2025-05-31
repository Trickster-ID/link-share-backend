package dto

import (
	"linkshare/app/global/model"
	"linkshare/app/models"
)

type GetByRefreshTokenChan struct {
	Data   *models.RefreshTokenSession
	ErrLog *model.ErrorLog
}
