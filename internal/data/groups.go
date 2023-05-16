package data

import (
	"database/sql"
	"strconv"
	"time"
)

type Groups struct {
	db *Db
}

const (
	adminGroup = "admin"
)

func NewGroups(conn *Db) *Groups {
	return &Groups{
		db: conn,
	}
}

func (g *Groups) IsMember(group, userID string) bool {
	var (
		res *sql.Rows
		err error
	)

	if group == "all" {
		res, err = g.db.Conn.Query(
			"SELECT user_id, name FROM users WHERE active = 1 and user_id = $1 LIMIT 1", userID)
	} else {
		res, err = g.db.Conn.Query(
			"SELECT groups.title, valid_till FROM membership INNER JOIN groups ON membership.id_group = groups.id WHERE groups.title = $1 and user_id = $2 LIMIT 1", group, userID)
	}

	if err != nil {
		return false
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return false
	}

	return res.Next()
}

func (g *Groups) MemberOf(userID string) []string {
	memberOf := make([]string, 0)

	res, err := g.db.Conn.Query(
		"SELECT groups.title, membership.valid_till FROM membership INNER JOIN groups ON membership.id_group = groups.id WHERE user_id = $1", userID)

	if err != nil {
		return memberOf
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return memberOf
	}

	for res.Next() {
		var (
			title     string
			validTill sql.NullInt64
		)

		if err := res.Scan(&title, &validTill); err != nil {
			continue
		}

		if validTill.Int64 > 0 && (validTill.Int64 > time.Now().Unix()) {
			continue
		}

		memberOf = append(memberOf, title)
	}

	return memberOf
}

func (g *Groups) List() map[string]string {
	list := make(map[string]string)

	res, err := g.db.Conn.Query("SELECT id, title FROM groups")
	if err != nil {
		return list
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return list
	}

	for res.Next() {
		var (
			groupId int64
			title   string
		)

		if err := res.Scan(&groupId, &title); err != nil {
			continue
		}

		id := strconv.FormatInt(groupId, 10)
		list[id] = title
	}

	return list
}

func (g *Groups) Members(title string) map[string]string {
	list := make(map[string]string, 0)

	res, err := g.db.Conn.Query(
		"SELECT users.user_id, users.name, membership.valid_till FROM users INNER JOIN membership ON membership.user_id = users.user_id INNER JOIN groups ON membership.id_group = groups.id WHERE title = $1 and users.active = 1", title)

	if err != nil {
		return list
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return list
	}

	for res.Next() {
		var (
			userId    int64
			validTill sql.NullInt64
			userName  string
		)

		if err := res.Scan(&userId, &userName, &validTill); err != nil {
			continue
		}

		if validTill.Int64 > 0 && validTill.Int64 > time.Now().Unix() {
			continue
		}

		id := strconv.FormatInt(userId, 10)
		list[id] = userName
	}

	return list
}

func (g *Groups) DeleteMember(userID, groupID string) bool {
	// Delete from membership table for all groups
	var (
		res sql.Result
		err error
	)

	if groupID == "*" {
		res, err = g.db.Conn.Exec("DELETE FROM membership WHERE user_id = $1", userID)
	} else {
		res, err = g.db.Conn.Exec("DELETE FROM membership WHERE user_id = $1 and id_group = $2", userID, groupID)
	}

	if err != nil {
		return false
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return cnt > 0
}

func (g *Groups) AddMember(groupID, userID string) bool {
	// Add into membership table
	res, err := g.db.Conn.Exec("INSERT INTO membership (user_id, id_group) VALUES($1, $2)", userID, groupID)

	if err != nil {
		return false
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return cnt > 0
}

func (g *Groups) Title(groupID string) string {
	title := ""

	res, err := g.db.Conn.Query("SELECT title FROM groups WHERE id = $1", groupID)
	if err != nil {
		return title
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return title
	}

	if res.Next() {
		if err := res.Scan(&title); err != nil {
			return ""
		}
	}

	return title
}

func (g *Groups) Delete(groupID string) bool {
	// Delete from membership table for all groups
	tx, err := g.db.Conn.Begin()
	if err != nil {
		return false
	}
	res, err := tx.Exec("DELETE FROM membership WHERE id_group = $1", groupID)

	if err == nil {
		res, err = tx.Exec("DELETE FROM groups WHERE id = $1", groupID)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return false
	}
	err = tx.Commit()
	if err != nil {
		return false
	}

	return cnt > 0
}

func (g *Groups) IsAdmin(userID string) bool {
	return g.IsMember(adminGroup, userID)
}
