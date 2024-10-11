package source

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// read font file function
func ReadFontFile(w http.ResponseWriter, banner string) string {
	file, err := os.ReadFile(banner)
	if err != nil {
		CheckError(w, http.StatusInternalServerError, fmt.Sprintln(err))
		return ""
	}
	content := string(file[1:])
	return content
}

// function responsible for parsong the font file
func ParseFont(w http.ResponseWriter, data string, font string) map[rune][]string {
	startChar := ' '
	var blocks []string
	if font == "thinkertoy" {
		data = data[1:]
		blocks = strings.Split(data, "\r\n\r\n")
	} else {
		blocks = strings.Split(data, "\n\n")
	}
	fontMap := make(map[rune][]string)

	for i, block := range blocks {
		var lines []string
		if font == "thinkertoy" {
			lines = strings.Split(block, "\r\n")
		} else {
			lines = strings.Split(block, "\n")
		}
		if len(lines) > 0 {
			char := rune(startChar + rune(i))
			fontMap[char] = lines
		} else {
			CheckError(w, http.StatusInternalServerError, fmt.Sprintln("warning: empty or malformed block at index %d", i))
		}

	}

	return fontMap
}
func GenerateTextFile(content string) error {
	file, err := os.Create("download/output.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	UserData.FileLength = len(content)

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func GenerateHTMLFile(content string) error {
	htmlTemplate := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>ASCII Art</title>
    </head>
    <body>
        <pre>{{.}}</pre>
    </body>
    </html>
    `
	file, err := os.Create("download/output.html")
	if err != nil {
		return err
	}
	defer file.Close()
	htmlTemplate = strings.ReplaceAll(htmlTemplate, "{{.}}", content)
	UserData.FileLength = len(htmlTemplate)

	_, err = file.WriteString(htmlTemplate)
	if err != nil {
		return err
	}
	return nil
}


func CheckError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	Error.StatusCode = statusCode
	Error.ErrorMessage = msg
	Template.ExecuteTemplate(w, "error.html", Error)
}
