package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aplJake/reals-course/server/models"
	"github.com/aplJake/reals-course/server/utils"
	"log"
	"net/http"
)

func GetPropertyQueue(w http.ResponseWriter, r *http.Request) {
	var propertyPageData models.PropertyCtxData

	propertyPageData = r.Context().Value("propertyData").(models.PropertyCtxData)
	resp := utils.Message(true, "Queues are sended")
	resp["property_queue_data"] = propertyPageData
	utils.Respond(w, resp)
}

func AddUserToQueue(w http.ResponseWriter, r *http.Request) {
	q := &models.Queue{}

	// Decode the request to server
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		panic(err.Error())
		utils.Respond(w, utils.Message(false, "Cannot decode recieved json object"))
		return
	}

	fmt.Print(q)
	// CreateSeller new User and UserProfile
	err = q.Create()
	if err != nil {
		utils.Respond(w, utils.Message(false, "Cannot add new queue object to the database"))
		log.Fatal(err.Error())
		return
	}

	resp := utils.Message(true, "New Queue user was successfully added")
	utils.Respond(w, resp)
}

func DeleteQueueUser(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("DELETE OPERATION ")
	var queue models.Queue

	queue = r.Context().Value("queueData").(models.Queue)

	db := models.InitDB()

	_, err := db.Exec("DELETE FROM property_queue WHERE user_id=? AND property_id=?",
		queue.UserID, queue.PropertyID)
	defer db.Close()

	if err != nil {
		utils.Respond(w, utils.Message(true, "Queue Delete error from db"))
	}

	utils.Respond(w, utils.Message(true, "Queue was successfully deleted"))
}