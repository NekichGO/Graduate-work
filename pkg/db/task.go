package db

import (
	"database/sql"
	"errors"
)

// ErrNotFound is returned when TaskByID and UpdateTask return zero rows
var ErrNotFound = errors.New("task not found")

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	res, err := db.Exec("insert into scheduler (date, title, comment, repeat) values (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task

	rows, err := db.Query("select id, date, title, comment, repeat from scheduler order by date asc limit :limit",
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func SearchTasks(search string, limit int) ([]*Task, error) {
	var tasks []*Task

	rows, err := db.Query("select id, date, title, comment, repeat from scheduler where title like :search or comment like :search order by date limit :limit",
		sql.Named("search", "%"+search+"%"),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func SearchTasksByDate(date string, limit int) ([]*Task, error) {
	var tasks []*Task

	rows, err := db.Query("select id, date, title, comment, repeat from scheduler where date = :date order by id limit :limit",
		sql.Named("date", date),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func TaskByID(id int64) (*Task, error) {
	var task Task

	row := db.QueryRow("select id, date, title, comment, repeat from scheduler where id = :id",
		sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	res, err := db.Exec("update scheduler set date = :date, title = :title, comment = :comment, repeat = :repeat where id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteTask(id int64) error {
	rows, err := db.Exec("delete from scheduler where id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func UpdateTaskDate(next string, id int64) error {
	res, err := db.Exec("update scheduler set date = :next where id = :id",
		sql.Named("next", next),
		sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
