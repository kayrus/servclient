package shared

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const logFormat = "%s [%s] \"%s\" %d %d %.10f\n"

type JsonDocument struct {
	sync.RWMutex
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type LogRecord struct {
	http.ResponseWriter

	Ip                    string
	Method, Uri, Protocol string
	Status                int
	Time                  time.Time
	ResponseBytes         int64
	ElapsedTime           time.Duration
}

func HandleSignal() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		fmt.Fprintf(os.Stderr, "Caught %d signal, exiting\n", s)
		os.Exit(0)
	}()
}

func NewLogRecord(w http.ResponseWriter, r *http.Request) *LogRecord {
	return &LogRecord{
		ResponseWriter: w,
		Ip:             r.RemoteAddr,
		Time:           time.Now(),
		Method:         r.Method,
		Uri:            r.RequestURI,
		Protocol:       r.Proto,
		Status:         http.StatusOK,
		ElapsedTime:    time.Duration(0),
		ResponseBytes:  0,
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	record := NewLogRecord(w, r)

	switch r.Method {
	case "GET":
		written, err := w.Write([]byte("ok"))
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

func (r *LogRecord) Log() {
	finishTime := time.Now()
	finishTimeT := finishTime.UTC()
	r.ElapsedTime = finishTime.Sub(r.Time)

	clientIP := r.Ip
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	timeFormatted := finishTimeT.Format("02/Jan/2006 03:04:05")
	requestLine := fmt.Sprintf("%s %s %s", r.Method, r.Uri, r.Protocol)
	fmt.Printf(logFormat, clientIP, timeFormatted, requestLine, r.Status, r.ResponseBytes,
		r.ElapsedTime.Seconds())
}
