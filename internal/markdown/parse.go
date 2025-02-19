package markdown

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/samuellando/gositter"
)

var COMPOENTS, _ = template.New("components").ParseGlob("./templates/components/*")

func ToHtml(md string) (template.HTML, error) {
	tree, err := G.Parse(md)
	if err != nil {
		return "", err
	}
	return parseTree(tree)
}

type a struct {
	Href  string
	Inner template.HTML
}

type img struct {
	Src    string
	Alt    string
	Height string
	Width  string
}

type list struct {
	Lis []template.HTML
}

func parseTags(out io.Writer, t gositter.SyntaxTree) error {
	nodes := t.Nodes()
	// If this is a leaf, just return it's value
	if len(nodes) == 0 {
		out.Write([]byte(t.Value()))
		return nil
	}
	// Otherwise check if there is an exiting template
	var tag string
	var data interface{}
	var err error
	switch t.Tag() {
	case "header":
		sub := t.Find("span")[0]
		tag = t.Nodes()[0].Tag()
		data, err = parseTree(sub)
	case "blockquote":
		sub := t.Find("p")[0]
		tag = t.Tag()
		data, err = parseTree(sub)
	case "codeblock":
		lines := t.Find("text")
		inner := new(strings.Builder)
		for _, line := range lines {
			inner.Write([]byte(fmt.Sprintf("%s\n", line.Value())))
		}
		tag = "codeblock"
		data = inner.String()
	case "p", "span":
		tag = t.Tag()
		sub := t.Nodes()[0]
		data, err = parseTree(sub)
	case "a":
		var inner template.HTML
		ht := t.Find("href")[0]
		imgs := t.Find("img")
		if len(imgs) == 0 {
			at := t.Find("alt")[0]
			inner = template.HTML(at.Value())
		} else {
			inner, err = parseTree(imgs[0])
		}
		tag = "a"
		data = a{Href: ht.Value(), Inner: inner}
	case "img":
		img := new(img)
		img.Alt = t.Find("alt")[0].Value()
		img.Src = t.Find("href")[0].Value()
		params := t.Find("param")
		if len(params) >= 1 {
			img.Height = params[0].Value()
		}
		if len(params) >= 2 {
			img.Width = params[1].Value()
		}
		tag = "img"
		data = img
	case "ul", "ol":
		lis := t.Find("li")
		list := new(list)
		var inner template.HTML
		for _, li := range lis {
			inner, err = parseTree(li)
			if err != nil {
				break
			}
			list.Lis = append(list.Lis, inner)
		}
		tag = t.Tag()
		data = list
	default:
		// If there is no template, continue traversing the tree.
		for _, node := range nodes {
			err := parseTags(out, node)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// Execute the template
	if err != nil {
		return err
	}
	return COMPOENTS.ExecuteTemplate(out, tag, data)
}

func parseTree(t gositter.SyntaxTree) (template.HTML, error) {
	s := new(strings.Builder)
	err := parseTags(s, t)
	return template.HTML(s.String()), err
}
