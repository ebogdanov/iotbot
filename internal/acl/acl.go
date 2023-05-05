package acl

import (
	"su27bot/internal/config"
	"sync"
)

const (
	GroupAdmin = "admin"
)

type Default struct {
	m   sync.Mutex
	cfg *config.Config

	members sync.Map
	users   sync.Map
}

type MemberList map[string]bool

func New(cfg *config.Config) *Default {
	s := &Default{
		cfg:     cfg,
		m:       sync.Mutex{},
		members: sync.Map{},
		users:   sync.Map{},
	}

	// Load data about users and groups
	for name, g := range cfg.Acl.Groups {
		members := make(MemberList, 0)

		for _, u := range g {
			members[u] = true

			s.users.Store(u, u)
		}

		s.members.Store(name, members)
	}

	for userId, userName := range cfg.Acl.Users {
		s.users.Store(userId, userName)
	}

	return s
}

func (u *Default) AddInvite(id string) {
	u.m.Lock()
	defer u.m.Unlock()

	u.cfg.Acl.Invites = append(u.cfg.Acl.Invites, id)
	_ = u.cfg.SaveFile()
}

func (u *Default) CheckInvite(id string) bool {
	u.m.Lock()
	defer u.m.Unlock()

	found := false

	if u.cfg.Acl != nil {
		for _, item := range u.cfg.Acl.Invites {
			if id == item {
				found = true
				break
			}
		}
	}

	return found
}

func (u *Default) RemoveInvite(id string) {
	u.m.Lock()
	defer u.m.Unlock()

	for i, item := range u.cfg.Acl.Invites {
		if id == item {
			// Move item from last position to found position
			u.cfg.Acl.Invites[i] = u.cfg.Acl.Invites[len(u.cfg.Acl.Invites)-1]
			// Reduce length of item array
			u.cfg.Acl.Invites = u.cfg.Acl.Invites[:len(u.cfg.Acl.Invites)-1]
			break
		}
	}

	_ = u.cfg.SaveFile()
}

func (u *Default) Add(userId, userName string) {
	u.m.Lock()
	defer u.m.Unlock()
	// Add user to storage
	u.cfg.Acl.Users[userId] = userName
	// Add user to our data structure
	u.users.Store(userId, userName)

	// If this is first user, and there is no admin group yet - this is our Parent, add him to admins group
	if len(u.cfg.Acl.Users) == 1 {
		_, ok := u.cfg.Acl.Groups["admin"]
		if !ok {
			u.cfg.Acl.Groups["admin"] = []string{userId}
		}
	}

	_ = u.cfg.SaveFile()
}

func (u *Default) Remove(userId string) bool {
	u.m.Lock()
	defer u.m.Unlock()

	for i, item := range u.cfg.Acl.Users {
		if userId == item || i == userId {
			delete(u.cfg.Acl.Users, i)

			// Find in groups and delete this userId
			groupMembers := u.MemberOfGroups(userId)
			for _, group := range groupMembers {
				if members, ok := u.members.Load(group); ok {
					delete(members.(MemberList), userId)
					u.members.Store(group, members)
				}
			}

			for groupName, groupMembers := range u.cfg.Acl.Groups {
				for i, groupUser := range groupMembers {
					if groupUser == userId {
						u.cfg.Acl.Groups[groupName][i] = groupMembers[len(groupMembers)-1]
						u.cfg.Acl.Groups[groupName] = groupMembers[:len(groupMembers)-1]
					}
				}
			}

			_ = u.cfg.SaveFile()
			return true
		}
	}

	return false
}

func (u *Default) GetUserList() map[string]string {
	return u.cfg.Acl.Users
}

func (u *Default) IsMember(group, userID string) bool {
	if group == "all" {
		if _, ok := u.users.Load(userID); ok {
			return ok
		}
	} else {
		if data, ok := u.members.Load(group); ok {
			if enabled, ok := data.(MemberList)[userID]; ok {
				return enabled
			}
		}
	}

	return false
}

func (u *Default) MemberOfGroups(userID string) []string {
	memberOf := make([]string, 0)
	u.members.Range(func(key, value any) bool {
		if _, ok := value.(MemberList)[userID]; ok {
			memberOf = append(memberOf, key.(string))
		}

		return true
	})

	return memberOf
}

func (u *Default) GetMembers() *sync.Map {
	return &u.members
}

func (u *Default) GetGroupMembers(group string) map[string]string {
	members := make(map[string]string, 0)

	if group, ok := u.members.Load(group); ok {
		for user := range group.(MemberList) {
			members[user] = user
		}
	}

	return members
}

func (u *Default) IotActions(userID string) []string {
	res := u.cfg.Acl.Actions.Allow

	return res
}

func (u *Default) IsAllowed(userID, actionName string) bool {
	// Check that user is registered
	enabled, ok := u.users.Load(userID)
	if !ok {
		return false
	}

	if enabled.(string) == "" {
		return false
	}

	if contains(u.cfg.Acl.Actions.Allow, actionName) {
		return true
	}

	if u.cfg.Acl.Actions.Only != nil {
		groups := u.MemberOfGroups(userID)
		for _, g := range groups {
			if _, ok := u.cfg.Acl.Actions.Only[g]; ok {
				return contains(u.cfg.Acl.Actions.Only[g], actionName)
			}
		}
	}

	return false
}

func (u *Default) IsAdmin(userId string) bool {
	if u.IsMember(GroupAdmin, userId) {
		// Check that userId is member of Admin group
		if _, ok := u.users.Load(userId); ok {
			return true
		}
	}

	return false
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
