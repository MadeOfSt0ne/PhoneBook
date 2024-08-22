package postgresql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

func openConnectoin() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// возвращает ID пользователя username или -1 если пользователь не существует
func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnectoin()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" WHERE username = '%s'`, username)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}
	defer rows.Close()
	return userID
}

// Добавляет нового пользователя в бд и возвращает новый ID или -1 если произошла ошибка
func addUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnectoin()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", Username)
		return -1
	}

	insertStatement := `INSERT INTO "users" ("username") VALUES ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	insertStatement = `INSERT INTO "userdata" ("userid", "name", "surname", "description") VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return userID
}

// Удаляет существующего пользователя
func deleteUser(id int) error {
	db, err := openConnectoin()
	if err != nil {
		return err
	}
	defer db.Close()

	statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
	rows, err := db.Query(statement)
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	deleteStatement := `DELETE FROM "userdata" WHERE userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	deleteStatement = `DELETE FROM "users" WHERE id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	return nil
}
