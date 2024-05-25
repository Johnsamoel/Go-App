package db

import (
	"database/sql"
	"example.com/fintech-app/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
	"time"
)

type UserRepository struct{}

var UserRepo = UserRepository{}

func (u UserRepository) CreateNewUser(userData *models.User) (*models.User, error) {
	query := `
        INSERT INTO users (name, email, password, phoneNumber)
        VALUES (?, ?, ?, ?)
    `

	stmt, err := DB.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error creating new user: %v", err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(userData.Name, userData.Email, userData.Password, userData.PhoneNumber)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error creating new user: %v", err)
	}

	// Get the ID of the newly inserted user
	userID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Set the ID of the user in the userData object
	userData.ID = userID

	_, err =  WalletRepo.CreateNewWallet(userID)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error in creating new wallet: %v", err)
	}

	return userData, nil
}

func (u UserRepository) DeleteUser(userId int64) error {

	DeleteUserQuery := `
	DELETE FROM users WHERE id = ?
`

	DeleteWalletQuery := `
	DELETE FROM wallets WHERE userId = ?
`

	// Delete the wallet associated with the user
	userStmt, err := DB.Prepare(DeleteUserQuery)

	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	walletStmt, err := DB.Prepare(DeleteWalletQuery)

	if err != nil {
		return fmt.Errorf("error deleting wallet: %v", err)
	}

	defer userStmt.Close()
	defer walletStmt.Close()

	_, err = walletStmt.Exec(userId)
	if err != nil {
		return fmt.Errorf("error deleting wallet: %v", err)
	}

	// Delete the user
	_, err = userStmt.Exec(userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

func (u UserRepository) UpdateUserData(userId int64, userData map[string]interface{}) error {
	// Prepare the SQL query to update user data
	query := `
        UPDATE users
        SET name = ?, email = ?, password = ?, phoneNumber = ?
        WHERE id = ?
    `

	updateUsersStmt, err := DB.Prepare(query)

	if err != nil {
		return fmt.Errorf("error updating user data: %v", err)
	}

	defer updateUsersStmt.Close()

	user, err := u.GetUserById(userId)

	if err != nil {
		return fmt.Errorf("error updating user data: %v", err)
	}

	// Get the reflect.Value of the user struct
	userValue := reflect.ValueOf(user)

	// Iterate over the fields in the map
	for key, value := range userData {
		// Get the reflect.Value of the field in the user struct
		field := userValue.Elem().FieldByName(key)
		if !field.IsValid() {
			return fmt.Errorf("field %s does not exist in user struct", key)
		}

		// Check if the field is settable and assignable
		if field.CanSet() && field.Type().AssignableTo(reflect.TypeOf(value)) {
			// Set the value of the field
			field.Set(reflect.ValueOf(value))
		} else {
			return fmt.Errorf("field %s is not settable or assignable", key)
		}
	}

	// hashing user password
	err = user.HashPassword()
	if err != nil {
		return fmt.Errorf("error updating user data: %v", err)
	}

	// Execute the SQL query with user data
	result, err := updateUsersStmt.Exec(user.Name, user.Email, user.Password, user.PhoneNumber, userId)
	if err != nil {
		return fmt.Errorf("error updating user data: %v", err)
	}

	// Check if any rows were affected by the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error updating user data: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userId)
	}

	// Set the ID of the user in the userData object
	user.ID = userId

	return nil
}

func (u UserRepository) GetUserById(userId int64) (*models.User, error) {
	// Prepare the SQL query to fetch user data
	query := `
        SELECT * FROM users
        WHERE id = ?
    `

	// Prepare the SQL statement
	userStmt, err := DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}

	defer userStmt.Close()

	// Execute the SQL query with user data
	row := userStmt.QueryRow(userId)

	// Initialize a new User object to store the fetched data
	var user models.User

	// Scan the row into the user object
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	return &user, nil
}

func (u UserRepository) Login(w http.ResponseWriter, r *http.Request, email, password string) (string, error) {
	// Prepare the SQL query with placeholders for email and password
	query := `SELECT * FROM users WHERE email = ?`
	stmt, err := DB.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with the provided email
	row := stmt.QueryRow(email)

	// Initialize a new User object to store the fetched user data
	var user models.User

	// Scan the row into the user object
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found with the given email
			return "", fmt.Errorf("user with email %s not found", email)
		}
		// Other error occurred
		return "", fmt.Errorf("error fetching user: %v", err)
	}

	// Compare the stored password hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Passwords do not match
		return "", fmt.Errorf("invalid Credentials")
	}

	// get user wallet 

	userWallet , err := WalletRepo.GetWalletByUserId(user.ID)

	if err != nil {
		return "", fmt.Errorf("something went wrong")
	}


	var token string
	var expirationDate time.Duration
	var tokenName string
	otpStr := ""

	// implement logic for sending otp to activate the wallet
	if userWallet.Status != "ACTIVE" {
		tokenName = "id_token"
		expirationDate = (5 * time.Minute)
		// generate wallet activation token
		token, err = GenerateToken(user.ID, userWallet.ID ,expirationDate)
		if err != nil {
			return "", fmt.Errorf("something went wrong")
		}

		// check if the user has an active opt

		userActiveOTP, err := OTPRepo.HasActiveOTP(user.ID)

		if err != nil {
			return "", fmt.Errorf("something went wrong")
		}

		if userActiveOTP == "" {
			// generate the otp
			otp, err := OTPRepo.GenerateOtp(user.ID, userWallet.ID)
			if err != nil {
				return "", fmt.Errorf("something went wrong")
			}

			otpStr = otp.Otp
		} else {
			otpStr = userActiveOTP
		}

	} else {
		// generate token
		token, err = GenerateToken(user.ID, userWallet.ID ,time.Hour)
		if err != nil {
			return "", fmt.Errorf("something went wrong")
		}
		tokenName = "id_token"
		expirationDate = 24 * time.Hour

	}

	// Set the token in an HTTP-only cookie
	cookie := &http.Cookie{
		Name:     tokenName,
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(expirationDate),
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	return otpStr, nil
}

func (u UserRepository) Logout(w http.ResponseWriter) {
	// Set the MaxAge of the token cookie to a negative value to expire immediately
	expiredCookie := &http.Cookie{
		Name:     "id_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, expiredCookie)
}
