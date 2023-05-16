package data

type Invites struct {
	db *Db
}

func NewInvites(conn *Db) *Invites {
	return &Invites{db: conn}
}

func (i *Invites) Add(code string) (bool, error) {
	res, err := i.db.Conn.Exec("INSERT INTO invites (code, active) VALUES ($1, $2)", code, 1)
	if err != nil {
		return false, err
	}
	cnt, err := res.LastInsertId()
	if err != nil {
		return false, err
	}
	return cnt > 1, nil
}

func (i *Invites) Check(code string) bool {
	res, err := i.db.Conn.Query("SELECT code FROM invites WHERE active = 1 and code = $1 LIMIT 1", code)

	if err != nil {
		return false
	}
	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return false
	}

	return res.Next()
}

func (i *Invites) Delete(code string) bool {
	res, err := i.db.Conn.Exec("DELETE FROM invites WHERE code = $1", code)
	if err != nil {
		return false
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return false
	}

	return cnt > 0
}
