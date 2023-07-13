package acl

import (
	"su27bot/internal/config"
	"su27bot/internal/data"
)

type Default struct {
	cfg *config.Config
	s   *data.Storage
}

type MemberList map[string]bool

func New(cfg *config.Config, s *data.Storage) *Default {
	return &Default{
		cfg: cfg,
		s:   s,
	}
}

func (d *Default) IsAllowed(userID, actionName string) bool {
	if contains(d.cfg.Acl.Actions.Allow, actionName) {
		return d.s.Users.Check(userID)
	}

	if d.cfg.Acl.Actions.Only != nil {
		for _, g := range d.s.Groups.MemberOf(userID) {
			if _, ok := d.cfg.Acl.Actions.Only[g]; ok {
				return contains(d.cfg.Acl.Actions.Only[g], actionName)
			}
		}
	}

	return false
}

func (d *Default) SupportActions() []string {
	list := d.cfg.Acl.Actions.Allow

	if d.cfg.Acl.Actions.Only != nil {
		for _, a := range d.cfg.Acl.Actions.Only {
			list = append(list, a...)
		}
	}

	return list
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
