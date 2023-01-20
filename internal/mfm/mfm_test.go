package mfm_test

import (
	"os"
	"testing"

	"github.com/gizmo-ds/misstodon/internal/mfm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	m.Run()
}

func TestParse(t *testing.T) {
	_, err := mfm.Parse("Hello, world!")
	assert.NoError(t, err)
}

func TestToHtml(t *testing.T) {
	t.Run("Plain", func(t *testing.T) {
		s, err := mfm.ToHtml("Hello, world!")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><span>Hello, world!</span></p>", s)
	})
	t.Run("Quote", func(t *testing.T) {
		s, err := mfm.ToHtml("> abc")
		assert.NoError(t, err)
		assert.Equal(t,
			`<p><blockquote><span>abc</span></blockquote></p>`, s)
	})
	t.Run("InlineCode", func(t *testing.T) {
		s, err := mfm.ToHtml("`abc`")
		assert.NoError(t, err)
		assert.Equal(t,
			`<p><code>abc</code></p>`, s)
	})
	t.Run("Search", func(t *testing.T) {
		s, err := mfm.ToHtml("MFM 書き方 Search")
		assert.NoError(t, err)
		assert.Equal(t,
			`<p><a href="https://www.google.com/search?q=MFM 書き方">MFM 書き方 Search</a></p>`, s)
	})
	t.Run("Text", func(t *testing.T) {
		s, err := mfm.ToHtml("hello world")
		assert.NoError(t, err)
		assert.Equal(t, s,
			`<p><span>hello world</span></p>`)
	})
	//NOTE: 当前版本的mfm.js(0.23.3)不支持, 这里使用了与当前版本行为一致的测试用例
	t.Run("BlockCode", func(t *testing.T) {
		s, err := mfm.ToHtml("```js\nabc\n````")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><span>```js<br>abc<br>````</span></p>", s)
	})
	t.Run("MathBlock", func(t *testing.T) {
		s, err := mfm.ToHtml("\\[a = 1\\]")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><code>a = 1</code></p>", s)

		s, err = mfm.ToHtml("\\[\na = 2\n\\]")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><code>a = 2</code></p>", s)
	})
	t.Run("Center", func(t *testing.T) {
		s, err := mfm.ToHtml("<center>abc</center>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><div><span>abc</span></div></p>", s)

		s, err = mfm.ToHtml("<center>\nabc\ndef\n</center>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><div><span>abc<br>def</span></div></p>", s)
	})
	t.Run("Fn?", func(t *testing.T) {
		s, err := mfm.ToHtml("***big!***")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i><span>big!</span></i></p>", s)
	})
	t.Run("Bold", func(t *testing.T) {
		s, err := mfm.ToHtml("**bold**")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><b><span>bold</span></b></p>", s)

		s, err = mfm.ToHtml("__bold__")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><b><span>bold</span></b></p>", s)

		s, err = mfm.ToHtml("<b>bold</b>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><b><span>bold</span></b></p>", s)
	})
	t.Run("Small", func(t *testing.T) {
		s, err := mfm.ToHtml("<small>small</small>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><small><span>small</span></small></p>", s)
	})
	t.Run("Strike", func(t *testing.T) {
		s, err := mfm.ToHtml("~~strike~~")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><del><span>strike</span></del></p>", s)

		s, err = mfm.ToHtml("<s>strike</s>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><del><span>strike</span></del></p>", s)
	})
	t.Run("Italic", func(t *testing.T) {
		s, err := mfm.ToHtml("<i>italic</i>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i><span>italic</span></i></p>", s)

		s, err = mfm.ToHtml("*italic*")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i><span>italic</span></i></p>", s)

		s, err = mfm.ToHtml("_italic_")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i><span>italic</span></i></p>", s)
	})
	t.Run("EmojiCode", func(t *testing.T) {
		s, err := mfm.ToHtml(":thinking_ai:")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p>\u200B:thinking_ai:\u200B</p>", s)
	})
	t.Run("UnicodeEmoji", func(t *testing.T) {
		s, err := mfm.ToHtml("$[shake 🍮]")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i>🍮</i></p>", s)

		s, err = mfm.ToHtml("$[spin.alternate 🍮]")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i>🍮</i></p>", s)

		s, err = mfm.ToHtml("$[shake.speed=1s 🍮]")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i>🍮</i></p>", s)

		s, err = mfm.ToHtml("$[flip.h,v MisskeyでFediverseの世界が広がります]")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><i><span>MisskeyでFediverseの世界が広がります</span></i></p>", s)
	})
	t.Run("Hashtag", func(t *testing.T) {
		s, err := mfm.ToHtml("#hello")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://misskey.io/tags/hello\" rel=\"tag\">#hello</a></p>", s)
	})
	t.Run("MathInline", func(t *testing.T) {
		s, err := mfm.ToHtml("\\(y = 2x\\)")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><code>y = 2x</code></p>", s)
	})
	t.Run("Link", func(t *testing.T) {
		s, err := mfm.ToHtml("[Misskey.io](https://misskey.io/)")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://misskey.io/\"><span>Misskey.io</span></a></p>", s)

		s, err = mfm.ToHtml("?[Misskey.io](https://misskey.io/)")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://misskey.io/\"><span>Misskey.io</span></a></p>", s)
	})
	t.Run("Mention", func(t *testing.T) {
		s, err := mfm.ToHtml("@user@misskey.io")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://misskey.io/@user\" class=\"u-url mention\">@user@misskey.io</a></p>", s)

		s, err = mfm.ToHtml("@user")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://misskey.io/@user\" class=\"u-url mention\">@user</a></p>", s)

		s, err = mfm.ToHtml("@gizmo_ds@liuli.lol")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://liuli.lol/@gizmo_ds\" class=\"u-url mention\">@gizmo_ds@liuli.lol</a></p>", s)
	})
	t.Run("URL", func(t *testing.T) {
		s, err := mfm.ToHtml("https://misskey.io/@ai")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://misskey.io/@ai\">https://misskey.io/@ai</a></p>", s)

		s, err = mfm.ToHtml("http://hoge.jp/abc")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"http://hoge.jp/abc\">http://hoge.jp/abc</a></p>", s)

		s, err = mfm.ToHtml("<https://misskey.io/@ai>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"https://misskey.io/@ai\">https://misskey.io/@ai</a></p>", s)

		s, err = mfm.ToHtml("<http://藍.jp/abc>")
		assert.NoError(t, err)
		assert.Equal(t,
			"<p><a href=\"http://藍.jp/abc\">http://藍.jp/abc</a></p>", s)
	})
}
