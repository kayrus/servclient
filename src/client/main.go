package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"shared"
)

var baseURL string = "http://localhost:8080"

func postJsonDocument(doc *shared.JsonDocument) error {
	j, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	_, err = doRequest(req)
	return err
}

func doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

func getJsonDocument() (*shared.JsonDocument, error) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := doRequest(req)
	if err != nil {
		return nil, err
	}
	var data shared.JsonDocument
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func call() (string, error) {
	t, err := getJsonDocument()

	if err != nil {
		return "", err
	}

	t.Message = reverse(t.Message)
	t.Id++
	/*
		err = postJsonDocument(t)

		if err != nil {
			return "", err
		}

		t, err = getJsonDocument()

		if err != nil {
			return "", err
		}

		j, err := json.MarshalIndent(t, "", "\t")

		if err != nil {
			panic(err)
		}
	*/
	return t.Message, err
}

func reverseProxy(w http.ResponseWriter, r *http.Request) {
	record := shared.NewLogRecord(w, r)

	switch r.Method {
	case "GET":
		j, err := call()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			record.Status = http.StatusServiceUnavailable
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		written, err := w.Write([]byte(j))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			record.ResponseBytes += int64(written)
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

	fmt.Fprintln(os.Stderr, "Go!")

	url := os.Getenv("APP_URL")
	if url != "" {
		baseURL = url
	}

	http.HandleFunc("/", reverseProxy)
	http.HandleFunc("/healthz", shared.HealthHandler)
	http.ListenAndServe(":8081", nil)
}
