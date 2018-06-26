package util

import (
	"os"
)

// CheckErr handleErr
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

// CheckFileIsExist check wether the file is exist,
// @filename string path of the file,
// @return bool wether the file is exist.
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func CheckDirIsExist(dirName string) bool {
	var exist = true
	if st, err := os.Stat(dirName); os.IsNotExist(err) || !st.IsDir() {
		exist = false
	}
	return exist
}
