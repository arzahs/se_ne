package models

type ResetRequest struct {
	UserId int `db:"user_id"`
	Token  string
}

func (rr *ResetRequest) Insert() error {
	query := "INSERT INTO reset_request(user_id, token) VALUES (?, ?)"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rr.UserId, rr.Token)
	return err
}

func GetResetRequestByToken(token string) (*ResetRequest, error) {
	req := ResetRequest{}
	err := Storage.DB.QueryRowx("SELECT * FROM reset_request WHERE token=? LIMIT 1", token).StructScan(&req)
	return &req, err
}

func DeleteResetRequestByToken(token string) error {
	query := "DELETE FROM reset_request WHERE token=?"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(token)
	return err
}
