package models

type AddressInfo struct {
	Address string
}

type ProfileInfo struct{
	FirstName string `db:"first_name"`
	LastName string `db:"last_name"`
	Telephone string
}

type Credentials struct {
	Email string
	Password string `db:"password_hash"`
}

type User struct{
	Id int
	IsActive bool `db:"is_active"`
	ProfileInfo
	AddressInfo
	Credentials
}

func GetUserById(userId int) (*User, error){
	user := User{}
	err := Storage.DB.QueryRowx("SELECT * FROM user WHERE id=? LIMIT 1", userId).StructScan(&user)
	return &user, err
}


func GetUserByEmail(email string) (*User, error){
	user := User{}
	err := Storage.DB.QueryRowx("SELECT * FROM user WHERE email=? LIMIT 1", email).StructScan(&user)
	return &user, err
}

func GetUserByCredentials(cr *Credentials) (*User, error){
	user := User{}
	err := Storage.DB.QueryRowx( "SELECT * FROM user WHERE email=? and password_hash=? LIMIT 1", cr.Email, cr.Password).StructScan(&user)
	return &user, err
}

func (user *User) Update() error{
	query := "UPDATE user SET email=?, first_name=?, last_name=?, address=?, telephone=?, is_active=? WHERE id=?"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(user.Email, user.FirstName, user.LastName, user.Address, user.Telephone, user.IsActive, user.Id)
	return err
}

func (user *User) Remove() error{
	query := "DELETE FROM user WHERE id=?"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(user.Id)
	return err
}

func (user *User) Insert() error{
	query := "INSERT INTO user(email, password_hash, first_name, last_name, address, telephone) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(user.Email, user.Password, user.FirstName, user.LastName, user.Address, user.Telephone)
	return err
}

func (cr *Credentials) Insert() error{
	query := "INSERT INTO user(email, password_hash) VALUES (?, ?)"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(cr.Email, cr.Password)
	return err
}

func (cr *Credentials) Update() error{
	query := "UPDATE user SET password_hash=? WHERE email=?"
	stmt, err := Storage.DB.Prepare(query)
	if err != nil{
		return err
	}
	_, err = stmt.Exec(cr.Password, cr.Email)
	return err
}
