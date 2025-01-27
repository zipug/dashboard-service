package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
)

func (d *DashboardService) SendOTP(email string) error {
	ctx := context.Background()
	user, err := d.user.GetUserByEmail(ctx, email)
	if err != nil {
		d.log.Log("error", "error occured while getting user by email", logger.WithStrAttr("email", email), logger.WithErrAttr(err))
		return err
	}
	if err := d.otp.SendOTP(ctx, int64(user.Id), email, string(user.Name)); err != nil {
		d.log.Log("error", "error occured while sending otp code to email", logger.WithErrAttr(err))
		return err
	}
	d.log.Log(
		"info",
		"otp code successfully sended to user email",
		logger.WithStrAttr("email", email),
	)
	return nil
}

func (d *DashboardService) VerifyUser(user_id int64, code models.OTPCode) error {
	ctx := context.Background()
	if err := d.otp.VerifyOTP(ctx, user_id, code); err != nil {
		d.log.Log("error", "error occured while verifying otp code", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "otp code successfully verified")
	if err := d.user.VerifyUser(ctx, models.Id(user_id)); err != nil {
		d.log.Log("error", "error occured while verifying user", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "user successfully verified")
	return nil
}
