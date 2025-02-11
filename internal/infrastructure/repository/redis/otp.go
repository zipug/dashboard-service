package redis

import (
	"context"
	"dashboard/internal/core/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func (repo *RedisRepository) SaveUserOTP(ctx context.Context, user_id int64, email, username string, otp models.OTPCode) error {
	key := strconv.FormatInt(user_id, 10)
	if key == "" {
		return fmt.Errorf("%w: user_id: %d", ErrKeyIsNotValid, user_id)
	}
	code, err := otp.ToStr()
	if err != nil {
		return nil
	}
	user_id_str := strconv.FormatInt(user_id, 10)
	evt := models.Event{
		Type:      "otp",
		Timestamp: time.Now().Unix(),
		Payload: models.OTPPayload{
			Type:     models.Verify,
			UserID:   user_id,
			UserName: username,
			Email:    email,
			Code:     code,
		},
	}
	message, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	if err := repo.rdb.Set(ctx, user_id_str, code, repo.expire).Err(); err != nil {
		return err
	}
	if err := repo.rdb.Publish(ctx, "otp", message).Err(); err != nil {
		return err
	}
	return nil
}

func (repo *RedisRepository) GetUserOTP(ctx context.Context, user_id int64) (models.OTPCode, error) {
	key := strconv.FormatInt(user_id, 10)
	if key == "" {
		return models.OTPCode(0), fmt.Errorf("%w: user_id: %d", ErrKeyIsNotValid, user_id)
	}
	val, err := repo.rdb.Get(ctx, key).Result()
	if err != nil {
		return models.OTPCode(0), err
	}
	otp, err := models.StrToOTPCode(val)
	if err != nil {
		return models.OTPCode(0), err
	}
	return models.OTPCode(otp), nil
}
