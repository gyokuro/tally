package tally

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

// Returns a http server from given service object
// Registration of URL routes to handler functions that will invoke the service's methods to do CRUD.
func HttpServer(service CabService) *http.Server {
	router := mux.NewRouter()

	// Create / Update Request
	router.Methods("PUT", "POST").Path("/cabs/{cabId}").HandlerFunc(handleCreateUpdate(service))
	// Get Request
	router.Methods("GET").Path("/cabs/{cabId}").HandlerFunc(handleGet(service))
	// Query
	router.Methods("GET").Path("/cabs").HandlerFunc(handleQuery(service))
	// Destroy Request
	router.Methods("DELETE").Path("/cabs/{cabId}").HandlerFunc(handleDelete(service))
	// Destroy All Request
	router.Methods("DELETE").Path("/cabs").HandlerFunc(handleDeleteAll(service))

	// These are added as hacks to work around issues with browsers problems with PUT and DELETE
	// Destroy Request
	router.Methods("POST").Path("/delete/{cabId}").HandlerFunc(handleDelete(service))
	// Destroy All Request
	router.Methods("POST").Path("/deleteAll").HandlerFunc(handleDeleteAll(service))

	return &http.Server{
		Handler: router,
	}
}

// Runs the http server.  This server offers more control than the standard go's default http server
// in that when a 'true' is sent to the stop channel, the listener is closed to force a clean shutdown.
func RunServer(server *http.Server, stop chan bool) (stopped chan bool) {
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}
	stopped = make(chan bool)

	// This will be set to true if a shutdown signal is received. This allows us to detect
	// if the server stop is intentional or due to some error.
	fromSignal := false

	// The main goroutine where the server listens on the network connection
	go func(fromSignal *bool) {
		// Serve will block until an error (e.g. from shutdown, closed connection) occurs.
		err := server.Serve(listener)
		if !*fromSignal {
			log.Println("Warning: server stops due to error", err)
		}
		stopped <- true
	}(&fromSignal)

	// Another goroutine that listens for signal to close the network connection
	// on shutdown.  This will cause the server.Serve() to return.
	go func(fromSignal *bool) {
		select {
		case <-stop:
			listener.Close()
			*fromSignal = true // Intentially stopped from signal
			return
		}
	}(&fromSignal)
	return
}

// Adds all the basic headers such as json content-type for REST endpoint
// and CORS header for cross-domain resource sharing
// TODO - make the CORS domains specified from the command line to tighten security.
func addHeaders(w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
}

func handleCreateUpdate(service CabService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addHeaders(&w)

		params := mux.Vars(r)

		cabId, err := strconv.ParseUint(params["cabId"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		cab := Cab{}
		err = json.Unmarshal(body, &cab)

		// A quick check to make sure we have id that matches
		if cab.Id == Id(0) {
			cab.Id = Id(cabId) // fill in the missing Id from the URL
		}

		if cab.Id != Id(cabId) {
			http.Error(w, "Cab Id and URL mismatch", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = service.Upsert(cab)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleGet(service CabService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addHeaders(&w)

		params := mux.Vars(r)
		cabId, err := strconv.ParseUint(params["cabId"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		cab, err := service.Read(Id(cabId))
		switch err {
		case nil:
			if jsonStr, err2 := json.Marshal(cab); err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Write(jsonStr)
				return
			}

		case ErrorNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleQuery(service CabService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addHeaders(&w)

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		radius, err := strconv.ParseFloat(r.FormValue("radius"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		limit := uint64(8)
		if len(r.FormValue("limit")) > 0 {
			limit, _ = strconv.ParseUint(r.FormValue("limit"), 10, 64)
		}

		cabs, err := service.Query(GeoWithin{
			Center: Location{
				Longitude: longitude,
				Latitude:  latitude,
			},
			Radius: radius,
			Unit:   Meters,
			Limit:  int(limit)})
		switch err {
		case nil:
			if jsonStr, err2 := json.Marshal(cabs); err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Write(jsonStr)
				return
			}

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handleDelete(service CabService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addHeaders(&w)

		params := mux.Vars(r)
		cabId, err := strconv.ParseUint(params["cabId"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = service.Delete(Id(cabId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleDeleteAll(service CabService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addHeaders(&w)

		err := service.DeleteAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
