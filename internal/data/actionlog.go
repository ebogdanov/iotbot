package data

import (
	"database/sql"
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
}

type ActionLog struct {
	db *Db
}

func NewActionLog(conn *Db) *ActionLog {
	return &ActionLog{db: conn}
}

func (a *ActionLog) Add(userID, cmd, handlerName string, result bool) {
	go func() {
		cmd1 := cmd
		if len(cmd1) > 30 {
			cmd1 = cmd[:30]
		}
		_, _ = a.db.Conn.Exec(
			"INSERT INTO actions (user_id, cmd, handler, result, execute_time) VALUES ($1, $2, $3, $4, strftime('%s', 'now'))",
			userID, cmd1, handlerName, result)
	}()
}

func (i *ActionLog) Flood(cmd, userID string) bool {
	cnt := i.Count(cmd, userID)

	return cnt > floodLimit
}

func (i *ActionLog) Count(cmd, userID string) int {
	period := time.Now().Add(floodPeriod).Unix()

	var (
		res *sql.Rows
		err error
	)

	if cmd == "*" {
		res, err = i.db.Conn.Query(
			"SELECT count(id) as cnt FROM actions WHERE user_id = $1 and execute_time > $2", userID, period)
	} else {
		res, err = i.db.Conn.Query(
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

func (i *ActionLog) List(limit int) ([]Action, error) {
	res, err := i.db.Conn.Query(
		"SELECT actions.user_id, cmd, handler, result, execute_time, users.name FROM actions LEFT JOIN users on actions.user_id = users.user_id WHERE actions.handler != $1 ORDER BY actions.id DESC LIMIT $2",
		"admin", limit)

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
		)

		if err := res.Scan(&item.UserID, &item.Cmd, &item.HandlerName, &item.Result, &unixTime, &item.User); err != nil {
			continue
		}

		dateTime := time.Unix(unixTime, 0)
		item.EventTime = dateTime.Format(time.RFC822)

		list = append(list, *item)
	}

	return list, nil
}

func (i *ActionLog) ListGroup(limit int, group string) ([]Action, error) {
	res, err := i.db.Conn.Query(
		"SELECT actions.user_id, cmd, handler, result, execute_time, users.name FROM actions LEFT JOIN users on actions.user_id = users.user_id WHERE actions.handler = $2 ORDER BY actions.id DESC LIMIT $1",
		limit, group)

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
		)

		if err := res.Scan(&item.UserID, &item.Cmd, &item.HandlerName, &item.Result, &unixTime, &item.User); err != nil {
			continue
		}

		dateTime := time.Unix(unixTime, 0)
		item.EventTime = dateTime.Format(time.RFC822)

		list = append(list, *item)
	}

	return list, nil
}

func (i *ActionLog) ListUser(userID string, limit int) ([]Action, error) {
	res, err := i.db.Conn.Query(
		"SELECT actions.user_id, cmd, handler, result, execute_time, users.name FROM actions LEFT JOIN users on actions.user_id = users.user_id WHERE actions.user_id = $1 ORDER BY actions.id DESC LIMIT $2",
		userID, limit)

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
		)

		if err := res.Scan(&item.UserID, &item.Cmd, &item.HandlerName, &item.Result, &unixTime, &item.User); err != nil {
			continue
		}

		dateTime := time.Unix(unixTime, 0)
		item.EventTime = dateTime.Format(time.RFC822)

		list = append(list, *item)
	}

	return list, nil
}
