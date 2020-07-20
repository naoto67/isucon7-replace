package main

import (
	"fmt"
	"io/ioutil"
)

const PUBLIC_DIR = "/home/isucon/isubata/webapp/public/icons/"

func writeAvatarIcon(filename string, data []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("%s%s", PUBLIC_DIR, filename), data, 0666)
}
