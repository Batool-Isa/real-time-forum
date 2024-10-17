package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func CreateDB(databaseName string) error {
	var err error
	// open the databse , create one if the db doesn't exist
	db, err = sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}



func ExecuteSQLQuery(db *sql.DB, query string) error {
	var err error
	_, err = db.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateTables() {
	categoryTable := `
	CREATE TABLE IF NOT EXISTS Category (
		category_Id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_name VARCHAR NOT NULL
	);`

	commentTable := `
	CREATE TABLE IF NOT EXISTS Comment (
		comment_Id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_content VARCHAR NOT NULL
	);`

	postTable := `
	CREATE TABLE IF NOT EXISTS Post (
		post_Id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_content VARCHAR NOT NULL,
		likesNum VARCHAR DEFAULT 0,
		dislikesNum VARCHAR DEFAULT 0
	);`

	postCategoryTable := `
	CREATE TABLE IF NOT EXISTS Post_Category (
		post_Id INTEGER,
		category_Id INTEGER,
		PRIMARY KEY (post_Id, category_Id),
		FOREIGN KEY (post_Id) REFERENCES Post(post_Id),
		FOREIGN KEY (category_Id) REFERENCES Category(category_Id)
	);`

	userTable := `
	CREATE TABLE IF NOT EXISTS User (
		user_Id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR NOT NULL UNIQUE,
		email VARCHAR NOT NULL,
		age INTEGER NOT NULL,
		gender VARCHAR NOT NULL,
		firstName VARCHAR NOT NULL,
		lastName VARCHAR NOT NULL,
		password VARCHAR NOT NULL
	);`

	messageTable := `
	CREATE TABLE IF NOT EXISTS Message (
		message_Id INTEGER PRIMARY KEY AUTOINCREMENT,
		content VARCHAR NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		sender_Id INTEGER,
		receiver_Id INTEGER,
		FOREIGN KEY (sender_Id) REFERENCES User(user_Id),
		FOREIGN KEY (receiver_Id) REFERENCES User(user_Id)
	);`

	commentLikesTable := `
	CREATE TABLE IF NOT EXISTS Comment_Like (
		comment_Id INTEGER,
		user_Id INTEGER,
		PRIMARY KEY (comment_Id, user_Id),
		FOREIGN KEY (comment_Id) REFERENCES Comment(comment_Id),
		FOREIGN KEY (user_Id) REFERENCES User(user_Id)
	);`

	commentDislikesTable := `
	CREATE TABLE IF NOT EXISTS Comment_Dislike (
		comment_Id INTEGER,
		user_Id INTEGER,
		PRIMARY KEY (comment_Id, user_Id),
		FOREIGN KEY (comment_Id) REFERENCES Comment(comment_Id),
		FOREIGN KEY (user_Id) REFERENCES User(user_Id)
	);`

	sessionTable := `
	CREATE TABLE IF NOT EXISTS Session (
		session_Id INTEGER PRIMARY KEY AUTOINCREMENT,
		session VARCHAR NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		user_Id INTEGER,
		FOREIGN KEY (user_Id) REFERENCES User(user_Id)
	);`

	postDislikeTable := `
	CREATE TABLE IF NOT EXISTS Post_Dislike (
		post_Id INTEGER,
		user_Id INTEGER,
		PRIMARY KEY (post_Id, user_Id),
		FOREIGN KEY (post_Id) REFERENCES Post(post_Id),
		FOREIGN KEY (user_Id) REFERENCES User(user_Id)
	);`

	postLikeTable := `
	CREATE TABLE IF NOT EXISTS Post_Like (
		post_Id INTEGER,
		user_Id INTEGER,
		PRIMARY KEY (post_Id, user_Id),
		FOREIGN KEY (post_Id) REFERENCES Post(post_Id),
		FOREIGN KEY (user_Id) REFERENCES User(user_Id)
	);`

	// Execute the table creation queries
	ExecuteSQLQuery(db, userTable)
	ExecuteSQLQuery(db, categoryTable)
	ExecuteSQLQuery(db, postTable)
	ExecuteSQLQuery(db, commentTable)
	ExecuteSQLQuery(db, commentLikesTable)
	ExecuteSQLQuery(db, commentDislikesTable)
	ExecuteSQLQuery(db, postLikeTable)
	ExecuteSQLQuery(db, postDislikeTable)
	ExecuteSQLQuery(db, sessionTable)
	ExecuteSQLQuery(db, postCategoryTable)
	ExecuteSQLQuery(db, messageTable)
}

