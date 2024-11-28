package database

import (
	//"fmt"
	"fmt"
	"log"
	//"net/http"
	"real-time-forum/backend/utils"
	"time"
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
	fmt.Println("DEBUG: Validating post content and categories")
	postErr := utils.ValidatePost(post_data, categoryName)
	if postErr != nil {
		return postErr
	}

	fmt.Println("DEBUG: Preparing to insert post into database")
	stmt, err := db.Prepare("INSERT INTO Post(user_Id, post_content) VALUES (?, ?)")
	if err != nil {
		fmt.Println("ERROR: Preparing statement failed", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user_Id, post_data)
	if err != nil {
		fmt.Println("ERROR: Executing statement failed", err)
		return err
	}

	postID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("ERROR: Retrieving last insert ID failed", err)
		return err
	}

	fmt.Println("DEBUG: Post inserted successfully, ID:", postID)
	for _, category := range categoryName {
		fmt.Println("DEBUG: Processing category:", category)
		var categoryID int
		err := db.QueryRow("SELECT category_Id FROM Category WHERE category_name = ?", category).Scan(&categoryID)
		if err != nil {
			fmt.Println("ERROR: Retrieving category ID failed", err)
			return err
		}

		err = InsertPostCategories(int(postID), categoryID)
		if err != nil {
			fmt.Println("ERROR: Inserting post-category mapping failed", err)
			return err
		}
	}
	return nil
}

func InsertComment(comment string, user_id int, postId int) error {
	CommentErr := utils.ValidateInput(map[string]string{"comment": comment})
	if CommentErr != nil {
		return CommentErr
	}
	stmt, err := db.Prepare("INSERT INTO Comment(comment_content, user_Id, post_Id) values (?, ?, ?)")
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
	fmt.Println("DEBUG: Inserting post-category mapping")
	stmt, err := db.Prepare("INSERT INTO Post_Category(post_Id, category_Id) values (?, ?)")
	if err != nil {
		fmt.Println("ERROR: Preparing post-category statement failed", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post_id, category_id)
	if err != nil {
		fmt.Println("ERROR: Executing post-category insert failed", err)
		return err
	}
	fmt.Println("DEBUG: Post-category mapping inserted successfully")
	return nil
}


func InsertLikes(post_id int, user_id int) error {
    log.Printf("Inserting like for PostID: %d, UserID: %d", post_id, user_id)

    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM Post_Like WHERE post_Id = ? AND user_Id = ?", post_id, user_id).Scan(&count)
    if err != nil {
        log.Println("Error querying Post_Like:", err)
        return err
    }

    if count == 1 {
        log.Println("Like already exists, removing it")
        stmt, err := db.Prepare("DELETE FROM Post_Like WHERE post_Id = ? AND user_Id = ?")
        if err != nil {
            log.Println("Error preparing DELETE statement:", err)
            return err
        }
        _, err = stmt.Exec(post_id, user_id)
        if err != nil {
            log.Println("Error executing DELETE statement:", err)
            return err
        }
    } else {
        log.Println("Like does not exist, adding it")
        stmt, err := db.Prepare("INSERT INTO Post_Like(post_Id, user_Id) VALUES (?, ?)")
        if err != nil {
            log.Println("Error preparing INSERT statement:", err)
            return err
        }
        _, err = stmt.Exec(post_id, user_id)
        if err != nil {
            log.Println("Error executing INSERT statement:", err)
            return err
        }
    }

    return nil
}

func InsertDislikes(post_id int, user_id int) error {
	dislikeserr := utils.ValidateInput(map[string]string{"post_Id": string(post_id), "user_Id": string(user_id)})
	if dislikeserr != nil {
		return dislikeserr
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Post_Dislike WHERE post_Id = ? AND user_Id = ?", post_id, user_id).Scan(&count)
	if err != nil {
		log.Println(err)
		return err
	}
	if count == 1 {
		// If the like exists, delete it (unlike)
		stmt, err := db.Prepare("DELETE FROM Post_Dislike WHERE post_Id = ? AND user_Id = ?")
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
		stmt, err := db.Prepare("INSERT INTO Post_Dislike(post_Id, user_Id) values (?, ?)")
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
