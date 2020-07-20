package main

func queryMessagesWithUser(chanID, lastID int64) ([]Message, error) {
	return scanMessagesWithUser(db.Query("SELECT * FROM message INNER JOIN user ON user.id = message.user_id WHERE message.id > ? AND message.channel_id = ? ORDER BY message.id DESC LIMIT 100",
		lastID, chanID))
}
