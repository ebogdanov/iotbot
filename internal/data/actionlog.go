package data

import (
	"database/sql"
	"strconv"
	"time"
)

var (
	floodPeriod = -1 * time.Hour
	floodLimit  = 20
)

type Action struct {
	UserID, Cmd, HandlerName, User string
	Result                         bool
	EventTime                      string
	ID                             int
}

type ActionLog struct {
	db *Db
}

func NewActionLog(conn *Db) *ActionLog {
	return &ActionLog{db: conn}
}

func (a *ActionLog) Add(userID, userName, cmd, handlerName string, result bool) {
	go func() {
		cmd1 := []rune(cmd)
		cmd2 := cmd

		if len(cmd1) > 30 {
			cmd2 = string(cmd1[:30]) + "..."
		}
		_, _ = a.db.Conn.Exec(
			"INSERT INTO actions (user_id, username, cmd, handler, result, execute_time) VALUES ($1, $2, $3, $4, $5, strftime('%s', 'now'))",
			userID, userName, cmd2, handlerName, result)
	}()
}

func (a *ActionLog) Flood(cmd, userID string) bool {
	cnt := a.Count(cmd, userID)

	return cnt > floodLimit
}

func (a *ActionLog) Count(cmd, userID string) int {
	period := time.Now().Add(floodPeriod).Unix()

	var (
		res *sql.Rows
		err error
	)

	if cmd == "*" {
		res, err = a.db.Conn.Query(
			"SELECT count(id) as cnt FROM actions WHERE user_id = $1 and execute_time > $2", userID, period)
	} else {
		res, err = a.db.Conn.Query(
			"SELECT count(id) as cnt FROM actions WHERE user_id = $1 and cmd = $2 and execute_time > $3", userID, cmd, period)
	}

	if err != nil {
		return 0
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return 0
	}

	if res.Next() {
		var (
			cnt int
		)

		if err := res.Scan(&cnt); err != nil {
			return 0
		}

		return cnt
	}

	return 0
}

func (a *ActionLog) List(limit int) ([]Action, error) {
	return a.executeListQuery("actions.handler != $1", []interface{}{"admin"}, limit)
}

func (a *ActionLog) ListGroup(limit int, group string) ([]Action, error) {
	return a.executeListQuery("actions.handler = $1", []interface{}{group}, limit)
}

func (a *ActionLog) ListUser(userID string, limit int) ([]Action, error) {
	return a.executeListQuery("actions.user_id = $1", []interface{}{userID}, limit)
}

func (a *ActionLog) ListUnknown(limit int) ([]Action, error) {
	return a.executeListQuery("users.id is null", []interface{}{}, limit)
}

func (a *ActionLog) executeListQuery(cond string, args []interface{}, limit int) ([]Action, error) {
	if cond == "" {
		cond = "1=1"
	}

	args = append(args, limit)
	cntLimit := strconv.Itoa(len(args))

	//goland:noinspection SqlResolve
	query := "SELECT actions.user_id, cmd, handler, result, execute_time, coalesce(actions.username, users.name, actions.user_id) as username, users.id FROM actions " +
		" LEFT JOIN users on actions.user_id = users.user_id WHERE " + cond + " ORDER BY actions.id DESC LIMIT $" + cntLimit

	res, err := a.db.Conn.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return nil, err
	}

	list := make([]Action, 0)
	for res.Next() {
		item := &Action{}

		var (
			unixTime int64
			uID      sql.NullInt64
			userName sql.NullString
		)

		if err := res.Scan(&item.UserID, &item.Cmd, &item.HandlerName, &item.Result, &unixTime, &userName, &uID); err != nil {
			return list, err
		}

		item.ID = int(uID.Int64)
		item.User = userName.String

		dateTime := time.Unix(unixTime, 0)
		item.EventTime = dateTime.Format(time.RFC822)

		list = append(list, *item)
	}

	return list, nil
}
