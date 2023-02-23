package mfm

import "golang.org/x/net/html"

type mfmNodeType string

const (
	nodeTypeBold         mfmNodeType = "bold"
	nodeTypeSmall        mfmNodeType = "small"
	nodeTypeStrike       mfmNodeType = "strike"
	nodeTypeItalic       mfmNodeType = "italic"
	nodeTypeFn           mfmNodeType = "fn"
	nodeTypeBlockCode    mfmNodeType = "blockCode"
	nodeTypeCenter       mfmNodeType = "center"
	nodeTypeEmojiCode    mfmNodeType = "emojiCode"
	nodeTypeUnicodeEmoji mfmNodeType = "unicodeEmoji"
	nodeTypeHashtag      mfmNodeType = "hashtag"
	nodeTypeInlineCode   mfmNodeType = "inlineCode"
	nodeTypeMathInline   mfmNodeType = "mathInline"
	nodeTypeMathBlock    mfmNodeType = "mathBlock"
	nodeTypeLink         mfmNodeType = "link"
	nodeTypeMention      mfmNodeType = "mention"
	nodeTypeQuote        mfmNodeType = "quote"
	nodeTypeText         mfmNodeType = "text"
	nodeTypeUrl          mfmNodeType = "url"
	nodeTypeSearch       mfmNodeType = "search"
	nodeTypePlain        mfmNodeType = "plain"
)

type (
	MfmNode struct {
		Type     mfmNodeType
		Props    map[string]any
		Children []MfmNode
	}
	Option struct {
		Url            string
		HashtagHandler func(*html.Node, MfmNode, string)
	}
)
