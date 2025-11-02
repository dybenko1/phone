package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "dyb_demp"
	password = "postgres"
)

func main() {

	db := connecting_db()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!\n")

	rawNumbers := get_rawNumbers(db)

	fmt.Println("\nCleaned numbers:")
	cleanedNumbers := cleanNumbers(rawNumbers)
	for _, number := range cleanedNumbers {
		fmt.Println(number)
	}

}

func connecting_db() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func get_rawNumbers(db *sql.DB) []string {
	rows, err := db.Query("SELECT first_name FROM raw_numbers LIMIT $1", 150)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()

	var rawNumbers []string
	for rows.Next() {
		var phone string
		err = rows.Scan(&phone)
		if err != nil {
			panic(err)
		}
		rawNumbers = append(rawNumbers, phone)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return rawNumbers
}

func cleanNumber(input string) string {
	// Regular expression to match non-digit characters
	re := regexp.MustCompile(`[^\d]`)
	// Replace all non-digit characters with empty string
	return re.ReplaceAllString(input, "")
}

func cleanNumbers(rawNumbers []string) []string {
	cleanNumbers := make([]string, len(rawNumbers))
	for i, rawNumber := range rawNumbers {

		cleanNumbers[i] = cleanNumber(rawNumber)
	}
	return cleanNumbers
}
