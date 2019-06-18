package models

import (
	"errors"
	"fmt"
	"time"
)

type Queue struct {
	UserID     string    `json:"user_id"`
	PropertyID string    `json:"property_id"`
	QueueTime  time.Time `json:"queue_time"`
}

func (q *Queue) Create() error {
	fmt.Printf("userid %s, proeprtyID %s", q.UserID, q.PropertyID)
	db := InitDB()
	_, err := db.Exec("INSERT INTO property_queue(user_id, property_id) VALUES(?,?)",
		q.UserID, q.PropertyID)
	defer db.Close()

	if err != nil {
		panic(err.Error())
		return errors.New("Insert Queue db error")
	}

	return nil
}
