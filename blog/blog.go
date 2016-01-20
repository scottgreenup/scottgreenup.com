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
    if r := regexp.MustCompile(`\[[\w\s']+\]\([\w:/.?=]+\)`); r.FindIndex([]byte(str)) != nil {
        return TAG_ANCHOR
    }

    return NONE
}

func ParseHTML(lines []string) ([]string, MetaData) {
    first := 0
    length := len(lines)
    var markup []string
    var metadata MetaData

    // This is a second parse for inline markdown like anchors
    for i := 0; i < length; i++ {
        curr := getSpanSemantic(lines[i]);
        if curr == NONE {
            continue
        }

        switch curr {
        case TAG_ANCHOR:
            r := regexp.MustCompile(`\[[\w\s']+\]\([\w:/.?=]+\)`);
            m := r.FindAllIndex([]byte(lines[i]), -1);

            line_md := lines[i];
            line_mu := line_md[0:m[0][0]]
            for j, v := range m {
                part := line_md[v[0]:v[1]]
                re_link := regexp.MustCompile(`\[[\w\s']+\]`);
                in_link := re_link.FindIndex([]byte(part))
                re_text := regexp.MustCompile(`\([\w:/.?=]+\)`);
                in_text := re_text.FindIndex([]byte(part))

                link := part[in_link[0]+1:in_link[1]-1]
                text := part[in_text[0]+1:in_text[1]-1]

                line_mu += "<a href=\"" + text + "\">" + link + "</a>";
                if j != len(m)-1 {
                    line_mu += line_md[m[j][1]:m[j+1][0]]
                }
            }
            line_mu += line_md[m[len(m)-1][1]:len(line_md)]
            lines[i] = line_mu;
        }
    }

    for i := 0; i < length; i++ {
        curr := getBlockSemantic(lines[i])
        if curr == NONE {
            first = i + 1;
            continue;
        }

        // There are elements which do not care about the next element.
        switch curr {
        case TAG_HEADER6:
            markup = append(markup, "<h6>" + lines[i][7:] + "</h6>")
            first = i + 1;
            continue;

        case TAG_HEADER5:
            markup = append(markup, "<h5>" + lines[i][6:] + "</h5>")
            first = i + 1;
            continue;

        case TAG_HEADER4:
            markup = append(markup, "<h4>" + lines[i][5:] + "</h4>")
            first = i + 1;
            continue;

        case TAG_HEADER3:
            markup = append(markup, "<h3>" + lines[i][4:] + "</h3>")
            first = i + 1;
            continue;

        case TAG_HEADER2:
            markup = append(markup, "<h2>" + lines[i][3:] + "</h2>")
            first = i + 1;
            continue;

        case TAG_HEADER1:
            markup = append(markup, "<h1 id=\"blogtitle\">" + lines[i][2:] + "</h1>")
            first = i + 1;
            continue;

        case TAG_META:
            if strings.HasPrefix(lines[i], "::title:") {
                metadata.Title = lines[i][8:]
            }
            if strings.HasPrefix(lines[i], "::timestamp:") {
                timestamp, _ := strconv.Atoi(lines[i][12:])
                metadata.Timestamp = uint64(timestamp)
            }
            first = i + 1;
            continue
        }


        next := NONE

        // If it's the last line, there is no next.
        if i < length-1 {
            next = getBlockSemantic(lines[i+1])
        }

        // Are we still in a block?
        if curr == next {
            continue
        }

        // There are elements which come in blocks
        switch curr {
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
                r := regexp.MustCompile(`^ [\d]+. `);
                indices := r.FindIndex([]byte(lines[j]));
                markup = append(markup, "<li>" + lines[j][indices[1]:] + "</li>")
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

    return markup, metadata
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

