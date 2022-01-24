package gee

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestContext(t *testing.T) {

	get := func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	}

	getWithName := func(c *Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	}

	login := func(c *Context) {
		c.JSON(http.StatusOK, H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	}

	server := func() {
		builder := New()
		builder.Get("/", get)
		builder.Get("/hello", getWithName)
		builder.Post("/login", login)
		fmt.Println("ready to start")
		log.Fatal(builder.Run(":8000"))
		fmt.Println("exit")
	}

	// start the server
	go server()

	t.Run("test for Get 0", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8000/")
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "<h1>Hello Gee</h1>" {
			t.Errorf("expect %s, get %s", "<h1>Hello Gee</h1>", string(s))
		}
	})

	t.Run("test for Get 1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8000/hello?name=yyl494")
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "hello yyl494, you're at /hello\n" {
			t.Errorf("expect %s, get %s", "hello yyl494, you're at /hello\n", string(s))
		}
	})

	t.Run("test for Post", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Post("http://127.0.0.1:8000/login", "application/x-www-form-urlencoded",
			strings.NewReader("username=hhh&password=hhh"))
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "{\"password\":\"hhh\",\"username\":\"hhh\"}\n" {
			t.Errorf("expect %s, get %s", "{\"password\":\"hhh\",\"username\":\"hhh\"}\n", string(s))
		}
	})

	t.Run("test for invalid URL", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.PostForm("http://127.0.0.1:8000", url.Values{"user": {"hhh"}, "pass": {"123"}})
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

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/geektutu")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
}

func TestDynamicRouter(t *testing.T) {
	go func() {
		r := New()

		r.Get("/hello/:name", func(c *Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		r.Get("/assets/*filepath", func(c *Context) {
			c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
		})

		log.Fatal(r.Run(":8001"))
	}()

	t.Run("test for Get 1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8001/assets/a/b/c")
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "{\"filepath\":\"a/b/c\"}\n" {
			t.Errorf("expect %s, get %s", "{\"filepath\":\"a/b/c\"}\n", string(s))
		}
	})

	t.Run("test for Get 0", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8001/hello/yyl494")
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "hello yyl494, you're at /hello/yyl494\n" {
			t.Errorf("expect %s, get %s", "hello yyl494, you're at /hello/yyl494\n", string(s))
		}
	})

}

func TestGroup(t *testing.T) {
	go func() {
		r := New()
		r.Get("/index", func(c *Context) {
			c.HTML(http.StatusOK, "<h1>Index Page</h1>")
		})
		v1 := r.Group("/v1")
		{
			v1.Get("/hello", func(c *Context) {
				// expect /hello?name=geektutu
				c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
			})
		}
		v2 := r.Group("/v2")
		{
			v2.Get("/hello/:name", func(c *Context) {
				// expect /hello/geektutu
				c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
			})
		}

		r.Run(":8002")
	}()

	t.Run("test for Get 1", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8002/v1/hello")
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "hello , you're at /v1/hello\n" {
			t.Errorf("expect %s, get %s", "hello , you're at /v1/hello\n", string(s))
		}
	})

	t.Run("test for Get 0", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8002/v2/hello/yyl")
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "hello yyl, you're at /v2/hello/yyl\n" {
			t.Errorf("expect %s, get %s", "hello yyl, you're at /v2/hello/yyl\n", string(s))
		}
	})

}

func TestMiddleware(t *testing.T) {
	go func() {
		r := New()
		r.Use(func(c *Context) {
			log.Println("access to ", c.Req.URL.Path)
		}) // global midlleware
		r.Get("/", func(c *Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v2 := r.Group("/v2")
		v2.Use(func(c *Context) {
			log.Println("v2 access to ", c.Req.URL.Path)
		}) // v2 group middleware
		{
			v2.Get("/hello/:name", func(c *Context) {
				c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
			})
		}

		r.Run(":8003")
	}()

	t.Run("test for middleware", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://127.0.0.1:8003/v2/hello/yyl")
		if err != nil {
			t.Errorf("get err %s", err)
		}

		defer response.Body.Close()
		s, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("get err %s", err)
		}
		if string(s) != "hello yyl, you're at /v2/hello/yyl\n" {
			t.Errorf("expect %s, get %s", "hello yyl, you're at /v2/hello/yyl\n", string(s))
		}
	})
}
