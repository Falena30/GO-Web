package main

import (
	"fmt"
	"html/template"
	"net/http"
)

//catatan jika struct harus public
type Info struct {
	Affilation string
	Address    string
}

//GetAffiliationDetailInfo method untuk akses pada bagian html
//ingat harus berada diatas main atau pisahkan filenya
func (t Info) GetAffiliationDetailInfo() string {
	return "have 31 divisions"
}

type Person struct {
	Name    string
	Gender  string
	Hobbies []string
	Info    Info
}

func (p Person) SayHello(from string, message string) string {
	return fmt.Sprintf("%s said: \"%s\"", from, message)
}

var funcMap = template.FuncMap{
	"unescaped": func(s string) template.HTML {
		return template.HTML(s)
	},
	"avg": func(n ...int) int {
		var total = 0
		for _, each := range n {
			total += each
		}
		return total / len(n)
	},
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data = Person{
			Name:    "Dimas Adi",
			Gender:  "Male",
			Hobbies: []string{"Basket", "Programing", "Reding"},
			Info:    Info{"Aious Production", "Bakaran Wetan"},
		}
		//jika berada di sub folder jangan lupa pathnya
		var tmpl = template.Must(template.ParseFiles("view/view.html"))
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/newView", func(w http.ResponseWriter, r *http.Request) {
		var templ = template.Must(template.New("newView.html").
			Funcs(funcMap).ParseFiles("view/newView.html"))
		if err := templ.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/spesific", func(rw http.ResponseWriter, r *http.Request) {
		var tmple = template.Must(template.New("spesific").ParseFiles("view/spesificView.html"))
		if err := tmple.Execute(rw, nil); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/Stest", func(w http.ResponseWriter, r *http.Request) {
		var tmpl = template.Must(template.New("Stest").ParseFiles("view/spesificView.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println("server start at localhost:8080")
	http.ListenAndServe(":8080", nil)
}
