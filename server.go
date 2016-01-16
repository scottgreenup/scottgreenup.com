package main

import (
    "html/template"
    "log"
    "net/http"
)

var templates = template.Must(template.ParseGlob("content/template/*"))

type Page struct {
    Title string
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string) error {
    err := templates.ExecuteTemplate(w, name, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
    }

    return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    err := renderTemplate(w, r, "index")
    if err != nil {
        log.Println(err.Error())
        http.Error(w, http.StatusText(500), 500)
        return
    }
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        return
    }

    if len(r.URL.RawQuery) == 0 {
        err := renderTemplate(w, r, "blog")
        if err != nil {
            log.Println(err.Error())
        }
    }

    err := renderTemplate(w, r, "blog")
    if err != nil {
        log.Println(err.Error())
        http.Error(w, http.StatusText(500), 500)
        return
    }
}

func main() {
    // TODO - Remove the redudancy in serving static traffic
    fs := http.FileServer(http.Dir("content/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.Handle("/blog/static/", http.StripPrefix("/blog/static/", fs))

    // Must be ordered in least specific to most specific.
    http.HandleFunc("/blog/", blogHandler)
    http.HandleFunc("/", indexHandler)

    log.Println("Listening...")
    http.ListenAndServe(":8080", nil)
}
