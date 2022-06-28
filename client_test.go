package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func setup() *httptest.Server {
	rand.Seed(time.Now().UnixNano())
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arg := r.URL.Query().Get("pct")
		if arg == "" {
			arg = "50"
		}

		p, err := strconv.Atoi(arg)
		if err != nil {
			p = 50
		}

		if r := rand.Intn(100); r < p {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "ERROR")
		}
	}))
}

func TestClient(t *testing.T) {
	ts := setup()
	defer ts.Close()

	resp, err := DoRequest(context.Background(), ts.URL+"?pct=10")
	if err := err; err != nil {
		t.Error(err)
	}
	if resp != "OK" {
		t.Error("expected OK, got", resp)
	}
}
