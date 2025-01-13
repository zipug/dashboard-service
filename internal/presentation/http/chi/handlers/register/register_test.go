package register

import (
	"bytes"
	"dashboard/internal/common/service/auth"
	"dashboard/internal/common/service/config"
	"dashboard/internal/presentation/http/chi/handlers"
	"dashboard/internal/presentation/http/chi/handlers/register/mocks"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

type Payload struct {
	Email          string `json:"email,omitempty"`
	Password       string `json:"password,omitempty"`
	RepeatPassword string `json:"repeat_password,omitempty"`
	Name           string `json:"name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
}

// TestRegisterUser function 
func TestRegisterUser(t *testing.T) {
	payload := Payload{
		Email:          "test@test.com",
		Password:       "password",
		RepeatPassword: "password",
		Name:           "Maksim",
		LastName:       "Volodin",
	}
	req, err := json.Marshal(payload)
	config := config.NewConfigService()
	config.JwtSecretKey = "supersecret"
	config.AccessTokenExpiration = 15 * time.Minute
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(req)

	ctrl := gomock.NewController(t)

	app := mocks.NewMockDashboardService(ctrl)
	app.EXPECT().RegisterUser(gomock.Any()).Return(int64(1), nil)
	log := mocks.NewMockLogger(ctrl)
	auth := auth.New(config)

	server := httptest.NewServer(RegisterUser(app, log, auth, config.AccessTokenExpiration))
	resp, err := http.Post(server.URL, "application/json", r)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	var res handlers.Response
	if err := json.Unmarshal(b, &res); err != nil {
		t.Error(err)
	}
	if res.Status != handlers.Success {
		t.Errorf("expected status success, got %s", res.Status)
	}

	data, ok := res.Data.(map[string]interface{})
	if !ok {
		t.Error("expected data to be map[string]interface{}")
	}
	dataId, ok := data["id"].(float64)
	if !ok {
		t.Error("expected id to be float64")
	}
	if dataId != 1 {
		t.Errorf("expected id 1, got %d", int(dataId))
	}
}
