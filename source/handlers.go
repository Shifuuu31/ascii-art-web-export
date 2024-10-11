package source

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type (
	Data struct {
		Input      string
		Banner     string
		Result     string
		FileLength int
	}

	Gerror struct {
		ErrorMessage string
		StatusCode   int
	}
)

var (
	UserData Data
	Template = template.Must(template.ParseGlob("./html/*.html"))
	Error    Gerror
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		CheckError(w, http.StatusNotFound, "Page Not Found :(")
		return
	}

	if r.Method != "GET" {
		CheckError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	if err := Template.ExecuteTemplate(w, "index.html", UserData); err != nil {
		CheckError(w, http.StatusInternalServerError, "Error Executing Template")
		return
	}
}

func HandleAscii(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art-web" {
		CheckError(w, http.StatusNotFound, "Page Not Found :(")
		return
	}

	if r.Method != "POST" {
		CheckError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	err := r.ParseForm()
	if err != nil {
		CheckError(w, http.StatusInternalServerError, fmt.Sprintln(err))
		return
	}

	templDir := "./templates/"
	UserData.Input, UserData.Banner = r.FormValue("input"), r.FormValue("banner")
	parsedFont := ParseFont(w, ReadFontFile(w, templDir+UserData.Banner+".txt"), UserData.Banner)
	UserData.Result = GenerateAsciiArt(UserData.Input, parsedFont)

	http.Redirect(w, r, "/", http.StatusFound)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/css/", http.FileServer(http.Dir("css")))
	_, err := os.Stat("." + r.URL.Path)
	if strings.HasSuffix(r.URL.Path, "/") || err != nil {
		CheckError(w, http.StatusForbidden, "Forbidden: Access Denied")
		return
	}
	fs.ServeHTTP(w, r)
}

func Download(w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/download/", http.FileServer(http.Dir("download")))
	_, err := os.Stat("." + r.URL.Path)
	if strings.HasSuffix(r.URL.Path, "/") || err != nil {
		CheckError(w, http.StatusForbidden, "Forbidden: Access Denied")
		return
	}
	fs.ServeHTTP(w, r)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	filePath := "download/output."
	format := r.URL.Query().Get("format")

	filePath += format

	file, err := os.Open(filePath)
	if err != nil {
		CheckError(w, http.StatusNotFound, "File not found.")
		return
	}
	defer file.Close()

	switch format {
	case "txt":
		if err := GenerateTextFile(UserData.Result); err != nil {
			CheckError(w, http.StatusInternalServerError, fmt.Sprintln(err))
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+UserData.Input+".txt")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", UserData.FileLength))

		http.ServeFile(w, r, filePath)

	case "html":
		if err := GenerateHTMLFile(UserData.Result); err != nil {
			CheckError(w, http.StatusInternalServerError, fmt.Sprintln(err))
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+UserData.Input+".html")
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", UserData.FileLength))

		http.ServeFile(w, r, filePath)

	default:

		CheckError(w, http.StatusBadRequest, "Invalid file format")

	}
}
