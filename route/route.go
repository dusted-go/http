package route

import "strings"

// ShiftPath splits off the first component of path.
// Head will never contain a slash and tail will always be a rooted path without trailing slash.
func ShiftPath(path string) (head, tail string) {
	i := strings.Index(path[1:], "/") + 1
	if i <= 0 {
		return path[1:], "/"
	}
	return path[1:i], path[i:]
}
