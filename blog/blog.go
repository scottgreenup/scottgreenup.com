package blog

import (
    "bufio"
    "os"
    "regexp"
    "strconv"
    "strings"
)

const (
    NONE = iota
    TAG_HEADER1
    TAG_HEADER2
    TAG_HEADER3
    TAG_HEADER4
    TAG_HEADER5
    TAG_HEADER6
    TAG_PARAPGRAPH
    TAG_UNORDERED_LIST
    TAG_ORDERED_LIST
    TAG_CODE
    TAG_PARAGRAPH
    TAG_ANCHOR
    TAG_META
)

type MetaData struct {
    Timestamp uint64
    Title string
}

func ParseHTML(lines []string) ([]string, MetaData) {
    first := 0
    length := len(lines)
    var markup []string
    var metadata MetaData

    for i := 0; i < length; i++ {
        curr := getBlockSemantic(lines[i])
        next := getBlockSemantic(lines[i+1])

        if curr != next || i == (length-2) {
            if i == (length-2) {
                i++;
            }
            switch curr {
            case TAG_HEADER6:
                markup = append(markup, "<h6>" + lines[i][7:] + "</h6>")

            case TAG_HEADER5:
                markup = append(markup, "<h5>" + lines[i][6:] + "</h5>")

            case TAG_HEADER4:
                markup = append(markup, "<h4>" + lines[i][5:] + "</h4>")

            case TAG_HEADER3:
                markup = append(markup, "<h3>" + lines[i][4:] + "</h3>")

            case TAG_HEADER2:
                markup = append(markup, "<h2>" + lines[i][3:] + "</h2>")

            case TAG_HEADER1:
                markup = append(markup, "<h1>" + lines[i][2:] + "</h1>")

            case TAG_META:
                if strings.HasPrefix(lines[i], "::title:") {
                    metadata.Title = lines[i][8:]
                }
                if strings.HasPrefix(lines[i], "::timestamp:") {
                    timestamp, _ := strconv.Atoi(lines[i][12:])
                    metadata.Timestamp = uint64(timestamp)
                }

            case TAG_PARAGRAPH:
                markup = append(markup, "<p>")
                for j := first; j <= i; j++ {
                    markup = append(markup, lines[j])
                }
                markup = append(markup, "</p>")

            case TAG_UNORDERED_LIST:
                markup = append(markup, "<ul>")
                for j := first; j <= i; j++ {
                    markup = append(markup, "<li>" + lines[j][3:] + "</li>")
                }
                markup = append(markup, "</ul>")

            case TAG_ORDERED_LIST:
                markup = append(markup, "<ol>")
                for j := first; j <= i; j++ {
                    // TODO fix this for numbers greater than a single digit
                    markup = append(markup, "<li>" + lines[j][4:] + "</li>")
                }
                markup = append(markup, "</ol>")

            case TAG_CODE:
                markup = append(markup, "<pre><code>")
                r := regexp.MustCompile(`\w`);
                index := r.FindIndex([]byte(lines[first]));
                if index == nil {
                    continue
                }
                for j := first; j <= i; j++ {
                    markup = append(markup, lines[j][index[0]:] + "\n")
                }
                markup = append(markup, "</code></pre>")
            }

            first = i + 1
        }
    }

    return markup, metadata
}

func getBlockSemantic(str string) int {
    if strings.HasPrefix(str, "::") {
        return TAG_META
    }

    if strings.HasPrefix(str, "###### ") {
        return TAG_HEADER6
    }
    if strings.HasPrefix(str, "##### ") {
        return TAG_HEADER5
    }
    if strings.HasPrefix(str, "#### ") {
        return TAG_HEADER4
    }
    if strings.HasPrefix(str, "### ") {
        return TAG_HEADER3
    }
    if strings.HasPrefix(str, "## ") {
        return TAG_HEADER2
    }
    if strings.HasPrefix(str, "# ") {
        return TAG_HEADER1
    }
    if r := regexp.MustCompile(`^ [*\-+] `); r.FindIndex([]byte(str)) != nil {
        return TAG_UNORDERED_LIST
    }
    if r := regexp.MustCompile(`^ [\d]+. `); r.FindIndex([]byte(str)) != nil {
        return TAG_ORDERED_LIST
    }
    if r := regexp.MustCompile(`^\t\t`); r.FindIndex([]byte(str)) != nil {
        return TAG_CODE
    }

    r := regexp.MustCompile(`^[\s]*$`)
    index := r.FindIndex([]byte(str))
    if index != nil && len(str) == (index[1] - index[0]) {
        return NONE
    }

    return TAG_PARAGRAPH
}

func getSpanSemantic(str string) int {
    if r := regexp.MustCompile(`\[[\w ]+\]([\w\:\/\.\?\=\#]+)`); r.FindIndex([]byte(str)) != nil {
        return TAG_ANCHOR
    }

    return NONE
}

func ParseHTMLFromFile(filename string) ([]string, MetaData, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, MetaData{}, err
    }
    defer file.Close()

    var lines []string

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    html, meta := ParseHTML(lines)

    return html, meta, nil
}

