package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//setting strings
const TimeFormat = time.RFC3339

//error strings
const (
	StaticFilePathNotSetError        = "STATIC_FILES_PATH not found \n using default"
	DatabaseConnectionUrlNotSetError = "POSTGRES_CONNECTION_STRING not found \n using default"
	ListenAddressNotSetError         = "LISTEN_ADDRESS not found \n using default"
)

var SessionServiceChannel = make(chan bool)
var SnapShotServiceChannel = make(chan bool)

// Asset Strings
const (
	StdKeyLength          = 512
	FilePathVarName       = "STATIC_FILES_PATH"
	ListeAddressVarName   = "LISTEN_ADDRESS"
	DatabaseConnectionVar = "POSTGRES_CONNECTION_STRING"
	DefaultListenAddress  = "0.0.0.0:9745"
	DefaultConnectionURL  = "postgres://postgres:postgres@localhost:5432/pork"
	DefaultPath           = "./static"
	SnapshotFile          = "./snapshot.bin"
)

func main() {

	// Getting settings from to OS to setup the microservice

	staticPath, isVariableSet := os.LookupEnv(FilePathVarName)
	if !isVariableSet {
		fmt.Print(StaticFilePathNotSetError)
		staticPath = DefaultPath
	}

	ConnectionURL, isVariableSet := os.LookupEnv(DatabaseConnectionVar)
	if !isVariableSet {
		fmt.Print(DatabaseConnectionUrlNotSetError)
		ConnectionURL = DefaultConnectionURL
	}

	ListenAddress, isVariableSet := os.LookupEnv(ListeAddressVarName)
	if !isVariableSet {
		fmt.Print(ListenAddressNotSetError)
		ListenAddress = DefaultListenAddress
	}

	// setting up session expire service
	ticker := time.NewTicker(time.Minute * 3)
	sessionExpireService(ticker)

	// setting up snapshot service
	snapShotTicker := time.NewTicker(time.Minute * 5)
	snapShotCartService(snapShotTicker)

	// connecting to database
	connection := connect(ConnectionURL)

	// settig up signals
	sigs := make(chan os.Signal, 1)
	sigDone := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// goroutine to listen to signals
	go func() {
		<-sigs
		sigDone <- true
	}()

	// go routine to gracefully shutdown on recieving kill signal
	go func() {
		<-sigDone                              // unblocks goroutine only when sigDone has value ie a signal
		SessionServiceChannel <- true          // sends a value to session expire session goroutine to stop
		SnapShotServiceChannel <- true         // sends signal to snapshot service
		connection.Close(context.Background()) // closing database connextion
		time.Sleep(time.Second)                // givining time for database connectin to close
		os.Exit(0)                             // exiting
	}()

	http.Handle("/", http.FileServer(http.Dir(staticPath)))
	fmt.Println("Server running now.....")

	log.Fatal(http.ListenAndServe(ListenAddress, nil))
}
