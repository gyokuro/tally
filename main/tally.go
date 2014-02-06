package main

import (
	"flag"
	"github.com/gyokuro/tally"
	"github.com/gyokuro/tally/impl"
//	webapp "github.com/gyokuro/tally/resources/webapp"
	"io"
	"log"
	"net/http"
	"os"
	"errors"
	"strconv"
)

// Flags from the command line
var (
	httpPort             = flag.Int("p", 8080, "http server port")
	webappPort           = flag.Int("wp", 8888, "web ui server port")
	noMongo              = flag.Bool("nomgo", false, "True to run without mongo db")
	mongoUrl             = flag.String("dbUrl", "localhost", "MongoDb url")
	mongoDbName          = flag.String("dbName", "tally", "MongoDb database name")
	mongoCollection      = flag.String("dbColl", "cabs", "MongoDb collection name")
	currentWorkingDir, _ = os.Getwd()
)

type fileSystemWrapper int

// Implements the http.FileSystem interface and try to open a local file.  If not found,
// defer to embedded
func (f *fileSystemWrapper) Open(path string) (file http.File, err error) {
	if file, err = http.Dir(currentWorkingDir + "/webapp").Open(path); err == nil {
		return
	}
	return nil, errors.New("Not found") //webapp.Dir(".").Open(path)
}

// Starts a separate server for the web ui.
func startWebUi(port int) {
	http.Handle("/", http.FileServer(new(fileSystemWrapper)))
	webappListen := ":" + strconv.Itoa(port)
	go func() {
		err := http.ListenAndServe(webappListen, nil)
		if err != nil {
			panic(err)
		}
	}()
}

func main() {

	flag.Parse()

	shutdownc := make(chan io.Closer, 1)
	go tally.HandleSignals(shutdownc)

	// Uses the mongodb as backend datastore.
	var service tally.CabService
	if *noMongo {
		service = impl.NewSimpleCabService()
		log.Println("Runing without MongoDb. Using simple / in memory service.")
	} else {
		var err error
		service, err = impl.NewMongoDbCabService(*mongoUrl, *mongoDbName, *mongoCollection)
		if err != nil {
			panic(err)
		}
	}

	httpServer := tally.HttpServer(service)
	httpServer.Addr = ":" + strconv.Itoa(*httpPort)

	// Run the http server in a separate go routine
	// When stopping, send a true to the httpDone channel.
	// The channel done is used for getting notification on clean server shutdown.
	httpDone := make(chan bool)
	done := tally.RunServer(httpServer, httpDone)
	log.Println("Server listening on", *httpPort)

	// Start the UI server
	startWebUi(*webappPort)
	log.Println("Web UI Server listening on", *webappPort)

	// Here is a list of shutdown hooks to execute when receiving the OS signal
	shutdownc <- tally.ShutdownSequence{
		tally.ShutdownHook(func() error {
			// Clean up database connections
			service.Close()
			return nil
		}),
		tally.ShutdownHook(func() error {
			httpDone <- true
			return nil
		}),
	}

	<-done // This just blocks until a bool is sent on the channel
}
