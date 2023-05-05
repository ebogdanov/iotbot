package iot

import (
	"context"
	"errors"
	"fmt"
	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/constant"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/env/extension"
	logger2 "github.com/tuya/tuya-connector-go/connector/logger"
	"strconv"
	"su27bot/internal/config"
	"su27bot/internal/model"
)

type Tuya struct {
	actions []Action
	UserId  string
}

func NewTuyaClient(apiUrl, msgUrl string, cfg *config.Tuya) *Tuya {
	logInstance := func() extension.ILogger {
		return logger2.NewDefaultLogger("tysdk", false)
	}

	extension.SetLog(constant.TUYA_LOG, logInstance)

	connector.InitWithOptions(env.WithApiHost(apiUrl),
		env.WithMsgHost(msgUrl),
		env.WithAccessID(cfg.AccessId),
		env.WithAccessKey(cfg.AccessKey),
		env.WithDebugMode(false))

	return &Tuya{UserId: cfg.UserId}
}

func (c *Tuya) DeviceInfo(ctx context.Context, deviceId string) (*model.DeviceResponse, error) {
	requestUrl := fmt.Sprintf("/v1.0/devices/%s", deviceId)

	resp := &model.DeviceResponse{}
	err := connector.MakeGetRequest(ctx,
		connector.WithAPIUri(requestUrl),
		connector.WithResp(resp))

	if resp.Code > 0 {
		return nil, errors.New(resp.Msg)
	}

	return resp, err
}

func (c *Tuya) DeviceList(ctx context.Context) (*model.DeviceList, error) {
	requestUrl := "/v1.3/iot-03/devices?source_type=homeApp"

	resp := &model.DeviceList{}
	err := connector.MakeGetRequest(ctx,
		connector.WithAPIUri(requestUrl),
		connector.WithResp(resp))

	if !resp.Success {
		return nil, errors.New("failed to get device list")
	}

	return resp, err
}

func (c *Tuya) UserHomes(ctx context.Context, userId string) (*model.UserHomesResponse, error) {
	requestUrl := fmt.Sprintf("/v1.0/users/%s/homes", userId)

	resp := &model.UserHomesResponse{}

	err := connector.MakeGetRequest(ctx,
		connector.WithAPIUri(requestUrl),
		connector.WithResp(resp))

	return resp, err
}

func (c *Tuya) Scenarios(ctx context.Context, homeId string) (*model.ScenariosResponse, error) {
	requestUrl := fmt.Sprintf("/v1.1/homes/%s/scenes", homeId)
	resp := &model.ScenariosResponse{}
	err := connector.MakeGetRequest(ctx,
		connector.WithAPIUri(requestUrl),
		connector.WithResp(resp))

	return resp, err
}

func (c *Tuya) StartScenario(ctx context.Context, homeId, scenarioId string) (bool, error) {
	resp := &model.BaseResponse{}

	requestUrl := fmt.Sprintf("/v1.0/homes/%s/scenes/%s/trigger", homeId, scenarioId)
	err := connector.MakePostRequest(ctx,
		connector.WithAPIUri(requestUrl),
		connector.WithResp(resp))

	return resp.Success, err
}

func (c *Tuya) ScenarioInfo(ctx context.Context, homeId, scenarioId string) (*model.DeviceItem, error) {
	resp := &model.SceneDeviceResponse{}

	requestUrl := fmt.Sprintf("/v1.0/homes/%s/scene/devices", homeId)

	err := connector.MakeGetRequest(ctx,
		connector.WithAPIUri(requestUrl),
		connector.WithResp(resp))

	if resp.Result != nil {
		return &resp.Result[0], errors.New(resp.Msg) // TODO: this need to be refactored
	}

	return nil, err
}

func (c *Tuya) Actions(ctx context.Context) ([]Action, error) {
	if len(c.actions) != 0 {
		return c.actions, nil
	}

	c.actions = make([]Action, 0)

	homes, err := c.UserHomes(ctx, c.UserId)

	if err != nil {
		return nil, err
	}

	if !homes.Success {
		return nil, fmt.Errorf("%s", homes.Msg)
	}

	for _, home := range homes.Result {
		id := strconv.Itoa(home.HomeId)
		scenarios, err := c.Scenarios(ctx, id)

		if err != nil && len(scenarios.Result) == 0 {
			continue
		}

		for _, scene := range scenarios.Result {
			c.actions = append(c.actions, Action{
				ID:       scene.SceneId,
				HomeID:   home.HomeId,
				Name:     scene.Name,
				DeviceID: scene.Actions[0].EntityId,
			})
		}
	}

	return c.actions, err
}
