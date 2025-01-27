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
	cfg   *config.AppConfig
	user  ports.UserService
	auth  *auth.Auth
	otp   ports.OTPService
	role  ports.RolesService
	log   *logger.Logger
	state State
}

func NewDashboardService(
	cfg *config.AppConfig,
	user ports.UserService,
	auth *auth.Auth,
	otp ports.OTPService,
	role ports.RolesService,
) *DashboardService {
	d := &DashboardService{
		cfg:  cfg,
		user: user,
		auth: auth,
		otp:  otp,
		role: role,
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
