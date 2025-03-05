package application

import (
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/ports"
)

type State int

const (
	Created State = iota + 1
	Running
	Ready
	Down
)

type DashboardService struct {
	cfg              *config.AppConfig
	user             ports.UserService
	auth             *auth.Auth
	otp              ports.OTPService
	minio            ports.MinioService
	role             ports.RolesService
	project          ports.ProjectsService
	article          ports.ArticlesService
	attachment       ports.AttachmentsService
	bot              ports.BotsService
	report           ports.ReportsService
	generated_report ports.GeneratedReportsService
	chat             ports.ChatsService
	telegram         ports.TelegramService
	log              *logger.Logger
	state            State
}

func NewDashboardService(
	cfg *config.AppConfig,
	user ports.UserService,
	auth *auth.Auth,
	otp ports.OTPService,
	minio ports.MinioService,
	role ports.RolesService,
	project ports.ProjectsService,
	article ports.ArticlesService,
	attachment ports.AttachmentsService,
	bot ports.BotsService,
	report ports.ReportsService,
	generated_report ports.GeneratedReportsService,
	chat ports.ChatsService,
	telegram ports.TelegramService,
) *DashboardService {
	d := &DashboardService{
		cfg:              cfg,
		user:             user,
		auth:             auth,
		otp:              otp,
		minio:            minio,
		role:             role,
		project:          project,
		article:          article,
		attachment:       attachment,
		bot:              bot,
		report:           report,
		generated_report: generated_report,
		chat:             chat,
		telegram:         telegram,
	}

	d.state = Created

	return d
}

func (d *DashboardService) GetConfig() *config.AppConfig {
	return d.cfg
}

func (d *DashboardService) GetLogger() *logger.Logger {
	return d.log
}

func (d *DashboardService) GetState() State {
	return d.state
}

func (d *DashboardService) GetAuth() *auth.Auth {
	return d.auth
}

func (d *DashboardService) Run() {
	d.state = Running
	d.log = logger.New(d.cfg.Env)
	d.state = Ready
}

func (d *DashboardService) Restart() {
	d.Run()
}

func (d *DashboardService) Stop() {
	d.state = Down
}
