package normalize

import "testing"

func TestHTML(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"plain text", "a  b\n c", "a b c"},
		{"whitespace around blocks", "\n\n<p>a</p>\n<p> b </p>\n", "<p>a</p>\n<p>b</p>"},
		{"inline whitespace kept", "<p>a <b>b</b> c</p>", "<p>a <b>b</b> c</p>"},
		{"void tags", "a<br />\nb<hr/>", "a<br> b\n<hr>"},
		{"attribute order", `<a href="x" class="y">t</a>`, `<a class="y" href="x">t</a>`},
		{"attribute whitespace", `<div style=" color: red;  ">x</div>`, `<div style="color: red;">x</div>`},
		{"entities", "<p>&#39;&quot;&lt;&amp;&nbsp;</p>", "<p>&#39;&#34;&lt;&amp;&nbsp;</p>"},
		{"pre keeps whitespace", "<pre>a\n  b</pre>", "<pre>a\n  b</pre>"},
		{"comments dropped", "a<!-- x -->b", "ab"},
		{"upper case tags", "<P CLASS=a>x</P>", `<p class="a">x</p>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTML(tt.in); got != tt.want {
				t.Errorf("HTML(%q)\n got: %q\nwant: %q", tt.in, got, tt.want)
			}
		})
	}
}
