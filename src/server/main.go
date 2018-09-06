package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"shared"
)

var doc *shared.JsonDocument = &shared.JsonDocument{
	Id:      1,
	Message: "Hello, World",
}

func docHandler(w http.ResponseWriter, r *http.Request) {
	record := shared.NewLogRecord(w, r)

	switch r.Method {
	case "GET":
		doc.RLock()
		j, _ := json.Marshal(doc)
		doc.RUnlock()
		written, err := w.Write(j)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			record.ResponseBytes += int64(written)
		}
	case "POST":
		d := json.NewDecoder(r.Body)
		n := &shared.JsonDocument{}
		err := d.Decode(n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			record.Status = http.StatusInternalServerError
		} else {
			doc.Lock()
			doc.Id = n.Id
			doc.Message = n.Message
			doc.Unlock()
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		record.Status = http.StatusMethodNotAllowed
		fmt.Fprintf(w, "Method is not allowed")
	}

	record.Log()
}

func main() {
	shared.HandleSignal()

	http.HandleFunc("/", docHandler)
	http.HandleFunc("/healthz", shared.HealthHandler)

	fmt.Fprintln(os.Stderr, "Go!")
	http.ListenAndServe(":8080", nil)
}
