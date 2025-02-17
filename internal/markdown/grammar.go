package markdown

import (
	. "github.com/samuellando/gositter"
)

var G = CreateGrammar("root", map[string]Expression{
	"root": Repeat1(
		Choice(
			Ref("tag"),
			Ref("newline"))),
	"tag": Choice(
        Ref("header"),
		Ref("blockquote"),
		Ref("codeblock"),
		Ref("ul"),
		Ref("ol"),
		Ref("p")),

    "header": Choice(
		Ref("h6"),
		Ref("h5"),
		Ref("h4"),
		Ref("h3"),
		Ref("h2"),
		Ref("h1"),
    ),

	"h1": Seq(Terminal("# "), Ref("span")),
	"h2": Seq(Terminal("## "), Ref("span")),
	"h3": Seq(Terminal("### "), Ref("span")),
	"h4": Seq(Terminal("#### "), Ref("span")),
	"h5": Seq(Terminal("##### "), Ref("span")),
	"h6": Seq(Terminal("###### "), Ref("span")),

	"p": Ref("lines"),
    "lines": Choice(
			Seq(
				Ref("inner"),
				Ref("newline"),
				Ref("lines")),
			Seq(Ref("inner"))),
    "inner": Repeat1(Choice(
		Ref("a"),
		Ref("img"),
		Ref("textchar"))),
    "text": Repeat1(Ref("textchar")),
    "textchar": Choice(
		Ref("char"),
		Ref("whitespace")),
    "span": Seq(Ref("inner")),
	"char":       Regex("[^\\s`]"),
	"whitespace": Regex(` |\t`),
	"newline": Choice(
		Terminal("\n"),
		Terminal("\r\n"),
	),

	"a": Seq(
		Terminal("["),
		Choice(
			Ref("img"),
			Ref("alt")),
		Terminal("]"),
		Optional(Seq(
			Terminal("("),
			Ref("href"),
			Terminal(")")))),
	"img": Seq(
		Terminal("!["),
		Ref("alt"),
		Terminal("]"),
		Optional(Seq(
			Terminal("("),
			Ref("href"),
			Terminal(")"))),
		Optional(Ref("params"))),
	"alt":  Regex(`[^\]]*`),
	"href": Regex(`[^\)]*`),

	"params": Seq(
		Terminal("{"),
		Ref("paramset"),
		Terminal("}")),
	"paramset": Choice(
		Seq(
			Ref("param"),
			Terminal(","),
			Ref("paramset")),
		Ref("param")),
	"param": Regex("[^,}]*"),

	"blockquote": Seq(
		Terminal(">"),
		Ref("p")),
	"codeblock": Seq(
		Terminal("```"),
		Ref("newline"),
        Repeat1(Seq(
            Ref("text"),
            Ref("newline"))),
        Terminal("```"),
		Optional(Ref("newline"))),

	"ul": Repeat1(
		Seq(
			Terminal("- "),
			Ref("li"),
			Optional(Ref("newline")))),
	"ol": Repeat1(
		Seq(
			Regex(`\d*. `),
			Ref("li"),
			Optional(Ref("newline")))),
	"li": Ref("span"),
})
