package route

import "strings"

// ShiftPath splits off the first component of p.
// Head will never contain a slash and tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
