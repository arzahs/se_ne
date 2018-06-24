package models

type Session struct {
	UserId int `db:"user_id"`
	Token  string
}

func (session *Session) Insert() error {
	query := "INSERT INTO session(user_id, token) VALUES (?, ?)"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(session.UserId, session.Token)
	return err
}

func GetSessionByToken(token string) (*Session, error) {
	session := Session{}
	err := Storage.DB.QueryRowx("SELECT * FROM session WHERE token=? LIMIT 1", token).StructScan(&session)
	return &session, err
}

func DeleteSessionByToken(token string) error {
	query := "DELETE FROM session WHERE token=?"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(token)
	return err
}
