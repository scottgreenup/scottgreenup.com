
function wrap(x, color) {
    return "<span style='color: " + color + ";'>" + x + "</span>";
}

function stringStartsWith(str, prefix) {
    return str.slice(0, prefix.length) == prefix;
}

function stringEndsWith(str, suffix) {
    return str.slice(-suffix.length) == suffix;
}

var code_elements = document.getElementsByTagName("code");
console.log("code elements:")
console.log(code_elements)

var keywords = ["import", "from", "if", "for", "not", "in", "len",
"print", "True", "False"];
var keywords_alt = ["continue", "break"]

for (n = 0; n < code_elements.length; n++) {
    var inner_html = code_elements[n].innerHTML;
    var final_html = "";

    var splat = inner_html.split('\n')
    for (i = 0; i < splat.length; i++) {
        words = splat[i].split(' ');
        if (i != 0) {
            final_html += '\n';
        }


        for (j = 0; j < words.length; j++) {
            var w = words[j];
            for (k = 0; k < keywords.length; k++) {
                if (w == keywords[k]) {
                    words[j] = wrap(words[j], "blue");
                    break;
                }
            }
            for (k = 0; k < keywords_alt.length; k++) {
                if (w == keywords_alt[k]) {
                    words[j] = wrap(words[j], "red");
                    break;
                }
            }

            if (j != 0) {
                final_html += ' ' + words[j];
            } else {
                final_html += words[j];
            }
        }

    }
    code_elements[n].innerHTML = final_html;
}
