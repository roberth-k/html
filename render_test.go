package html_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratom/html"
)

func TestRender(t *testing.T) {
	layout := html.HTML(
		html.HEAD(
			html.TITLE(html.Textf("greetings, programs!"))),
		html.BODY(
			html.DIV(
				html.Id("test"),
				html.Textf("Test!"),
				html.B(html.Textf("test!")))))

	result, err := html.RenderString(layout)
	require.NoError(t, err)

	content, err := ioutil.ReadFile("testdata/simple.html")
	require.NoError(t, err)
	require.Equal(t, string(content), result)
}
