package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"path"
	"path/filepath"
)

type User struct {
	FirstName string
	LastName  string
	Edad      int
	Avatar    string
}

var users map[int]*User
var i int
var templates = template.Must(template.ParseGlob("templates/*"))

func TodoMainHandler(res http.ResponseWriter, req *http.Request) {
	basedir, _ := filepath.Abs(filepath.Dir("."))

	if req.Method == "POST" {
		firstName := req.FormValue("firstName")
		lastName := req.FormValue("lastName")
		edad, _ := strconv.Atoi(req.FormValue("edad"))
		m := req.MultipartForm
		files := m.File["avatar"]

		file, err := files[0].Open()
		defer file.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(path.Join(basedir, "static/images", files[0].Filename))

		defer dst.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		user := &User{firstName, lastName, edad, files[0].Filename}
		users[i] = user
		i++
	}
	err := templates.ExecuteTemplate(res, "base", users)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
func TodoDeleteHandler(res http.ResponseWriter, req *http.Request) {
	index, _ := strconv.Atoi(req.FormValue("i"))
	delete(users, index)
	http.Redirect(res, req, "/", 302)
}
func main() {
	users = make(map[int]*User, 0)
	i = 1
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", TodoMainHandler)
	http.HandleFunc("/eliminar", TodoDeleteHandler)
	http.ListenAndServe(":8000", nil)
}
