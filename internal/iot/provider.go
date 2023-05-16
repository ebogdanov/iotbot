package iot

import (
	"context"
	"su27bot/internal/model"
)

type Provider interface {
	DeviceInfo(context.Context, string) (*model.DeviceResponse, error)
	DeviceList(context.Context) (*model.DeviceList, error)
	UserHomes(context.Context, string) (*model.UserHomesResponse, error)
	Scenarios(context.Context, string) (*model.ScenariosResponse, error)
	StartScenario(context.Context, string, string) (bool, error)
	ScenarioInfo(context.Context, string, string) (*model.DeviceItem, error)
	Actions(context.Context, []string) ([]Action, error)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
