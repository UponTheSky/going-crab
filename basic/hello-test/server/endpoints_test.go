package server

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestEndpointWithServer(t *testing.T) {
	// asset
	mux := http.NewServeMux()
	logger := log.New(os.Stderr, "", log.LstdFlags)

	mux.HandleFunc("/capitalize", CapitalizeHandler(logger))

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	input := "Bradley Cooper"
	want := "BRADLEY COOPER"

	// act
	client := testServer.Client()
	body := strings.NewReader(input)

	response, err := client.Post(testServer.URL+"/capitalize", "text/plain", body)

	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	// assert
	if response.StatusCode != http.StatusOK {
		t.Errorf("the server response status code: got=%v, want=%v", response.StatusCode, http.StatusOK)
	}

	respBody := response.Body
	defer respBody.Close()

	gotByte, err := io.ReadAll(respBody)

	if err != nil {
		t.Fatalf("unexpected response read error: %v", err)
	}

	got := string(gotByte)

	if got != want {
		t.Errorf("Test /capitalize endpoint: got=%v, want=%v", got, want)
	}
}

func TestEndpointNoServer(t *testing.T) {
	// asset
	logger := log.New(os.Stderr, "", log.LstdFlags)
	handler := CapitalizeHandler(logger)
	input := "Bradley Cooper"
	want := "BRADLEY COOPER"

	body := strings.NewReader(input)
	request := httptest.NewRequest(http.MethodPost, "/capitalize", body)
	request.Header.Set("Content-Type", "text/plain")

	rw := httptest.NewRecorder()

	// act
	handler(rw, request)

	// assert
	response := rw.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("status code error - want=%v, got=%v", http.StatusOK, response.StatusCode)
		return
	}

	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	got := string(resBody)

	if got != want {
		t.Errorf("different return value - want=%v, got=%v", want, got)
		return
	}
}
