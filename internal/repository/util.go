package repository

import "strings"

// escapeLike makes %, _ and \ in user input match literally.
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
