package controllers

import (
	"github.com/aplJake/reals-course/server/models"
	"github.com/aplJake/reals-course/server/utils"
	"net/http"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	var (
		notificationArr []models.Notification
		err      error
	)

	userID := r.Context().Value("userID").(string)

	notificationArr, err = models.GetNotifications(userID)
	if err != nil {
		utils.Respond(w, utils.Message(true, err.Error()))
		return
	}

	resp := utils.Message(true, "Notifications are sended")
	resp["notifications"] = notificationArr
	// Respond to the client and ...
	utils.Respond(w, resp)
}

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	notificationID := r.Context().Value("notificationID").(string)

	err := models.RemoveNotification(notificationID)
	if err != nil {
		utils.Respond(w, utils.Message(true, "Error while removing notification from the db"))
		return
	}

	utils.Respond(w, utils.Message(true, "Notification was successfully removed"))

}