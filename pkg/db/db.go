package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

const (
	schema = `create table scheduler (
    id integer primary key,
    date varchar(8) not null default '',
    title varchar(128) not null,
    comment text not null default '',
    repeat varchar(128) not null default ''
);
create index idx_scheduler_date on scheduler (date);
`
)

var db *sql.DB

func Init(dbFile string) error {
	var install bool
	_, err := os.Stat(dbFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			install = true
		} else {
			return err
		}
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
