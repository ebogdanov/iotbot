package iot

import (
	"context"
	"github.com/ebogdanov/ewelink"
	"strconv"
	"su27bot/internal/model"
	"time"
)

type EwelinkWebSocket struct {
	actions  []Action
	session  *ewelink.Session
	client   *ewelink.Ewelink
	username string
	password string
	region   string
}

func NewEwelinkWebsocket(region, userName, password string) *EwelinkWebSocket {
	instance := ewelink.New()

	return &EwelinkWebSocket{
		client:   instance,
		username: userName,
		password: password,
		region:   region,
	}
}

func (e *EwelinkWebSocket) DeviceInfo(ctx context.Context, ID string) (*model.DeviceResponse, error) {
	if err := e.auth(); err != nil {
		return nil, err
	}

	info, err := e.client.GetDevice(ctx, e.session, ID)
	if err != nil {
		return nil, err
	}

	return &model.DeviceResponse{
		Code:    200,
		Msg:     "Success",
		Success: true,
		Result: model.DeviceModel{
			Sub:         false,
			ID:          info.DeviceID,
			LocalKey:    info.Devicekey,
			IP:          info.IP,
			TimeZone:    "",
			ProductName: info.BrandName,
			ActiveTime:  info.OnlineTime.Unix(),
			UpdateTime:  info.OnlineTime.Unix(),
			Name:        info.Name,
			Online:      info.Online,
			Model:       info.ProductModel,
			UID:         info.ID,
			Category:    info.Type,
		}}, nil
}

func (e *EwelinkWebSocket) DeviceList(ctx context.Context) (*model.DeviceList, error) {
	if err := e.auth(); err != nil {
		return nil, err
	}

	devices, err := e.client.GetDevices(ctx, e.session)
	if err != nil {
		return nil, err
	}

	response := &model.DeviceList{}

	for i := range devices.Devicelist {
		item := devices.Devicelist[i]

		devItem := model.Device{
			Sub:          false,
			CategoryName: item.Group,
			CreateTime:   item.CreatedAt.Unix(),
			LocalKey:     item.Devicekey,
			OwnerId:      "",
			Ip:           item.IP,
			Icon:         item.BrandLogoURL,
			Lon:          "",
			TimeZone:     "",
			ProductName:  item.BrandName,
			Uuid:         item.Apikey,
			GatewayId:    "",
			ActiveTime:   item.OnlineTime.Unix(),
			UpdateTime:   item.OnlineTime.Unix(),
			ProductId:    "",
			Name:         item.Name,
			Online:       item.Online,
			Model:        item.ProductModel,
			Id:           item.DeviceID,
			Category:     item.Type,
			Lat:          "",
		}

		response.Result.List = append(response.Result.List, devItem)
	}

	return response, err
}

func (e *EwelinkWebSocket) UserHomes(_ context.Context, _ string) (*model.UserHomesResponse, error) {
	return &model.UserHomesResponse{}, nil
}

func (e *EwelinkWebSocket) Scenarios(_ context.Context, _ string) (*model.ScenariosResponse, error) {
	return &model.ScenariosResponse{}, nil
}

func (e *EwelinkWebSocket) StartScenario(ctx context.Context, homeId, scenarioId string) (bool, error) {
	ctx1, cancelFunc := context.WithTimeout(ctx, 15*time.Second)
	defer func() { cancelFunc() }()

	if err := e.auth(); err != nil {
		return false, err
	}

	uid, err := strconv.Atoi(homeId)
	if err != nil {
		return false, err
	}

	device := &ewelink.Device{DeviceID: scenarioId, Uiid: uid}
	_, err = e.client.SetDevicePowerState(ctx1, e.session, device, true)

	return err == nil, err
}

func (e *EwelinkWebSocket) ScenarioInfo(_ context.Context, _, _ string) (*model.DeviceItem, error) {
	return nil, nil
}

func (e *EwelinkWebSocket) Actions(ctx context.Context) ([]Action, error) {
	if len(e.actions) != 0 {
		return e.actions, nil
	}

	// List devices
	if err := e.auth(); err != nil {
		return nil, err
	}

	e.actions = make([]Action, 0)
	devices, err := e.client.GetDevices(ctx, e.session)
	if err == nil {
		for i := range devices.Devicelist {
			e.actions = append(e.actions, Action{
				ID:       devices.Devicelist[i].ID,
				Name:     devices.Devicelist[i].Name,
				HomeID:   devices.Devicelist[i].Uiid,
				DeviceID: devices.Devicelist[i].DeviceID,
			})
		}
	}

	return e.actions, err
}

func (e *EwelinkWebSocket) auth() error {
	if e.session != nil {
		return nil
	}

	session, err := e.client.AuthenticateWithEmail(
		context.Background(), ewelink.NewConfiguration(e.region), e.username, e.password)

	if err == nil {
		e.session = session
	}

	return err
}
