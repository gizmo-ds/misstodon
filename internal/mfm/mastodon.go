package mfm

import "golang.org/x/net/html"

func MastodonHashtagHandler(node *html.Node, m MfmNode, serverUrl string) {
	a := &html.Node{
		Type: html.ElementNode,
		Data: "a",
	}
	a.Attr = append(a.Attr,
		html.Attribute{
			Key: "href",
			Val: serverUrl + "/tags/" + m.Props["hashtag"].(string),
		},
		html.Attribute{Key: "class", Val: "mention hashtag"},
		html.Attribute{Key: "rel", Val: "nofollow noopener noreferrer"},
		html.Attribute{Key: "target", Val: "_blank"},
	)
	tag := &html.Node{
		Type: html.ElementNode,
		Data: "span",
	}
	tag.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: m.Props["hashtag"].(string),
	})
	a.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: "#",
	})
	a.AppendChild(tag)
	node.AppendChild(a)
}
