package database

import (
	"real-time-forum/backend/structs"
)
func GetChatHistory(senderID, receiverID int) ([]structs.Message, error) {
    query := `
        SELECT content, sender_Id, receiver_Id, created_at 
        FROM Message 
        WHERE (sender_Id = ? AND receiver_Id = ?) 
        OR (sender_Id = ? AND receiver_Id = ?)
        ORDER BY created_at ASC
    `
    
    rows, err := db.Query(query, senderID, receiverID, receiverID, senderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []structs.Message
    for rows.Next() {
        var msg structs.Message
        err := rows.Scan(&msg.Content, &msg.SenderID, &msg.ReceiverID, &msg.CreatedAt)
        if err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }
    
    return messages, nil
}
