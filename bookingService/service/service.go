package service

import (
	"flag"
	"github.com/agelloz/reach/bookingService/configuration"
	"github.com/agelloz/reach/bookingService/listener"
	"github.com/agelloz/reach/bookingService/persistence"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type BookingServiceHandler struct {
	DBHandler   persistence.DBHandler
	Endpoint    string
	TLSEndpoint string
}

// ServeAPI is
func ServeAPI() (chan error, chan error) {
	confPath := flag.String("conf", `.\configuration\config.json`,
		"flag to set the path to the configuration json file")
	flag.Parse()
	conf, err := configuration.ExtractConfiguration(*confPath)
	if err != nil {
		panic(err)
	}
	log.Println("connecting to database...")
	dh, err := persistence.NewPersistenceLayer(conf.DBType, conf.DBConnection)
	if err != nil {
		log.Panic(err)
	}
	conn, err := amqp.Dial(conf.AMQPMessageBroker)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil
	}
	defer ch.Close()
	_, err = ch.QueueDeclare("events_queue", false, false, false, false, nil)
	if err != nil {
		log.Panic(err)
	}
	err = listener.Listen(conn, dh)
	if err != nil {
		log.Panic(err)
	}

	eh := &BookingServiceHandler{
		DBHandler:   dh,
		Endpoint:    conf.Endpoint,
		TLSEndpoint: conf.TLSEndpoint,
	}

	r := mux.NewRouter()
	s := r.PathPrefix("/events").Subrouter()
	s.Methods("GET").Path("").HandlerFunc(eh.GetAllEventsHandler)
	s.Methods("DELETE").Path("").HandlerFunc(eh.DeleteAllEventsHandler)
	s.Methods("GET").Path("/bookings").HandlerFunc(eh.GetAllBookingsHandler)
	s.Methods("POST").Path("/bookings/{eventID}").HandlerFunc(eh.AddBookingHandler)
	s.Methods("DELETE").Path("/bookings/{eventID}").HandlerFunc(eh.DeleteBookingHandler)
	httpErrChan := make(chan error)
	httpsErrChan := make(chan error)
	log.Println("bookingService listening...")
	go func() {
		httpsErrChan <- http.ListenAndServeTLS(eh.TLSEndpoint, "certificate/cert.pem", "certificate/key.pem", r)
	}()
	go func() {
		httpErrChan <- http.ListenAndServe(eh.Endpoint, r)
	}()
	return httpErrChan, httpsErrChan
}
