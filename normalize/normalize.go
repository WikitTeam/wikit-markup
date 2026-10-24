// Package normalize rewrites HTML into a canoical form.
package normalize

import (
	"regexp"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// 这个map是块级元素，浏览器不渲染他们周围的空白。
// Wikidot的输出里块级标签附近全是这些换行和缩进，所以说得把他们去掉。
var blockTags = map[string]bool{
	"address": true, "blockquote": true, "body": true, "dd": true, "div": true,
	"dl": true, "dt": true, "h1": true, "h2": true, "h3": true, "h4": true,
	"h5": true, "h6": true, "hr": true, "li": true, "ol": true, "p": true,
	"pre": true, "table": true, "tbody": true, "td": true, "tfoot": true,
	"th": true, "thead": true, "tr": true, "ul": true, "caption": true,
	"colgroup": true, "col": true,
}

// Wikidot输出里有些是那种自闭合的，举例比方说<br />。
// 但是别的地方可能会输出<br>，于是给他们都归一成<br>。
var voidTags = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"source": true, "track": true, "wbr": true,
}

var spaces = regexp.MustCompile(`[ \t\r\n\f]+`)

// HTML rewrites HTML into a canonical form.
// This is a tokenizer.
func HTML(s string) string {
	var toks []html.Token
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.CommentToken || tt == html.DoctypeToken {
			continue
		}
		toks = append(toks, z.Token())
	}

	var b strings.Builder
	pre := 0 // <pre>的depth
	for i, t := range toks {
		switch t.Type {
		case html.TextToken:
			if pre > 0 {
				b.WriteString(escapeText(t.Data))
				continue
			}
			text := spaces.ReplaceAllString(t.Data, " ")
			if inRawText(toks, i) {
				b.WriteString(strings.TrimSpace(text))
				continue
			}
			if isBlock(toks, i-1) {
				text = strings.TrimLeft(text, " ")
			}
			if isBlock(toks, i+1) {
				text = strings.TrimRight(text, " ")
			}
			b.WriteString(escapeText(text))
		case html.StartTagToken, html.SelfClosingTagToken:
			if blockTags[t.Data] {
				newline(&b)
			}
			writeTag(&b, t)
			if t.Data == "pre" && t.Type == html.StartTagToken {
				pre++
			}
		case html.EndTagToken:
			if voidTags[t.Data] {
				continue
			}
			if t.Data == "pre" && pre > 0 {
				pre--
			}
			b.WriteString("</" + t.Data + ">")
			if blockTags[t.Data] {
				newline(&b)
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func escapeText(s string) string {
	return strings.ReplaceAll(html.EscapeString(s), "\u00a0", "&nbsp;")
}

func isBlock(toks []html.Token, i int) bool {
	if i < 0 || i >= len(toks) {
		return false
	}
	t := toks[i]
	return t.Type != html.TextToken && blockTags[t.Data]
}

func inRawText(toks []html.Token, i int) bool {
	if i == 0 {
		return false
	}
	prev := toks[i-1]
	return prev.Type == html.StartTagToken && (prev.Data == "script" || prev.Data == "style")
}

func newline(b *strings.Builder) {
	if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") {
		b.WriteByte('\n')
	}
}

// 浏览器不会管属性的先后，然鹅Wikidot输出的属性可能是乱的，所以得特别处理一下。
// 另外wikidot还有一些多余的属性值空白也一同处理
func writeTag(b *strings.Builder, t html.Token) {
	attrs := slices.Clone(t.Attr)
	slices.SortStableFunc(attrs, func(x, y html.Attribute) int {
		return strings.Compare(x.Key, y.Key)
	})
	b.WriteString("<" + t.Data)
	for _, a := range attrs {
		val := strings.TrimSpace(spaces.ReplaceAllString(a.Val, " "))
		b.WriteString(" " + a.Key + `="` + html.EscapeString(val) + `"`)
	}
	b.WriteString(">")
	if voidTags[t.Data] {
		return
	}
	if t.Type == html.SelfClosingTagToken {
		b.WriteString("</" + t.Data + ">")
	}
}
