/*
webhook-server는 kubernetes의 admission webhook을 받아 화면에 출력하는 서버입니다.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"webhook-server/pkg/admission"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
)

func main() {
	setLogger()

	var cert, key string

	// set TLS
	if os.Getenv("TLS") == "" || strings.ToLower(os.Getenv("TLS")) == "false" {
		cert = "/etc/admission-webhook/tls/tls.crt"
		key = "/etc/admission-webhook/tls/tls.key"
	}

	http.HandleFunc("/validate-pods", ServeValidatePods)
	http.HandleFunc("/mutate-pods", ServeMutatePods)
	http.HandleFunc("/health", ServeHealth)

	fmt.Printf("Listening on port 443\n")
	http.ListenAndServeTLS(":443", cert, key, nil)
}

func setLogger() {
	logrus.SetLevel(logrus.DebugLevel)

	lev := os.Getenv("LOG_LEVEL")

	if lev != "" {
		llev, err := logrus.ParseLevel(lev)
		if err != nil {
			logrus.Fatalf("Failed to parse log level: %v", err)
		}
		logrus.SetLevel(llev)
	}
}

func ServeHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ServeValidatePods(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{"uri": "r.RequestURI", "kind": "validate"})
	logger.Debug("Validating pod")

	in, err := parseRequset(r) // admission.Request
	if err != nil {
		logrus.Errorf("Failed to parse admission request: %v", err)
		http.Error(w, fmt.Sprintf("Failed to parse admission request: %v", err), http.StatusBadRequest)
		return
	}

	// TODO: implement validation
	_ = admission.Admitter{
		Logger:  logger,
		Request: in.Request,
	}
	// adm.MutatePodReview()
}

func ServeMutatePods(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(logrus.Fields{"uri": r.RequestURI, "kind": "mutate"})
	in, err := parseRequset(r)
	if err != nil {
		logrus.Errorf("Failed to parse admission request: %v", err)
		http.Error(w, fmt.Sprintf("Failed to parse admission request: %v", err), http.StatusBadRequest)
		return
	}
	_ = admission.Admitter{
		Logger:  logger,
		Request: in.Request,
	}
}

// parseRequest extracts an AdmissionReview from an http.Request if possible
func parseRequset(r *http.Request) (*admissionv1.AdmissionReview, error) {
	if r.Header.Get("Conttent-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}
	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	if len(body) == 0 {
		return nil, fmt.Errorf("admission request has empty body")
	}

	var a admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal admission request json body: %v", err)
	}
	if a.Request == nil {
		return nil, fmt.Errorf("admission review request was nil")
	}
	return &a, nil
}
