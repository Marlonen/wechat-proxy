package wechat

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMessageServer(t *testing.T) {

	ts_get := []struct {
		Url    string
		Calls  []string
		Result string
	}{
		{
			Url:    "/msg?echostr=test",
			Calls:  []string{"/svc1"},
			Result: "test",
		},
		{
			Url:    "/msg?echostr=test",
			Calls:  []string{"/svc1", "/svc2"},
			Result: "test",
		},
	}

	ts_post := []struct {
		Url    string
		Calls  []string
		Body   string
		Result string
	}{
		{
			Url:    "/msg?",
			Calls:  []string{"/svc1"},
			Body:   "<xml>...</xml>",
			Result: "<xml>...</xml>",
		},
		{
			Url:    "/msg?",
			Calls:  []string{"/svc2"},
			Body:   "<xml>...</xml>",
			Result: "",
		},
		{
			Url:    "/msg?",
			Calls:  []string{"/svc1", "/svc2"},
			Body:   "<xml>...</xml>",
			Result: "<xml>...</xml>",
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/msg", NewMessageServer())
	mux.HandleFunc("/svc1", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		echostr := r.Form.Get("echostr")
		if echostr != "" {
			w.Write([]byte(echostr))
			return
		}

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(body)
	})
	mux.HandleFunc("/svc2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, v := range ts_get {
		url := ts.URL + v.Url
		for _, c := range v.Calls {
			url += "&call=" + ts.URL[7:] + c
		}
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if string(body) != v.Result {
			log.Printf("body: %s, result: %s\n", string(body), v.Result)
			t.Fatal()
		}
	}

	for _, v := range ts_post {
		url := ts.URL + v.Url
		for _, c := range v.Calls {
			url += "&call=" + ts.URL[7:] + c
		}
		resp, err := http.Post(url, "", strings.NewReader(v.Body))
		if err != nil {
			log.Fatal()
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if string(body) != v.Result {
			t.Fatal()
		}
	}
}
