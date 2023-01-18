package utils

import (
	"bytes"

	"github.com/yuin/goldmark"
)

// MfmToHtml converts MFM to HTML.
func MfmToHtml(mfm string) string {
	// FIXME: 因为并没有MFM的Golang实现, 所以这里使用了一个Markdown的解析器.
	// 但是这样会导致一些问题, 比如Markdown的语法和MFM的语法不一致, 所以这里需要一个MFM的Golang实现.
	var buf bytes.Buffer
	_ = goldmark.Convert([]byte(mfm), &buf)
	s := buf.String()
	return s[:len(s)-1]
}
