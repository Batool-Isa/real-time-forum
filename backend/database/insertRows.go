package database
import (
	//"fmt"
	"log"
	//"net/http"
	"time"
	"real-time-forum/backend/utils"
	//"golang.org/x/crypto/bcrypt"
)
func CreateUser(username string, email string, age int, gender string, firstName string, lastName string, pass string) error {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	// preparer statment to create user
	prepStatment, err := db.Prepare("INSERT INTO User(username, email, age, gender, firstName, lastName, password) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer prepStatment.Close()
	// excute the statment
	_, err = prepStatment.Exec(username, email, age, gender, firstName, lastName, pass)
	if err != nil {
		log.Println("Error executing insert statement:", err)
		return err
	}
	log.Println("User successfully inserted with username:", username)
	return nil
}
func InsertNewSession(session string, userID int) error {
	// preparer statment to create user
	prepStatment, err := db.Prepare("INSERT INTO Session(session, timestamp, user_Id) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	// excute the statment
	_, err = prepStatment.Exec(session, time.Now().Add(12*time.Hour), userID)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func InsertPost(user_Id int, post_data string, categoryName []string) error {
	postErr := utils.ValidatePost(post_data, categoryName)
	if postErr != nil {
		return postErr
	}
	stmt, err := db.Prepare("INSERT INTO Post(user_Id, post_content) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(user_Id, post_data)
	if err != nil {
		log.Println(err)
		return err
	}
	postID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, categoryName := range categoryName {
		var categoryID int
		err := db.QueryRow("SELECT category_Id FROM Category WHERE category_name = ?", categoryName).Scan(&categoryID)
		if err != nil {
			log.Println(err)
			return err
		}
		InsertPostCategories(int(postID), categoryID)
	}
	return nil
}
func InsertComment(comment string, user_id int, postId int) error {
	CommentErr := utils.ValidateInput(map[string]string{"comment": comment})
	if CommentErr != nil {
		return CommentErr
	}
	stmt, err := db.Prepare("INSERT INTO comments(comment, user_id, post_id) values (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(comment, user_id, postId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func InsertPostCategories(post_id int, category_id int) error {
	postCatErr := utils.ValidateInput(map[string]string{"post_Id": string(post_id), "category_Id": string(category_id)})
	if postCatErr != nil {
		return postCatErr
	}
	stmt, err := db.Prepare("INSERT INTO Post_Category(post_Id, category_Id) values (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(post_id, category_id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func InsertLikes(post_id int, user_id int) error {
	likeserr := utils.ValidateInput(map[string]string{"post_id": string(post_id), "user_id": string(user_id)})
	if likeserr != nil {
		return likeserr
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ?", post_id, user_id).Scan(&count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count == 1 {
		// If the like exists, delete it (unlike)
		stmt, err := db.Prepare("DELETE FROM likes WHERE post_id = ? AND user_id = ?")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(post_id, user_id)
		if err != nil {
			return err
		}
	} else {
		// If the like does not exist, insert it (like)
		stmt, err := db.Prepare("INSERT INTO likes(post_id, user_id) values (?, ?)")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(post_id, user_id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
func InsertDislikes(post_id int, user_id int) error {
	dislikeserr := utils.ValidateInput(map[string]string{"post_id": string(post_id), "user_id": string(user_id)})
	if dislikeserr != nil {
		return dislikeserr
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ? AND user_id = ?", post_id, user_id).Scan(&count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count == 1 {
		// If the like exists, delete it (unlike)
		stmt, err := db.Prepare("DELETE FROM dislikes WHERE post_id = ? AND user_id = ?")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(post_id, user_id)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		// If the like does not exist, insert it (like)
		stmt, err := db.Prepare("INSERT INTO dislikes(post_id, user_id) values (?, ?)")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(post_id, user_id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
func InsertCommentLikes(comment_id int, user_id int) error {
	commentLikeErr := utils.ValidateInput(map[string]string{"comment_id": string(comment_id), "user_id": string(user_id)})
	if commentLikeErr != nil {
		return commentLikeErr
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM likeComment WHERE comment_id = ? AND user_id = ?", comment_id, user_id).Scan(&count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count == 1 {
		// If the like exists, delete it (unlike)
		stmt, err := db.Prepare("DELETE FROM likeComment WHERE comment_id = ? AND user_id = ?")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(comment_id, user_id)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		// If the like does not exist, insert it (like)
		stmt, err := db.Prepare("INSERT INTO likeComment(comment_id, user_id) values (?, ?)")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(comment_id, user_id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
func InsertCommentDislikes(comment_id int, user_id int) error {
	commentDislikeErr := utils.ValidateInput(map[string]string{"comment_id": string(comment_id), "user_id": string(user_id)})
	if commentDislikeErr != nil {
		return commentDislikeErr
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM dislikeComment WHERE comment_id = ? AND user_id = ?", comment_id, user_id).Scan(&count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count == 1 {
		// If the like exists, delete it (unlike)
		stmt, err := db.Prepare("DELETE FROM dislikeComment WHERE comment_id = ? AND user_id = ?")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(comment_id, user_id)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		// If the like does not exist, insert it (like)
		stmt, err := db.Prepare("INSERT INTO dislikeComment(comment_id, user_id) values (?, ?)")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = stmt.Exec(comment_id, user_id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
func InsertSession(session string, user_id int) error {
	sessionErr := utils.ValidateInput(map[string]string{"session": session, "user_id": string(user_id)})
	if sessionErr != nil {
		return sessionErr
	}
	stmt, err := db.Prepare("INSERT INTO sessions(session, user_id, timestamp) values (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(session, user_id, time.Now().Add(12*time.Hour))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func CheckCategoryTable() error {
	var count = 0
	err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count == 0 {
		AddDummyData()
	}
	return nil
}
func InsertCategories(category_name string) error {
	categoryArray, err := GetCategories()
	if err != nil {
		log.Println(err)
		return err
	}
	catErr := utils.ValidateCategory(category_name, categoryArray)
	if catErr != nil {
		return catErr
	}
	stmt, err := db.Prepare("INSERT INTO Category(category_name) values(?)")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(category_name)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func AddDummyData() {
	// InsertUser(db, "zhashim", "12345678", "zahra@gmail.com")
	// InsertUser(db, "zee", "12345678", "z@gmail.com")
	// InsertUser(db, "zozo", "12345678", "zozo@gmail.com")
	InsertCategories("Sports")
	InsertCategories("Technology")
	InsertCategories("Education")
	InsertCategories("Health")
	InsertCategories("Entertainment")
	InsertCategories("Travel")
	InsertCategories("Finance")
	InsertCategories("Culture")
}


func SaveMessage(content string, senderID, receiverID int) error {
    query := `INSERT INTO Message (content, sender_Id, receiver_Id) VALUES (?, ?, ?)`
    stmt, err := db.Prepare(query)
    if err != nil {
        log.Printf("Error preparing save message statement: %v", err)
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(content, senderID, receiverID)
    if err != nil {
        log.Printf("Error executing save message: %v", err)
        return err
    }

    return nil
}
