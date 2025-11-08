package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Define constants for file extensions.
const (
	extHTML = ".html"
	extCSS  = ".css"
	extJS   = ".js"
	extTS   = ".ts"
	extVue  = ".vue"
	extGo   = ".go"
)

func removeHTMLComments(content []byte) ([]byte, error) {
	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return nil, err
	}

	var removeComments func(*html.Node)
	removeComments = func(n *html.Node) {
		for c := n.FirstChild; c != nil; {
			next := c.NextSibling
			if c.Type == html.CommentNode {
				n.RemoveChild(c)
			}
			removeComments(c)
			c = next
		}
	}

	removeComments(doc)

	var b bytes.Buffer
	if err := html.Render(&b, doc); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func removeCSSComments(content []byte) []byte {
	re := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	return re.ReplaceAll(content, []byte{})
}

func removeJSComments(content []byte) []byte {
	re := regexp.MustCompile(`(//.*)|(/\*[\s\S]*?\*/)`)
	return re.ReplaceAll(content, []byte{})
}

func removeVueComments(content []byte) []byte {
	templateRe := regexp.MustCompile(`(?s)<template>(.*)</template>`)
	scriptRe := regexp.MustCompile(`(?s)<script.*?>(.*)</script>`)
	styleRe := regexp.MustCompile(`(?s)<style.*?>(.*)</style>`)

	// Extract content
	templateMatch := templateRe.FindSubmatch(content)
	scriptMatch := scriptRe.FindSubmatch(content)
	styleMatch := styleRe.FindSubmatch(content)

	var newTemplate, newScript, newStyle []byte
	var err error

	if len(templateMatch) > 1 {
		newTemplate, err = removeHTMLComments(templateMatch[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error removing HTML comments from Vue template: %v\n", err)
			newTemplate = templateMatch[1]
		}
	}

	if len(scriptMatch) > 1 {
		newScript = removeJSComments(scriptMatch[1])
	}

	if len(styleMatch) > 1 {
		newStyle = removeCSSComments(styleMatch[1])
	}

	// Reconstruct the file
	var result []byte
	if len(templateMatch) > 0 {
		result = append(result, []byte("<template>")...)
		result = append(result, newTemplate...)
		result = append(result, []byte("</template>\n")...)
	}
	if len(scriptMatch) > 0 {
		scriptTag := scriptRe.Find(content)
		scriptTagString := string(scriptTag)
		endOfOpenTag := strings.Index(scriptTagString, ">")
		result = append(result, []byte(scriptTagString[:endOfOpenTag+1])...)
		result = append(result, newScript...)
		result = append(result, []byte("</script>\n")...)
	}
	if len(styleMatch) > 0 {
		styleTag := styleRe.Find(content)
		styleTagString := string(styleTag)
		endOfOpenTag := strings.Index(styleTagString, ">")
		result = append(result, []byte(styleTagString[:endOfOpenTag+1])...)
		result = append(result, newStyle...)
		result = append(result, []byte("</style>\n")...)
	}

	if len(result) == 0 {
		return content
	}

	return result
}
