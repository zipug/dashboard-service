package application

import (
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	"dashboard/internal/core/service/articles"
	"dashboard/internal/core/service/attachments"
	"dashboard/internal/core/service/bots"
	generatedreports "dashboard/internal/core/service/generated_reports"
	files "dashboard/internal/core/service/minio"
	"dashboard/internal/core/service/otp"
	"dashboard/internal/core/service/projects"
	"dashboard/internal/core/service/reports"
	"dashboard/internal/core/service/roles"
	"dashboard/internal/core/service/user"
	"dashboard/internal/infrastructure/repository/minio"
	"dashboard/internal/infrastructure/repository/postgres"
	"dashboard/internal/infrastructure/repository/redis"
)

var (
	configCommonService         = config.NewConfigService()
	postgresRepository          = postgres.NewPostgresRepository(configCommonService)
	redisRepository             = redis.NewRedisRepository(configCommonService)
	minioRepository             = minio.NewMinioRepository(configCommonService)
	usersCoreService            = user.NewUserService(postgresRepository)
	otpCoreService              = otp.NewOTPService(configCommonService, redisRepository)
	minioCoreService            = files.NewMinioService(minioRepository)
	rolesCoreService            = roles.NewRolesService(postgresRepository)
	projectsCoreService         = projects.NewProjectsService(postgresRepository)
	articlesCoreService         = articles.NewArticlesService(postgresRepository)
	attachmentsCoreService      = attachments.NewAttachmentsService(postgresRepository)
	botsCoreService             = bots.NewBotsService(postgresRepository)
	reportsCoreService          = reports.NewReportsService(postgresRepository, redisRepository)
	generatedReportsCoreService = generatedreports.NewGeneratedReportsService(postgresRepository)
	authModule                  = auth.New(configCommonService)
)

var DashboardAppService = NewDashboardService(
	configCommonService,
	usersCoreService,
	authModule,
	otpCoreService,
	minioCoreService,
	rolesCoreService,
	projectsCoreService,
	articlesCoreService,
	attachmentsCoreService,
	botsCoreService,
	reportsCoreService,
	generatedReportsCoreService,
)
