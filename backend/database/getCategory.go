package database

import (
	"real-time-forum/backend/struct"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetCategories() ([]structs.Category, error) {
	query := `SELECT category_name FROM Category;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categoryArray []structs.Category
	for rows.Next() {
		var cat structs.Category
		err := rows.Scan(&cat.Category)
		if err != nil {
			return nil, err
		}
		categoryArray = append(categoryArray, cat)
	}
	if err := rows.Err(); err != nil {
        log.Println(err)
		return nil, err
	}
    return categoryArray, nil

}
