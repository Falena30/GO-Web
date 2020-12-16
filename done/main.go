package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type M = map[string]interface{}

func main() {
	/*


		fmt.Println("server start at localhost:8080")
		http.ListenAndServe(":8080", nil)
	*/

	//stripPrefix digunakan untuk mengilangkan static dan hanya akan memanggil sesudahnya
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("asset"))))

	//ini tidak akan muncul karena di index.html di define dengan index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//masih kosong untuk sekarang
		var filePath = filepath.Join("view", "index.html")
		var tmpl, err = template.ParseFiles(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var data = map[string]interface{}{
			"tittle": "Learning Golang Web",
			"name":   "Dimas Adi",
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})
	//templete render partial
	/*
		var tmpl, err = template.ParseGlob("view/*")
		if err != nil {
			panic(err.Error())
			return
		}
	*/

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		var iData = M{"name": "Dimas Adi"}
		var tmplInside = template.Must(template.ParseFiles(
			"view/index.html",
			"view/_header.html",
			"view/_message.html",
		))
		err := tmplInside.ExecuteTemplate(w, "index", iData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		var aData = M{"name": "Adi Suyikno"}
		var tmplInside = template.Must(template.ParseFiles(
			"view/about.html",
			"view/_header.html",
			"view/_message.html",
		))
		err := tmplInside.ExecuteTemplate(w, "about", aData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("server start at localhost:8080")
	http.ListenAndServe(":8080", nil)
}
