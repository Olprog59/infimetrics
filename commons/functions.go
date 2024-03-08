package commons

import "strconv"

// convert string to uint
func StringToUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}
