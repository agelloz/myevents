package service

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/agelloz/reach/contracts"
	"net/http"
	"time"

	"github.com/agelloz/reach/eventsService/models"
)

func (eh *EventsServiceHandler) AddEventHandler(w http.ResponseWriter, r *http.Request) {
	event := models.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Cannot decode event data", http.StatusInternalServerError)
		return
	}
	id, err := eh.DbHandler.AddEvent(event)
	if nil != err {
		http.Error(w, fmt.Sprintf("Cannot add event ID: %s", id), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Added new event to database ID:%d\n", id)

	msg := contracts.EventCreatedEvent{
		ID:         hex.EncodeToString(id),
		Name:       event.Name,
		LocationID: event.Location.ID.String(),
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
	}
	err = eh.EventEmitter.Emit(&msg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot emit creation of event ID: %s", id), http.StatusInternalServerError)
		return
	}
	fmt.Print("Creation of event successfully emitted\n")
}
