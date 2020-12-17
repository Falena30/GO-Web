package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func routeFormGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("form").ParseFiles("view/FormView.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func routeSubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var tmpl = template.Must(template.New("result").ParseFiles("view/FormView.html"))
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var name = r.FormValue("name")
		var message = r.FormValue("message")

		var data = map[string]string{"name": name, "message": message}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func routeFormFileGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var tmpl = template.Must(template.ParseFiles("view/FormViewFile.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func routeSubmiFiletPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var alias = r.FormValue("alias")

	uploadfile, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadfile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := handler.Filename
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	}
	filelocation := filepath.Join(dir, "file", filename)
	targetfile, err := os.OpenFile(filelocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetfile.Close()

	if _, err := io.Copy(targetfile, uploadfile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("done"))
}

func handleJSONView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view/jsonView.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		payload := struct {
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Gender string `json:"gender"`
		}{}
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message := fmt.Sprintf(
			"hello, my name is %s. I'm %d year old %s",
			payload.Name,
			payload.Age,
			payload.Gender,
		)
		w.Write([]byte(message))
		return
	}

	http.Error(w, "Only accept POST request", http.StatusBadRequest)
}
func actionResponse(w http.ResponseWriter, r *http.Request) {
	data := []struct {
		Name string
		Age  int
	}{
		{"Dimas Adi", 24},
		{"Aulia Wahib", 22},
		{"Imam Abdul", 23},
		{"Nur Faiz", 22},
	}
	//jsoninbytes, err := json.Marshal(data)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.Write(jsoninbytes)
}

func handleMutli(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("view/multiView.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only accept POST request", http.StatusBadRequest)
		return
	}

	basePath, _ := os.Getwd()
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		fileLocation := filepath.Join(basePath, "file", part.FileName())
		dst, err := os.Create(fileLocation)
		if dst != nil {
			defer dst.Close()
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte(`all files uploaded`))
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

	http.HandleFunc("/testingPostAndGet", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			w.Write([]byte("post"))
		case "GET":
			w.Write([]byte("get"))
		default:
			http.Error(w, "", http.StatusBadRequest)
		}
	})

	http.HandleFunc("/form", routeFormGet)
	http.HandleFunc("/process", routeSubmitPost)
	http.HandleFunc("/formFile", routeFormFileGet)
	http.HandleFunc("/processFile", routeSubmiFiletPost)
	http.HandleFunc("/jsonView", handleJSONView)
	http.HandleFunc("/save", handleSave)
	http.HandleFunc("/jsonResonse", actionResponse)
	http.HandleFunc("/multiView", handleMutli)
	http.HandleFunc("/upload", handleUpload)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("asset"))))
	fmt.Println("server start at localhost:8080")
	http.ListenAndServe(":8080", nil)
}
