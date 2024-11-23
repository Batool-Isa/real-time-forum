package database

import (
	"log"
	"real-time-forum/backend/structs"
)

func GetChatHistory(user1ID, user2ID int) ([]structs.Message, error) {
	query := `
		SELECT message_Id, content, created_at, sender_Id, receiver_Id 
		FROM Message 
		WHERE (sender_Id = ? AND receiver_Id = ?) 
		OR (sender_Id = ? AND receiver_Id = ?)
		ORDER BY created_at ASC`
	
	rows, err := db.Query(query, user1ID, user2ID, user2ID, user1ID)
	if err != nil {
		log.Printf("Error querying chat history: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []structs.Message
	for rows.Next() {
		var msg structs.Message
		err := rows.Scan(&msg.MessageID, &msg.Content, &msg.CreatedAt, &msg.SenderID, &msg.ReceiverID)
		if err != nil {
			log.Printf("Error scanning message: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
