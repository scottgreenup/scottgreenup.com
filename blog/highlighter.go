package blog

import (
    "fmt"
    "regexp"
)

const (
    SEM_PRIM = iota
    SEM_CONTROL
    SEM_KEYWORD
)

func highlight(line string, words map[string]int) string {
    // TODO move colors into CSS, use IDs instead
    colors := make(map[int]string)
    colors[SEM_PRIM]    = "#81A2BE"
    colors[SEM_CONTROL] = "#B294BB"
    colors[SEM_KEYWORD] = "#CC6666"
    return highlightCustom(line, words, colors)
}

// TODO speed this function up, it's super slow and shit.
func highlightCustom(line string, words map[string]int, colors map[int]string) string {

    flipbit := true
    tmp := ""
    for i := 0; i < len(line); i++ {
        if line[i] == '"' {
            if flipbit {
                tmp += "<span style=\"color: #B5BD68;\">"
                tmp += string(line[i])
            } else {
                tmp += string(line[i])
                tmp += "</span>"
            }
            flipbit = !flipbit
            continue
        }

        tmp += string(line[i])
    }
    line = tmp

    for k, v := range words {
        pre := fmt.Sprintf("<span style=\"color: %s;\">", colors[v])
        post := "</span>"

        reg, _ := regexp.Compile("(^|[^\\w])" + k + "($|[^\\w])")

        ranges := reg.FindAllStringIndex(line, -1)

        // BUG finding the same string twice, it doesn't work.
        for _, r := range(ranges) {
            if r[0] == 0 && r[1] == len(line) {
                line = pre + k + post
            } else if r[0] == 0 {
                line = pre + k + post + line[r[1]-1:]
            } else if r[1] == len(line) {
                line = line[0:r[0]+1] + pre + k + post
            } else {
                line = line[0:r[0]+1] + pre + k + post + line[r[1]-1:]
            }
        }
    }

    return line;
}


func CHighlight (line string) string {
    words := make(map[string]int)

    words["int"]    = SEM_PRIM
    words["float"]  = SEM_PRIM
    words["double"] = SEM_PRIM
    words["char"]   = SEM_PRIM
    words["long"]   = SEM_PRIM
    words["if"]     = SEM_CONTROL

    return highlight(line, words)
}

func PythonHighlight (line string) string {
    words := make(map[string]int)

    words["import"] = SEM_PRIM
    words["from"]   = SEM_PRIM
    words["if"]     = SEM_PRIM
    words["for"]    = SEM_PRIM
    words["not"]    = SEM_PRIM
    words["in"]     = SEM_PRIM
    words["len"]    = SEM_PRIM
    words["print"]  = SEM_PRIM
    words["True"]   = SEM_PRIM
    words["False"]  = SEM_PRIM

    words["continue"] = SEM_CONTROL
    words["break"]    = SEM_CONTROL

    return highlight(line, words)
}
