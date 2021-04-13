package html_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratom/html"
)

func TestChildren(t *testing.T) {
	layout := html.Children{
		html.DIV(html.Id("foo")),
		html.DIV(html.Id("bar")),
	}

	result, err := html.RenderString(layout)
	require.NoError(t, err)

	content, err := ioutil.ReadFile("testdata/children.html")
	require.NoError(t, err)
	require.Equal(t, string(content), result)
}
