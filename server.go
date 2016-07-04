package main

import (
    "github.com/scottgreenup/scottgreenup.com/blog"
    "github.com/gorilla/mux"

    "bytes"
    "flag"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "sort"
    "strconv"
    "time"
)

var port = flag.Int("port", 80, "The port for the webserver to run on.")

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
    fmt.Printf("%s - index\n", r.RemoteAddr)

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

type ByTimestamp []blog.MetaData
func (b ByTimestamp) Len() int { return len(b) }
func (b ByTimestamp) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByTimestamp) Less(i, j int) bool { return b[i].Timestamp > b[j].Timestamp }

func blogHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("%s - blog\n", r.RemoteAddr)

    // Get all the posts out of the directory
    files, _ := ioutil.ReadDir("./content/posts");
    posts := make(map[blog.MetaData][]string)
    var meta_data []blog.MetaData
    for _, f := range files {
        markup, meta, err := blog.ParseHTMLFromFile("./content/posts/" + f.Name())
        if err != nil {
            log.Printf("Error from ParseHTMLFromFile(): %+v", err);
        }
        posts[meta] = markup;
        meta_data = append(meta_data, meta)
    }

    // Print them to a buffer, inserting HTML appropriately
    sort.Sort(ByTimestamp(meta_data))
    var buf bytes.Buffer
    buf.WriteString("{{define \"blog_content\"}}")
    for k, v := range meta_data {

        // Write the post to the buffer, insert the timestamp after header
        buf.WriteString("<article>")
        markup := posts[v];
        for i := 0; i < len(markup); i++ {
            buf.WriteString(markup[i])

            // TODO move this logic to blog package
            if i == 0 {
                tm := time.Unix(int64(v.Timestamp), 0)
                buf.WriteString(
                    "<h5 id=\"timestamp\">" + tm.Format(time.RFC1123) + "</h5>",
                )
            }
        }
        buf.WriteString("</article>")

        // Insert a divider between posts
        if k != len(meta_data)-1 {
            buf.WriteString("<hr />");
        }
    }
    buf.WriteString("{{end}}")
    templates.Parse(buf.String());

    if r.Method != "GET" {
        return
    }

    if len(r.URL.RawQuery) == 0 {
        err := renderTemplate(w, r, "blog")
        if err != nil {
            log.Println(err.Error())
        }
        return
    }

    err := renderTemplate(w, r, "blog")
    if err != nil {
        log.Println(err.Error())
        http.Error(w, http.StatusText(500), 500)
        return
    }
}

func main() {
    flag.Parse()

    r := mux.NewRouter()
    r.HandleFunc("/blog", blogHandler)
    r.HandleFunc("/blog/", blogHandler)
    r.HandleFunc("/", indexHandler)
    r.PathPrefix("/static/").Handler(
        http.StripPrefix(
            "/static",
            http.FileServer(http.Dir("content/static"))))

    port_string := strconv.Itoa(*port)
    http.ListenAndServe(":" + port_string, r)

}


