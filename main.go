package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var messageChannel chan string;

func prepareHeaderforSSE(responseWriter http.ResponseWriter){
	//!must have
	responseWriter.Header().Set("Content-Type", "text/event-stream");
	//optional
	responseWriter.Header().Set("Cache-Control", "no-cache");
	//!must have
	responseWriter.Header().Set("Connection","keep-alive");
	//optional / configure as needed
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*");
}

func loopEvents(responseWriter http.ResponseWriter, request *http.Request){

	prepareHeaderforSSE(responseWriter);

	flusher, ok := responseWriter.(http.Flusher);

	if !ok {
		http.Error(responseWriter,
			"No streaming for you! Renew sub! and we don't support stream",
			http.StatusInternalServerError,
			);
		return;
	}

	for i := 1000; i != 0; i-- {
		//* Must start with 'data:' to send the data as plain text
		//*Could be raw bytes or JSON
		//! the double \n\n signifies: End Of File
		fmt.Fprintf(responseWriter, "data: %d \n\n", i);
		flusher.Flush();
		time.Sleep(1 * time.Second)
	}
}

func getTime(responseWriter http.ResponseWriter, request *http.Request){
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*");

	if messageChannel != nil {
		message := time.Now().Format("15:04:05");
		messageChannel <- message;
	}
}

func sseHandler(responseWriter http.ResponseWriter, request *http.Request){
	prepareHeaderforSSE(responseWriter);

	messageChannel := make(chan string);

	defer func(){
		close(messageChannel);
		messageChannel = nil;
		log.Println("client closed connected");
	}();

	flusher, ok := responseWriter.(http.Flusher);

	if !ok {
		http.Error(responseWriter,
			"No streaming for you! Renew sub! and we don't support stream",
			http.StatusInternalServerError,
			);
		return;
	}

	free: for {
		select {
			case message := <-messageChannel:
				fmt.Fprintf(responseWriter, "data: %s \n\n", message);
				flusher.Flush();
			case <-request.Context().Done():
				break free;
		}
	}
}

//SSE => UTF-8 text only, so JSON is ok
//6 connections per browser, including other tabs
// Push notifications

func main(){

	router := http.NewServeMux(); //unnecessary //if none specified it goes default

	router.HandleFunc("/loop", loopEvents);
	router.HandleFunc("/stream", sseHandler);
	router.HandleFunc("/time", getTime);

	log.Fatal(http.ListenAndServe(":4000", router));

}