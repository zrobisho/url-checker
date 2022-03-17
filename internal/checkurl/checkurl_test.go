package checkurl

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckURL200(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer svr.Close()
	client := &http.Client{}

	urlChecker := URLChecker{
		Client: client,
	}

	err := urlChecker.CheckURL(svr.URL)
	if err != nil {
		t.Fail()
	}
}

func TestCheckURL500(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer svr.Close()
	client := &http.Client{}

	urlChecker := URLChecker{
		Client: client,
	}

	err := urlChecker.CheckURL(svr.URL)
	if err == nil {
		t.Fail()
	}
}
