package database

import (
	"log"
	"real-time-forum/backend/structs"
)

func GetChatHistory(senderID, receiverID int) ([]structs.Message, error) {
    query := `
        SELECT message_Id, content, created_at, sender_Id, receiver_Id
        FROM Message
        WHERE (sender_Id = ? AND receiver_Id = ?)
           OR (sender_Id = ? AND receiver_Id = ?)
        ORDER BY created_at ASC
    `
    rows, err := db.Query(query, senderID, receiverID, receiverID, senderID)
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
            log.Printf("Error scanning message row: %v", err)
            continue
        }
        messages = append(messages, msg)
    }

    return messages, nil
}

func GetAllChats(userID int) ([]structs.User, error) {
    query := `
        SELECT u.user_Id, u.firstName, u.lastName, u.username 
        FROM User u
        INNER JOIN (
            SELECT 
                CASE 
                    WHEN sender_Id = ? THEN receiver_Id 
                    ELSE sender_Id 
                END AS user_Id,
                MAX(created_at) AS last_message_time
            FROM Message
            WHERE sender_Id = ? OR receiver_Id = ?
            GROUP BY user_Id
        ) latest_messages
        ON u.user_Id = latest_messages.user_Id
        ORDER BY latest_messages.last_message_time DESC
    `

    rows, err := db.Query(query, userID, userID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []structs.User
    for rows.Next() {
        var user structs.User
        if err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}
