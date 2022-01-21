package gee

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get %s\n", r.URL.Path)
}

func post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "post %s\n", r.URL.Path)
}

func TestStaticRouter(t *testing.T) {
	server := func() {
		builder := New()
		builder.Get("/test1", get)
		builder.Post("/test1", post)
		fmt.Println("ready to start")
		log.Fatal(builder.Run(":8000"))
		fmt.Println("exit")
	}

	go server()

	t.Run("test for Get", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8000/test1")
		if err != nil {
			t.Errorf("get err %e", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %e", err)
		}
		if string(s) != "get /test1\n" {
			t.Errorf("expect %s, get %s", "get /test1\n", string(s))
		}
	})

	t.Run("test for Post", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Post("http://127.0.0.1:8000/test1", "", nil)
		if err != nil {
			t.Errorf("get err %e", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %e", err)
		}
		if string(s) != "post /test1\n" {
			t.Errorf("expect %s, get %s", "post /test1\n", string(s))
		}
	})

	t.Run("test for invalid URL", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Post("http://127.0.0.1:8000", "", nil)
		if err != nil {
			t.Errorf("get err %e", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %e", err)
		}
		if string(s) != "404 NOT FOUND: /\n" {
			t.Errorf("expect %s, get %s", "404 NOT FOUND: \n", string(s))
		}
	})
}
