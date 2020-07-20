package main

import (
	"database/sql"
	"fmt"
)

func scanMessagesWithUser(rows *sql.Rows, e error) ([]Message, error) {
	if e != nil {
		return nil, e
	}
	defer func() {
		rows.Close()
	}()

	msgs := []Message{}
	for rows.Next() {
		var m Message
		var u User
		if err := rows.Scan(&m.ID, &m.ChannelID, &m.UserID, &m.Content, &m.CreatedAt, &u.ID, &u.Name, &u.Salt, &u.Password, &u.DisplayName, &u.AvatarIcon, &u.CreatedAt); err != nil {
			fmt.Println("SCAN: ", err)
			return nil, err
		}
		m.User = &u
		msgs = append(msgs, m)
	}
	return msgs, nil
}
