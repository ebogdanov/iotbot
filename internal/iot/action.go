package iot

import "fmt"

type Action struct {
	ID       string
	Name     string
	HomeID   int
	DeviceID string
}

func (a *Action) MenuId() string {
	return fmt.Sprintf("DO_%d_%s", a.HomeID, a.ID)
}
