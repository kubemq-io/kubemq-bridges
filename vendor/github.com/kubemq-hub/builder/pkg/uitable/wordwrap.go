//nolint
package uitable

import (
	"bytes"
	"unicode"
)

// WrapString wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func WrapString(s string, lim uint) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current uint
	var wordBuf, spaceBuf bytes.Buffer
	var wordWidth, spaceWidth int

	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+uint(spaceWidth) > lim {

				} else {
					current += uint(spaceWidth)
					_, _ = spaceBuf.WriteTo(buf)
					spaceWidth += StringWidth(buf.String())
				}
				spaceBuf.Reset()
				spaceWidth = 0
			} else {
				current += uint(spaceWidth + wordWidth)
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
				spaceWidth = 0
				wordWidth = 0
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += uint(spaceWidth + wordWidth)
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
				spaceWidth = 0
				wordWidth = 0
			}

			spaceBuf.WriteRune(char)
			spaceWidth += RuneWidth(char)
		} else {
			wordBuf.WriteRune(char)
			wordWidth += RuneWidth(char)

			if current+uint(spaceWidth+wordWidth) > lim && uint(wordWidth) < lim {
				buf.WriteRune('\n')
				current = 0
				spaceBuf.Reset()
				spaceWidth = 0
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+uint(spaceWidth) <= lim {
			_, _ = spaceBuf.WriteTo(buf)
		}
	} else {
		_, _ = spaceBuf.WriteTo(buf)
		_, _ = wordBuf.WriteTo(buf)
	}

	return buf.String()
}
