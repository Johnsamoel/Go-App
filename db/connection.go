package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func OpenDbConnection() (*sql.DB, error) {
	// Open a connection to the MySQL database
	database, err := sql.Open("mysql", "root:5256@tcp(localhost:3306)/gowallet?parseTime=true")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, err
	}

	// Check if the connection is successful
	err = database.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return nil, err
	}

	fmt.Println("Connected to the database!")

	CreateDbTables(database)

	return database, nil
}

func CreateDbTables(Db *sql.DB) error {
	// Create channels to communicate completion status
	usersCreationDone := make(chan bool)
	walletsCreationDone := make(chan bool)
	otpsCreationDone := make(chan bool)
	transactionDone := make(chan bool)

	// Execute createUsersTable and createWalletsTable concurrently
	go func() {
		usersCreationDone <- createUsersTable(Db)
	}()

	go func() {
		walletsCreationDone <- createWalletsTable(Db)
	}()

	go func() {
		otpsCreationDone <- createOtpsTable(Db)
	}()

	go func() {
		transactionDone <- createTransactionsTable(Db)
	}()

	// Wait for both tasks to complete
	var usersErr, walletsErr, otpsErr, transErr  bool
	for i := 0; i < 2; i++ {
		select {
		case err := <-usersCreationDone:
			usersErr = err
		case err := <-walletsCreationDone:
			walletsErr = err
		case err := <-otpsCreationDone:
			otpsErr = err
		case err := <-otpsCreationDone:
			transErr = err
		}
	}

	// Check for errors
	if usersErr != true {
		return fmt.Errorf("error creating users table: %v", usersErr)
	}
	if walletsErr != true {
		return fmt.Errorf("error creating wallets table: %v", walletsErr)
	}
	if otpsErr != true {
		return fmt.Errorf("error creating otps table: %v", otpsErr)
	}
	if transErr != true {
		return fmt.Errorf("error creating transactions table: %v", transErr)
	}

	return nil
}

// generate users table if doesn't exist.
func createUsersTable(Db *sql.DB) bool {
	// Define the SQL query to create the table
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL,
			phoneNumber VARCHAR(50) NOT NULL UNIQUE
        )
    `

	// Execute the SQL query to create the table
	_, err := Db.Exec(query)
	if err != nil {
		return false
	}

	return true
}

// generate wallets table if doesn't exist.
func createWalletsTable(Db *sql.DB) bool {
	// Define the SQL query to create the table
	query := `
		CREATE TABLE IF NOT EXISTS wallets (
			id INT AUTO_INCREMENT PRIMARY KEY,
			balance DOUBLE DEFAULT 0,
			userId INT NOT NULL,
			status VARCHAR(50) DEFAULT 'ACTIVE',
			FOREIGN KEY (userId) REFERENCES users(id)
		)
	`

	// Execute the SQL query to create the table
	_, err := Db.Exec(query)
	if err != nil {
		return false
	}

	return true
}

// generate otps table if doesn't exist.
func createOtpsTable(Db *sql.DB) bool {
	// Define the SQL query to create the table
	query := `
		CREATE TABLE IF NOT EXISTS otps (
			id INT AUTO_INCREMENT PRIMARY KEY,
			userId INT NOT NULL,
			walletId INT NOT NULL,
			status VARCHAR(50),
			otp VARCHAR(50) NOT NULL,
			createdAt DATETIME NOT NULL,
			FOREIGN KEY (userId) REFERENCES users(id),
			FOREIGN KEY (walletId) REFERENCES wallets(id)
		)	
	`

	// Execute the SQL query to create the table
	_, err := Db.Exec(query)
	if err != nil {
		return false
	}

	return true
}

// generate transactions table if doesn't exist.
func createTransactionsTable(Db *sql.DB) bool {
	// Define the SQL query to create the table
	query := `
		CREATE TABLE IF NOT EXISTS transactions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			walletId INT NOT NULL,
			status VARCHAR(50),
			type VARCHAR(50) NOT NULL,
			createdAt DATETIME NOT NULL,
			amount DOUBLE,
			FOREIGN KEY (walletId) REFERENCES wallets(id)
		)	
	`

	// Execute the SQL query to create the table
	_, err := Db.Exec(query)
	if err != nil {
		return false
	}

	return true
}
