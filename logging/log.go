package logging

import (
	"io"
	"log"
	"net/http"
	"os"
)

// Configure sends standard log output to stderr and the configured log file.
func Configure(logFile string) (*os.File, error) {
	if logFile == "" {
		log.SetOutput(os.Stderr)
		return nil, nil
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return nil, err
	}

	log.SetOutput(io.MultiWriter(os.Stderr, file))
	return file, nil
}

// Failure emits the user and failed step in a consistent format.
func Failure(username, failedStep string, err error) {
	log.Printf("user=%s failed_step=%s error=%v", username, failedStep, err)
}

// LoginEvent captures one summarized backend login request log.
type LoginEvent struct {
	Method      string
	URL         string
	Target      string
	RemoteAddr  string
	ContentType string
	Username    string
	Body        string
	Outcome     string
	FailedStep  string
	Error       string
}

// NewLoginEvent builds the default log event for one backend login request.
func NewLoginEvent(r *http.Request, target string) LoginEvent {
	return LoginEvent{
		Method:      r.Method,
		URL:         requestURL(r),
		Target:      target,
		RemoteAddr:  r.RemoteAddr,
		ContentType: r.Header.Get("Content-Type"),
		Username:    "-",
		Outcome:     "failed",
	}
}

// LogLoginEvent emits one compact line for a backend login request.
func LogLoginEvent(event LoginEvent, status int) {
	log.Printf(
		"maas_login method=%s url=%q target=%q remote_addr=%q content_type=%q username=%q body=%q outcome=%s failed_step=%s error=%q status=%d",
		event.Method,
		event.URL,
		event.Target,
		event.RemoteAddr,
		event.ContentType,
		event.Username,
		event.Body,
		event.Outcome,
		event.FailedStep,
		event.Error,
		status,
	)
}

// StatusRecorder records the status code written by a handler.
type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *StatusRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(body)
}

func (r *StatusRecorder) StatusCode() int {
	return r.status
}

func requestURL(r *http.Request) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}
	return scheme + "://" + host + r.URL.RequestURI()
}
