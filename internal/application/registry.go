package application

import (
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	"dashboard/internal/core/service/otp"
	"dashboard/internal/core/service/roles"
	"dashboard/internal/core/service/user"
	"dashboard/internal/infrastructure/repository/postgres"
	"dashboard/internal/infrastructure/repository/redis"
)

var (
	configCommonService = config.NewConfigService()
	postgresRepository  = postgres.NewPostgresRepository(configCommonService)
	userCoreService     = user.NewUserService(postgresRepository)
	otpCoreService      = otp.NewOTPService(configCommonService, redis.NewRedisRepository(configCommonService))
	rolesCoreService    = roles.NewRolesService(postgresRepository)
	authModule          = auth.New(configCommonService)
)

var DashboardAppService = NewDashboardService(configCommonService, userCoreService, authModule, otpCoreService, rolesCoreService)
