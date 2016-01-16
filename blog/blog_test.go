package blog

import (
    "testing"
)

func TestGetBlockSemantic_Header6(t *testing.T) {
    var test_good = []string{
        "###### Heading Goes Here",
        "###### Heading",
        "###### ",
        "###### #####",
    };
    for _, str := range(test_good) {
        if x := getBlockSemantic(str); x != TAG_HEADER6 {
            t.Errorf("%s should be a HEADER6\n", str)
        }
    }

    var test_bad = []string{
        "####### Heading Goes Here",
        "##### Heading",
        "######notheader ",
        "Words ###### more words",
    };
    for _, str := range(test_bad) {
        if x := getBlockSemantic(str); x == TAG_HEADER6 {
            t.Errorf("%s should not be a HEADER6\n", str)
        }
    }
}

func TestGetBlockSemantic_Header5(t *testing.T) {
    var test_good = []string{
        "##### Heading Goes Here",
        "##### Heading",
        "##### ",
        "##### #####",
    };
    for _, str := range(test_good) {
        if x := getBlockSemantic(str); x != TAG_HEADER5 {
            t.Errorf("%s should be a HEADER6\n", str)
        }
    }

    var test_bad = []string{
        "###### Heading Goes Here",
        "#### Heading",
        "#####notheader ",
        "Words ##### more words",
    };
    for _, str := range(test_bad) {
        if x := getBlockSemantic(str); x == TAG_HEADER5 {
            t.Errorf("%s should not be a HEADER5\n", str)
        }
    }
}

func TestGetBlockSemantic_UNORDERED_LIST(t *testing.T) {
    var test_good = []string{
        " * list element",
        " - list element",
        " + list element",
        " * ",
    };
    for _, str := range(test_good) {
        if x := getBlockSemantic(str); x != TAG_UNORDERED_LIST {
            t.Errorf("%s should be an UNORDERED_LIST\n", str)
        }
    }

    var test_bad = []string{
        " *list element",
        " -list element",
        " +list element",
        "* list element",
        "- list element",
        "+ list element",
        "* ",
        " *",
    };
    for _, str := range(test_bad) {
        if x := getBlockSemantic(str); x == TAG_UNORDERED_LIST {
            t.Errorf("%s should not be a UNORDERED_LIST\n", str)
        }
    }
}

func TestGetBlockSemantic_ORDERED_LIST(t *testing.T) {
    var test_good = []string{
        " 1. list element",
        " 21. list element",
        " 13. list element",
        " 0. ",
    };
    for _, str := range(test_good) {
        if x := getBlockSemantic(str); x != TAG_ORDERED_LIST {
            t.Errorf("%s should be an ORDERED_LIST\n", str)
        }
    }

    var test_bad = []string{
        " 1.list element",
        "1. list element",
        "1. ",
        " 1.",
        " a. ",
    };
    for _, str := range(test_bad) {
        if x := getBlockSemantic(str); x == TAG_ORDERED_LIST {
            t.Errorf("%s should not be a ORDERED_LIST\n", str)
        }
    }
}

func TestParse(t *testing.T) {
    input := []string{
        "###### Header Six",
        "Whatever is going on shall be hunted.",
        "If nothing is going on, then nothing shall be hunted.",
        "",
        "# Features of Travelling",
        "There are many features of travelling, here are some listed:",
        "",
        " * Excitement",
        " * Photography",
        " * Socialising",
        " * New experience",
        "",
        "## More Stuff",
        " 1. I can't believe.",
        " 2. Do it and do it now.",
    }

    expected := []string{
        "<h6>Header Six</h6>",
        "<p>",
        "Whatever is going on shall be hunted.",
        "If nothing is going on, then nothing shall be hunted.",
        "</p>",
        "<h1>Features of Travelling</h1>",
        "<p>",
        "There are many features of travelling, here are some listed:",
        "</p>",
        "<ul>",
        "<li>Excitement</li>",
        "<li>Photography</li>",
        "<li>Socialising</li>",
        "<li>New experience</li>",
        "</ul>",
        "<h2>More Stuff</h2>",
        "<ol>",
        "<li>I can't believe.</li>",
        "<li>Do it and do it now.</li>",
        "</ol>",
    }

    result := ParseHTML(input)
    smaller := len(result)
    if smaller > len(expected) {
        smaller = len(expected)
    }

    for i := 0; i < smaller; i++ {
        if result[i] != expected[i] {
            t.Errorf(
                "[%d] String 1 (result) does not equal String 2 (expected):\n1. <%s>\n2. <%s>",
                i,
                result[i],
                expected[i],
            );
        }
    }
}

