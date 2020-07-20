package main

import (
	"fmt"
	"io/ioutil"
)

const PUBLIC_DIR = "/home/isucon/isubata/webapp/public/icons/"

func writeAvatarIcon(filename string, data []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("%s%s", PUBLIC_DIR, filename), data, 0666)
}

type image struct {
	Name string `db:"name"`
	Data []byte `db:"data"`
}

func initImages() error {
	var images []image
	err := db.Select(&images, "SELECT * FROM images")

	if err != nil {
		return err
	}

	for _, v := range images {
		if err = writeAvatarIcon(v.Name, v.Data); err != nil {
			return err
		}
	}
	return nil
}
