package util

import "io/ioutil"

// WriteAllText 将给定text写入给定path
func WriteAllText(path, text string) error {
	return ioutil.WriteFile(path, []byte(text), 0o644)
}
