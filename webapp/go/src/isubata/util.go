package main

import (
	"fmt"
	"io/ioutil"
)

const PUBLIC_DIR = "/home/isucon/isubata/webapp/public/icons/"

func writeAvatarIcon(filename string, data []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("%s%s", PUBLIC_DIR, filename), data, 0666)
}

type Image struct {
	ID   int64  `json:"-" db:"id"`
	Name string `db:"name"`
	Data []byte `db:"data"`
}

func initImages() error {
	var images []Image
	err := db.Select(&images, "SELECT * FROM image")
	fmt.Println("DEBUG: INIT IMAGES")

	if err != nil {
		fmt.Println("ERROR: ", err)
		return err
	}

	for _, v := range images {
		if err = writeAvatarIcon(v.Name, v.Data); err != nil {
			fmt.Println("ERROR: ", err)
			return err
		}
	}
	fmt.Println("DEBUG: INIT IMAGES DONE")
	return nil
}

func jsonifyMessage(m Message) (map[string]interface{}, error) {
	r := make(map[string]interface{})
	r["id"] = m.ID
	r["user"] = *m.User
	r["date"] = m.CreatedAt.Format("2006/01/02 15:04:05")
	r["content"] = m.Content
	return r, nil
}
