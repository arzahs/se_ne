package models

import (
	"se_ne/db"
	"log"
)

var (
	Storage *db.Storage
)

func GetUserByToken(token string) (*User, error){
	session, err := GetSessionByToken(token)
	log.Print("session", session, err)
	if err != nil{
		return nil, err
	}

	user, err := GetUserById(session.UserId)
	if err != nil{
		return nil, err
	}

	return user, nil
}


func InitModels(cfg db.Config) (err error){
	Storage, err = db.NewStorage(cfg)
	return err
}


