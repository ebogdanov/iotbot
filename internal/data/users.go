package data

import "strconv"

type Users struct {
	db *Db
}

func NewUsers(conn *Db) *Users {
	return &Users{db: conn}
}

func (u *Users) Add(userID, userName string) (bool, error) {
	res, err := u.db.Conn.Exec("INSERT INTO users (user_id, name, active) VALUES ($1, $2, $3)", userID, userName, 1)
	if err != nil {
		return false, err
	}
	cnt, err := res.LastInsertId()
	if err != nil {
		return false, err
	}
	return cnt > 1, nil
}

func (u *Users) Delete(userID string) bool {
	res, err := u.db.Conn.Exec("UPDATE users SET active = $1 WHERE user_id = $2", 0, userID)
	if err != nil {
		return false
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return cnt > 0
}

func (u *Users) Check(userID string) bool {
	res, err := u.db.Conn.Query(
		"SELECT user_id, name FROM users WHERE active = 1 and user_id = $1 LIMIT 1", userID)

	if err != nil {
		return false
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return false
	}

	return res.Next()
}

func (u *Users) Info(userID string) (string, string, bool, error) {
	res, err := u.db.Conn.Query(
		"SELECT user_id, name, active FROM users WHERE user_id = $1 LIMIT 1", userID)

	if err != nil {
		return "", "", false, err
	}
	defer func() { _ = res.Close() }()

	if res.Err() == nil {
		if res.Next() {
			var (
				userID   string
				userName string
				active   int
			)

			if err = res.Scan(&userID, &userName, &active); err == nil {
				return userID, userName, active == 1, nil
			}
		}
	}

	return "", "", false, err
}

func (u *Users) Name(userID string) string {
	res, err := u.db.Conn.Query(
		"SELECT name FROM users WHERE user_id = $1 LIMIT 1", userID)

	if err != nil {
		return ""
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return ""
	}

	for res.Next() {
		var (
			userName string
		)

		if err := res.Scan(&userName); err != nil {
			continue
		}

		return userName
	}

	return ""
}

func (u *Users) Active() map[string]string {
	list := make(map[string]string)

	res, err := u.db.Conn.Query("SELECT user_id, name FROM users WHERE active = 1")
	if err != nil {
		return list
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return list
	}

	for res.Next() {
		var (
			userId   int64
			userName string
		)

		if err := res.Scan(&userId, &userName); err != nil {
			continue
		}

		id := strconv.FormatInt(userId, 10)
		list[id] = userName
	}

	return list
}
