package application

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
	"fmt"
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
	log   *logger.Logger
	state State
}

func NewDashboardService(cfg *config.AppConfig, user ports.UserService, auth *auth.Auth) *DashboardService {
	d := &DashboardService{
		cfg:  cfg,
		user: user,
		auth: auth,
	}

	d.state = Created

	return d
}

func (d *DashboardService) GetUserById(id int64) (models.User, error) {
	d.log.Log("info", fmt.Sprintf("getting user by id: %d", id))
	ctx := context.Background()
	usr, err := d.user.GetUserById(ctx, id)
	if err != nil {
		d.log.Log("error", fmt.Sprintf("error while getting user by id: %d", id), logger.WithErrAttr(err))
		return usr, err
	}
	d.log.Log("info", "user successfully get", logger.WithInt64Attr("user_id", id))
	return usr, nil
}

func (d *DashboardService) RegisterUser(user dto.UserDto) (int64, error) {
	d.log.Log("info", "registering user")
	ctx := context.Background()
	id, err := d.user.RegisterUser(ctx, user.ToValue())
	if err != nil {
		d.log.Log("error", "error while registering user", logger.WithErrAttr(err))
		return int64(dto.BadUserId), err
	}
	d.log.Log("info", "user successfully registered", logger.WithInt64Attr("user_id", id))
	return id, nil
}

func (d *DashboardService) LoginUser(user dto.UserDto) (models.User, error) {
	d.log.Log("info", "login user")
	ctx := context.Background()
	usr, err := d.user.LoginUser(ctx, user.ToValue())
	if err != nil {
		d.log.Log("error", "error while loginnig user", logger.WithErrAttr(err))
		return usr, err
	}
	d.log.Log("info", "user successfully logined", logger.WithInt64Attr("user_id", int64(usr.Id)))
	return usr, nil
}

func (d *DashboardService) GetAllUsers() ([]models.User, error) {
	d.log.Log("info", "get all users")
	ctx := context.Background()
	users, err := d.user.GetAllUsers(ctx)
	if err != nil {
		d.log.Log("error", "error while fetching users", logger.WithErrAttr(err))
		return users, err
	}
	d.log.Log("info", "users successfully fetched", logger.WithInt64Attr("user_count", int64(len(users))))
	return users, nil
}

func (d *DashboardService) SaveUser(user models.User) {
	fmt.Println("saving user")
	ctx := context.Background()
	if err := d.user.SaveUser(ctx, user); err != nil {
		fmt.Printf("error occured while saving user: %s\n", err.Error())
	}
	fmt.Println("user successfully saved")
}

func (d *DashboardService) DeleteUser(id int64) {
	fmt.Println("deleting user")
	ctx := context.Background()
	if err := d.user.DeleteUser(ctx, id); err != nil {
		fmt.Printf("error occured while deleting user with id: %d, %s\n", id, err.Error())
	}
	fmt.Println("user successfully deleted")
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
	d.log = logger.New()
	d.state = Ready
}

func (d *DashboardService) Restart() {
	d.Run()
}

func (d *DashboardService) Stop() {
	d.state = Down
}
