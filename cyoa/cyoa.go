/* Plans
 * Step1: Open json file and get its content
 * Step2: Build a web server
 * Step3: Modify the server to display the content
 */
package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type StoryHandler struct {
	tpl   *template.Template
	story map[string]Chapter
}

func (h *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]

	if chapter, ok := h.story[path]; ok {
		err := h.tpl.Execute(w, chapter)
		if err != nil {
			http.Error(w, "加载失败", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "没有此章节...", http.StatusNotFound)
	}
}

func main() {
	f, _ := os.Open("gopher.json")
	var storyData map[string]Chapter
	json.NewDecoder(f).Decode(&storyData)

	t := template.Must(template.ParseFiles("template.html"))

	h := &StoryHandler{
		tpl:   t,
		story: storyData,
	}

	http.Handle("/", h)
	http.ListenAndServe(":8000", nil)
}
