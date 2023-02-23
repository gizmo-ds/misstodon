//go:generate pnpm build
package mfm

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/dop251/goja"
	"github.com/gizmo-ds/misstodon/internal/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

//go:embed out.js
var script string
var DefaultMfmOption = Option{
	Url: "https://misskey.io",
}

var vm = goja.New()
var parseText func(string) string

func init() {
	_, err := vm.RunString(script)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run mfm.js")
	}
	if err = vm.ExportTo(vm.Get("parse"), &parseText); err != nil {
		log.Fatal().Err(err).Msg("Failed to export parse function")
	}
}

// Parse parses MFM to nodes.
func Parse(text string) ([]MfmNode, error) {
	var nodes []MfmNode
	err := json.Unmarshal([]byte(parseText(text)), &nodes)
	return nodes, err
}

// ToHtml converts MFM to HTML.
func ToHtml(text string, option ...Option) (string, error) {
	nodes, err := Parse(text)
	if err != nil {
		return "", err
	}
	return toHtml(nodes, option...)
}

func toHtml(nodes []MfmNode, option ...Option) (string, error) {
	node := &html.Node{
		Type: html.ElementNode,
		Data: "p",
	}
	if len(option) == 0 {
		option = append(option, DefaultMfmOption)
	}

	appendChildren(node, nodes, option...)

	var buf bytes.Buffer
	if err := html.Render(&buf, node); err != nil {
		return "", err
	}
	h := buf.String()
	// NOTE: misskey的br标签不符合XHTML 1.1，需要替换为<br>
	h = strings.ReplaceAll(h, "<br/>", "<br>")
	return h, nil
}

func appendChildren(parent *html.Node, children []MfmNode, option ...Option) {
	for _, child := range children {
		switch child.Type {
		case nodeTypePlain:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "span",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeText:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "span",
			}
			var text, ok = "", false
			if text, ok = child.Props["text"].(string); !ok {
				return
			}
			text = strings.ReplaceAll(text, "\r", "")

			arr := strings.Split(text, "\n")
			for i := 0; i < len(arr); i++ {
				if i > 0 && i < len(arr) {
					n.AppendChild(&html.Node{
						Type: html.ElementNode,
						Data: "br",
					})
				}
				n.AppendChild(&html.Node{
					Type: html.TextNode,
					Data: arr[i],
				})
			}

			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeBold:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "b",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeQuote:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "blockquote",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeInlineCode:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "code",
			}
			n.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: child.Props["code"].(string),
			})
			parent.AppendChild(n)

		case nodeTypeSearch:
			a := &html.Node{
				Type: html.ElementNode,
				Data: "a",
			}
			a.Attr = append(a.Attr, html.Attribute{
				Key: "href",
				Val: "https://www.google.com/search?q=" + child.Props["query"].(string),
			})
			a.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: child.Props["content"].(string),
			})
			parent.AppendChild(a)

		case nodeTypeMathBlock:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "code",
			}
			n.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: child.Props["formula"].(string),
			})
			parent.AppendChild(n)

		case nodeTypeCenter:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "div",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeFn:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "i",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeSmall:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "small",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeStrike:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "del",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeItalic:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "i",
			}
			appendChildren(n, child.Children, option...)
			parent.AppendChild(n)

		case nodeTypeBlockCode: // NOTE: 当前版本的mfm.js(0.23.3)不支持, 所以下面的代码没有进行测试
			pre := &html.Node{
				Type: html.ElementNode,
				Data: "pre",
			}
			inner := &html.Node{
				Type: html.ElementNode,
				Data: "code",
			}
			inner.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: child.Props["code"].(string),
			})
			pre.AppendChild(inner)
			parent.AppendChild(pre)

		case nodeTypeEmojiCode:
			parent.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: "\u200B:" + child.Props["name"].(string) + ":\u200B",
			})

		case nodeTypeUnicodeEmoji:
			parent.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: child.Props["emoji"].(string),
			})

		case nodeTypeHashtag:
			if option[0].HashtagHandler != nil {
				option[0].HashtagHandler(parent, child, option[0].Url)
				break
			}
			a := &html.Node{
				Type: html.ElementNode,
				Data: "a",
			}
			hashtag := child.Props["hashtag"].(string)
			a.Attr = append(a.Attr, html.Attribute{
				Key: "href",
				Val: option[0].Url + "/tags/" + hashtag,
			})
			a.Attr = append(a.Attr, html.Attribute{
				Key: "rel",
				Val: "tag",
			})
			a.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: "#" + hashtag,
			})
			parent.AppendChild(a)

		case nodeTypeMathInline:
			n := &html.Node{
				Type: html.ElementNode,
				Data: "code",
			}
			n.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: child.Props["formula"].(string),
			})
			parent.AppendChild(n)

		case nodeTypeLink:
			a := &html.Node{
				Type: html.ElementNode,
				Data: "a",
			}
			a.Attr = append(a.Attr, html.Attribute{
				Key: "href",
				Val: child.Props["url"].(string),
			})
			appendChildren(a, child.Children, option...)
			parent.AppendChild(a)

		case nodeTypeMention:
			a := &html.Node{
				Type: html.ElementNode,
				Data: "a",
			}
			acct := child.Props["acct"].(string)
			username, host := utils.AcctInfo(acct)
			if host == "" {
				host = option[0].Url[8:]
			}
			a.Attr = append(a.Attr,
				html.Attribute{
					Key: "href",
					Val: "https://" + host + "/@" + username,
				},
				html.Attribute{
					Key: "class",
					Val: "u-url mention",
				})
			a.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: acct,
			})
			parent.AppendChild(a)

		case nodeTypeUrl:
			a := &html.Node{
				Type: html.ElementNode,
				Data: "a",
			}
			a.Attr = append(a.Attr, html.Attribute{
				Key: "href",
				Val: child.Props["url"].(string),
			})
			a.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: child.Props["url"].(string),
			})
			parent.AppendChild(a)

		default:
			log.Warn().Str("type", string(child.Type)).Msg("unknown node type")
		}
	}
}
