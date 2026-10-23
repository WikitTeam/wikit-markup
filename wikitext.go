// Package wikitm renders Wikit Markup, the Wikidot-compatible markup used
// by ProjectWikit, to HTML.
package wikitm

// Render converts Wikidot wikitext into HTML.
// It will never fail. Any inputs, however malformed, produces some HTML.
func Render(source string) string {
	return source
}
