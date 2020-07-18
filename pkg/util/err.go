package util

import "log"

// CheckErr is a helper for handling errors
func CheckErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}
