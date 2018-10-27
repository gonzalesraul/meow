package util

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func processResult(handlerFn func(http.ResponseWriter, *http.Request)) (int, string) {
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFn)
	handler.ServeHTTP(recorder, nil)
	return recorder.Code, strings.TrimSpace(recorder.Body.String())
}

func TestResponseOk(t *testing.T) {
	type args struct {
		body interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"ParseMapToJson", args{map[string]bool{"camigol": true}}, `{"camigol":true}`},
		{"EchoStringMessage", args{`{"camigol": true}`}, `"{\"camigol\": true}"`},
		{"EmptyObject", args{`{}`}, `"{}"`},
		{"EmptyString", args{``}, `""`},
		{"NilCheck", args{nil}, `null`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := processResult(func(w http.ResponseWriter, req *http.Request) {
				ResponseOk(w, tt.args.body)
			})
			if got != tt.want {
				t.Errorf("TestResponseOk() = got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	type args struct {
		httpCode int
		message  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"EchoMessage", args{404, "not found"}, `{"error":"not found"}`},
		{"EmptyString", args{500, ""}, `{"error":""}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpCode, got := processResult(func(w http.ResponseWriter, req *http.Request) {
				ResponseError(w, tt.args.httpCode, tt.args.message)
			})
			if httpCode != tt.args.httpCode {
				t.Errorf("TestResponseOk() = got http code %v, want %v", httpCode, tt.args.httpCode)
			} else if got != tt.want {
				t.Errorf("TestResponseOk() = got %v, want %v", got, tt.want)
			}
		})
	}
}
