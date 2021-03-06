package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
)

var IRCDIR string

func parseURL(url string) string {
	paths := strings.Split(url, "channel/")
	path := paths[0]
	if len(paths) > 1 {
		path = paths[0] + "#" + paths[1]
	}
	return path
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	url := parseURL(r.URL.Path)
	buffer, _ := ioutil.ReadFile(IRCDIR + url + "/out")
	w.Write(buffer)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	url := parseURL(r.URL.Path)
	r.ParseForm()
	msg := r.FormValue("msg")
	f, _ := os.OpenFile(IRCDIR+url+"/in", os.O_RDWR|os.O_APPEND, 0660)
	defer f.Close()

	f.WriteString(msg + "\n")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getHandler(w, r)
	} else if r.Method == "POST" {
		postHandler(w, r)
	}
}

func auth(fn http.HandlerFunc, user string, pass string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checkUser, checkPass, _ := r.BasicAuth()
		if user != "" && pass != "" {
			if user != checkUser && pass != checkPass {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"IRC\"")
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}
		}
		fn(w, r)
	}
}

func main() {
	_auth := flag.String("a", "", "Authentication. Example: user:pass (default none)")
	path := flag.String("i", "", "IRC directory")
	port := flag.String("p", "8003", "Port of the server")
	flag.Parse()

	if *path == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal("Could not get home directory", err)
		}
		IRCDIR = usr.HomeDir + "/irc"
	} else {
		IRCDIR = *path
	}

	user := ""
	pass := ""
	if *_auth != "" {
		info := strings.Split(*_auth, ":")
		user = info[0]
		pass = info[1]
	}

	http.HandleFunc("/", auth(indexHandler, user, pass))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
