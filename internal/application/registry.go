package application

import (
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	"dashboard/internal/core/service/otp"
	"dashboard/internal/core/service/user"
	"dashboard/internal/infrastructure/repository/postgres"
	"dashboard/internal/infrastructure/repository/redis"
)

var (
	configCommonService = config.NewConfigService()
	userCoreService     = user.NewUserService(postgres.NewPostgresRepository(configCommonService))
	otpCoreService      = otp.NewOTPService(configCommonService, redis.NewRedisRepository(configCommonService))
	authModule          = auth.New(configCommonService)
)

var DashboardAppService = NewDashboardService(configCommonService, userCoreService, authModule, otpCoreService)
