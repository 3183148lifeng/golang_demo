package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//xiaorui.cc
const AddForm = `
<html><body>
<form method="POST" action="/add">
Name: <input type="text" name="name">
Age: <input type="text" name="age">
<input type="submit" value="Add">
</form>
</body></html>
`
const setform = `
<html><body>
<form method="POST" action="/set">
key: <input type="text" name="key">
value: <input type="text" name="value">
<input type="submit" value="set">
</form>
</body></html>
`

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path == "favicon.ico" {
		http.NotFound(w, r)
		return
	}
	if path == "" {
		path = "index.html"
	}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}
	fmt.Fprintf(w, "%s\n", contents)
}

func Add(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	age := r.FormValue("age")
	if name == "" || age == "" {
		fmt.Fprint(w, AddForm)
		return
	}
	fmt.Fprintf(w, "Save : Your name is  %s , You age is %s", name, age)
}
func redisset(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")
	if key == "" || value == "" {
		fmt.Fprint(w, setform)
		return
	}
	spec, err := redis.Dial("tcp", ":9527")
	defer spec.Close()
	client, e := redis.NewSynchClientWithSpec(spec)
	if e != nil {
		log.Println("服务器连接有异常", e)
		return
	}
	inva := []byte(value)
	client.Set(key, inva)
	fmt.Fprintf(w, "哥们，你输入的key  %s 和value  %s 已经插入到redis里面了", key, key)
}
func redisget(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" {
		fmt.Fprint(w, setform)
		return
	}
	spec := redis.DefaultSpec().Db(0).Password("")
	client, e := redis.NewSynchClientWithSpec(spec)
	if e != nil {
		log.Println("服务器连接有异常", e)
		return
	}

	value, e := client.Get(key)
	fmt.Fprintf(w, "哥们，你要查询的key  %s 和value  %s ", key, value)
}
func valueget(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	user := params.Get("user")
	fmt.Fprintf(w, "you are get user %s", user)
}

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/redisset", redisset)
	http.HandleFunc("/redisget", redisget)
	http.HandleFunc("/valueget", valueget)
	s := &http.Server{
		Addr:           ":8888",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
