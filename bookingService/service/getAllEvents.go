package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (eh *BookingServiceHandler) GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.DBHandler.GetAllEvents()
	if err != nil {
		http.Error(w, "Cannot get all events", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		http.Error(w, "Cannot encode all events to JSON", http.StatusInternalServerError)
	}
	fmt.Println("Got all events")
}
