package functions

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Data struct {
	Result string
	banner string
	text   string
	format string
}

var D Data

type ERRORS struct {
	PageTitle string
	Message   string
	ErrCde    int
}

var ERR ERRORS

// The welcome page
func Welcom(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	tmpl2, err2 := template.ParseFiles("./templates/errors.html")
	if err != nil || err2 != nil {
		if err2 != nil {
			ChooseErr(500, w)
			tmpl2.Execute(w, ERR)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if r.URL.Path != "/" {
		ChooseErr(404, w)
		tmpl2.Execute(w, ERR)
		return
	}

	// fmt.Println(r.Method)
	if r.Method != "GET" {
		ChooseErr(405, w)
		tmpl2.Execute(w, ERR)
		return
	}
	tmpl.Execute(w, nil)
}

// The result page
func Last(w http.ResponseWriter, r *http.Request) {
	// r.Method = http.MethodPost

	tmpl2, err2 := template.ParseFiles("./templates/errors.html")
	if err2 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if r.URL.Path != "/ascii-art" {
		ChooseErr(404, w)
		tmpl2.Execute(w, ERR)
		return
	}
	// fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		ChooseErr(405, w)
		tmpl2.Execute(w, ERR)
		return
	}

	check := r.FormValue("check")
	D.text = r.FormValue("ljomla")
	D.banner = r.FormValue("banner")
	// fmt.Println(d.text, d.banner)
	D.text = strings.ReplaceAll(D.text, "\r\n", "\n")
	// fmt.Println(len(D.text),D.banner)
	if len(D.text) > 250 {
		ChooseErr(777, w)
		tmpl2.Execute(w, ERR)
		return
	}
	if D.text == "" || D.banner == "" {
		ChooseErr(400, w)
		tmpl2.Execute(w, ERR)
		return
	}

	D.Result = FS(D.banner, D.text)
	if D.Result == "ERORR" {
		ChooseErr(400, w)
		tmpl2.Execute(w, ERR)
		return
	}

	if check == "sub" {
		tmpl, err := template.ParseFiles("./templates/result.html")
		if err != nil {
			ChooseErr(500, w)
			tmpl2.Execute(w, ERR)
			return
		}
		tmpl.Execute(w, D)
	} else {
		D.format = r.FormValue("format")
		http.Redirect(w, r, "/output", http.StatusSeeOther)
	}
}

func Output(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.ParseFiles("./templates/download.html")
	tmpl2, err2 := template.ParseFiles("./templates/errors.html")

	if err2 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/output" {
		ChooseErr(404, w)
		tmpl2.Execute(w, ERR)
		return
	}

	// fmt.Println("FORMAT:", D.format)
	if D.format == "" {
		ChooseErr(400, w)
		tmpl2.Execute(w, ERR)
		return
	}
	if D.format == ".html" {
		D.Result = "<pre>" + D.Result + "</pre>"
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Disposition", "attachment; filename=file"+D.format)
	} else {
		D.Result = strings.TrimSuffix(strings.TrimPrefix(D.Result, "<pre>"), "</pre>")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=file"+D.format)
	}

	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(D.Result))))
	w.Write([]byte(D.Result))
}

func ServeStyle(w http.ResponseWriter, r *http.Request) {
	tmpl2, err2 := template.ParseFiles("./templates/errors.html")
	if err2 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fs := http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles")))
	if r.URL.Path == "/styles/" {
		ChooseErr(403, w)
		tmpl2.Execute(w, ERR)
		return
	}

	fs.ServeHTTP(w, r)
}

func ChooseErr(code int, w http.ResponseWriter) {
	if code == 404 {
		ERR.PageTitle = "Error 404"
		ERR.Message = "The page web doesn't exist\nError 404"
		ERR.ErrCde = code
		w.WriteHeader(code)
	} else if code == 405 {
		ERR.PageTitle = "Error 405"
		ERR.Message = "The method is not alloweded\nError 405"
		ERR.ErrCde = code
		w.WriteHeader(code)
	} else if code == 400 {
		ERR.PageTitle = "Error 400"
		ERR.Message = "Bad Request\nError 400"
		ERR.ErrCde = code
		w.WriteHeader(code)
	} else if code == 500 {
		ERR.PageTitle = "Error 500"
		ERR.Message = "Internal Server Error\nError 500"
		ERR.ErrCde = code
		w.WriteHeader(code)
	} else if code == 403 {
		ERR.PageTitle = "Error 403"
		ERR.Message = "You didn't authorized to see the resource\nError 403"
		ERR.ErrCde = code
		w.WriteHeader(code)
	}else if code == 777 {
		ERR.PageTitle = "Error 400"
		ERR.Message = "Bad Request: String have more than 250 characters\nError 400"
		ERR.ErrCde = 400
		w.WriteHeader(400)
	}
}
