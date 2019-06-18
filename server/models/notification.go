package models

import (
	"errors"
	"fmt"
	"time"
)

type Notification struct {
	NotificationID uint      `json:"notification_id,string"`
	UserID         uint      `json:"user_id,string"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"created_at"`
}

func GetNotifications(userID string) ([]Notification, error) {
	db := InitDB()

	res, err := db.Query("SELECT * FROM notifications WHERE user_id=? ORDER BY user_id DESC", userID)
	if err != nil {
		return nil, errors.New("Error to get notifications from db")
	}
	defer db.Close()

	notification := Notification{}
	notificationArr := []Notification{}
	for res.Next() {
		err = res.Scan(&notification.NotificationID, &notification.UserID,
			&notification.Text, &notification.CreatedAt)
		if err != nil {
			return nil, errors.New("Error to scan notifications from db")
		}
		notificationArr = append(notificationArr, notification)
	}

	return notificationArr, nil
}

func RemoveNotification(notificationID string) error {
	fmt.Print("REMOVE NOTIFICATION")
	db := InitDB()
	_, err := db.Exec("DELETE FROM notifications where notification_id=?;", notificationID)
	defer db.Close()

	if err != nil {
		return errors.New("Delete notification error")
	}
	return nil
}

//var ADMIN_PROMOTION_MSG = "You have been promoted to manager user"
//var ADMIN_DEMOTION_MSG = "You have been demoted to user"
//
//var USER_ADDING_TO_QUEUE_MSG = "A new user was added to your property listing queue"
//var USER_REMOVING_TO_QUEUE_MSG = "A user was removed from your property listing queue"
//
//func AddNotification(userID string, notification string) error  {
//	db := InitDB()
//	_, err := db.Exec("INSERT INTO notifications(user_id, text) VALUES (?,?);",
//		userID, notification)
//	defer db.Close()
//
//	if err != nil {
//		return errors.New("Notification insertion error")
//	}
//	return nil
//}
