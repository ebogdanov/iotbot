package data

type Codes struct {
	db *Db
}

type Code struct {
	UserID      string `db:"user_id"`
	Title       string `db:"title"`
	Code        string `db:"code"`
	Attempts    int    `db:"attempts"`
	MaxAttempts int    `db:"max_attempts"`
}

func NewCodes(conn *Db) *Codes {
	return &Codes{db: conn}
}

func (c *Codes) Add(code *Code) (bool, error) {
	res, err := c.db.Conn.Exec(
		"INSERT INTO codes (code, user_id, title, max_attempts) VALUES ($1, $2, $3, $4)", code.Code, code.UserID, code.Title, code.MaxAttempts)

	if err != nil {
		return false, err
	}
	cnt, err := res.LastInsertId()
	if err != nil {
		return false, err
	}
	return cnt > 1, nil
}

func (c *Codes) Info(code string) (*Code, error) {
	res, err := c.db.Conn.Query(
		"SELECT code, user_id, title, attempts, max_attempts FROM codes WHERE code = $1 AND (max_attempts = 0 OR attempts < max_attempts) LIMIT 1", code)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return nil, res.Err()
	}

	if res.Next() {
		c := &Code{}

		err := res.Scan(&c.Code, &c.UserID, &c.Title, &c.Attempts, &c.MaxAttempts)
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	return nil, nil
}

func (c *Codes) Check(code string) bool {
	res, err := c.db.Conn.Query(
		"SELECT code FROM codes WHERE code = $1 AND (max_attempts = 0 OR attempts < max_attempts) LIMIT 1", code)

	if err != nil {
		return false
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return false
	}

	return res.Next()
}

func (c *Codes) Use(code string) bool {
	res, err := c.db.Conn.Exec(
		"UPDATE codes SET attempts = attempts + 1 WHERE code = $1 AND (max_attempts = 0 OR attempts < max_attempts)", code)

	if err != nil {
		return false
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return cnt > 1
}

func (c *Codes) Delete(code, userId string) bool {
	res, err := c.db.Conn.Exec("DELETE FROM codes WHERE code = $1 AND user_id = $2", code, userId)
	if err != nil {
		return false
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return cnt > 0
}

func (c *Codes) List(userID string) ([]*Code, error) {
	res, err := c.db.Conn.Query("SELECT code, user_id, title, attempts, max_attempts  FROM codes WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return nil, res.Err()
	}

	list := make([]*Code, 0)
	for res.Next() {
		c := &Code{}

		err = res.Scan(&c.Code, &c.UserID, &c.Title, &c.Attempts, &c.MaxAttempts)
		if err != nil {
			continue
		}

		list = append(list, c)
	}

	return list, err
}
