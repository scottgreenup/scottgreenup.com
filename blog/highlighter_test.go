package blog

import (
    "strings"
    "testing"
)

func TestHighlight(t *testing.T) {

    colors := make(map[int]string)
    colors[0] = "#000000"
    colors[1] = "#000001"
    colors[2] = "#000002"

    words := make(map[string]int)
    words["int"] = 0

    cases := []string{
        "int", "<span style=\"color: " + colors[0] + ";\">int</span>",
        "integer int other", "integer <span style=\"color: " + colors[0] + ";\">int</span> other",
    };

    for i := 0; i < len(cases); i++ {
        line := highlightCustom(cases[i], words, colors)
        i++
        if strings.Compare(line, cases[i]) != 0 {
            t.Errorf("Highlighter failed:\n" +
                "From   : %s\n" +
                "To     : %s\n" +
                "Expect : %s\n",
                cases[i-1],
                line,
                cases[i],
            )
        }
    }
}

