package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type Config struct{
	DatabaseHost string `env:"APP_DB_HOST,required"`
	DatabasePort int    `env:"APP_DB_PORT,default=3306"`
	DatabaseName string `env:"APP_DB_NAME,required"`
	DatabaseUsername string `env:"APP_DB_USERNAME,required"`
	DatabasePassword string `env:"APP_DB_PASSWORD,required"`
}

func (c *Config) GetSourceName() string{
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?multiStatements=true",
		c.DatabaseUsername,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName)
}


type Storage struct{
	DB *sqlx.DB
}

func NewStorage(cfg Config) (*Storage, error){
	connPool, err := sqlx.Connect(
		"mysql", cfg.GetSourceName())
	if err != nil{
		return nil, err
	}

	err = connPool.Ping()
	if err != nil{
		return nil, err
	}

	err = InstallScheme(connPool)
	return &Storage{
		DB: connPool,
	}, err
}

func InstallScheme(db *sqlx.DB) error{
	scheme := `
CREATE TABLE IF NOT EXISTS user (
	id int NOT NULL AUTO_INCREMENT,
	email varchar(255) NOT NULL,
	password_hash varchar(255) NOT NULL,
    first_name varchar(255) DEFAULT "",
    last_name varchar(255) DEFAULT "",
    address varchar(255) DEFAULT "",
	telephone varchar(16) DEFAULT "",
	is_active BOOLEAN DEFAULT FALSE,
	PRIMARY KEY (id),
	UNIQUE (email)
) ENGINE=InnoDB CHARACTER SET=UTF8;

CREATE TABLE IF NOT EXISTS session (
    user_id int NOT NULL,
    token varchar(255) NOT NULL,
	PRIMARY KEY (token),
	FOREIGN KEY (user_id) REFERENCES user(id)
) ENGINE=InnoDB CHARACTER SET=UTF8;

CREATE TABLE IF NOT EXISTS reset_request (
    user_id int NOT NULL,
    token varchar(255) NOT NULL,
	PRIMARY KEY (token),
	FOREIGN KEY (user_id) REFERENCES user(id)
) ENGINE=InnoDB CHARACTER SET=UTF8;`
	_, err := db.Exec(scheme)
	return err
}
