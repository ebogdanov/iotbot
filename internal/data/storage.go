package data

import "su27bot/internal/config"

type Storage struct {
	Users   *Users
	Codes   *Codes
	Groups  *Groups
	Actions *ActionLog
	Invites *Invites
	BotName string
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	db, err := NewDb(cfg)

	if err == nil {
		return &Storage{
			Users:   NewUsers(db),
			Codes:   NewCodes(db),
			Groups:  NewGroups(db),
			Actions: NewActionLog(db),
			Invites: NewInvites(db),
		}, nil
	}

	return nil, err
}
