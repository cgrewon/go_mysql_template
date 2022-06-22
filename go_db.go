package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	id        int
	username  string
	password  string
	createdAt time.Time
}

type CreateUser struct {
	username  string
	password  string
	createdAt time.Time
}

func createTable(db *sql.DB) {
	{ // Create a new table

		tableQuery := `
			SELECT * FROM users;
		`
		_, err := db.Query(tableQuery)
		if err != nil {
			log.Fatal(err)

		} else {
			log.Println("users Table was created already.")
			return
		}

		log.Println("Tables")

		log.Println("Create Table ...")
		query := `
            CREATE TABLE users (
                id INT AUTO_INCREMENT,
                username TEXT NOT NULL,
                password TEXT NOT NULL,
                created_at DATETIME,
                PRIMARY KEY (id)
            );`

		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}

}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/go_sample?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func insertNewUser(db *sql.DB, newUser CreateUser) int64 {
	log.Println("Insert new user ...")
	// username := "johndoe"
	// password := "secret"
	// createdAt := time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, newUser.username, newUser.password, newUser.createdAt)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	fmt.Println(id)

	return id
}

func getUser(db *sql.DB, _id int64) {
	log.Println("Query single user ...")
	var (
		id        int
		username  string
		password  string
		createdAt time.Time
	)

	query := "SELECT id, username, password, created_at FROM users WHERE id = ?"
	if err := db.QueryRow(query, _id).Scan(&id, &username, &password, &createdAt); err != nil {
		log.Fatal(err)
	}

	fmt.Println(id, username, password, createdAt)

}

func getAllUsers(db *sql.DB) []user {
	log.Println("Query all users ...")

	rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user

		err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", users)

	return users
}

func deleteUser(db *sql.DB, id int64) sql.Result {
	log.Println("Delete user from table ...")
	res, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func main() {

	db := connectDB()

	createTable(db)

	user1 := CreateUser{
		username:  "user1",
		password:  "asd",
		createdAt: time.Now(),
	}

	user2 := CreateUser{
		username:  "user2",
		password:  "asd2",
		createdAt: time.Now(),
	}

	id1 := insertNewUser(db, user1)
	id2 := insertNewUser(db, user2)

	getUser(db, id1)

	getAllUsers(db)

	deleteUser(db, id2)
}
