package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 5

var db *sqlx.DB
var ctx = context.Background()

// New is the function used to create an instance of the data package.
// It returns the type Model, which embeds all of the types we want to
// be available to our application.
func New(dbPool *sqlx.DB) Models {
	db = dbPool

	return Models{
		User:  User{},
		Token: Token{},
	}
}

// Models is the type for this package. Note that any model that is
// included as a member in this type is available to us throughout the
// application, anywhere that the app variable is used, provided that the
// model is also added in the New function.
type Models struct {
	User  User
	Token Token
}

// define type for NULL from database
type NullTime struct {
	mysql.NullTime
}

type NullString struct {
	sql.NullString
}

// User is the structure which holds one user from the database. Note
// that it embeds a token type.
type User struct {
	ID              int       `db:"id" json:"id"`
	UserID          string    `db:"user_id" json:"user_id" validate:"required"`
	FirstName       string    `db:"first_name" json:"first_name,omitempty"`
	LastName        string    `db:"last_name" json:"last_name,omitempty"`
	Email           string    `db:"email" json:"email,omitempty" validate:"required"`
	Password        string    `db:"password" json:"password,omitempty" validate:"required"`
	EmailVerifiedAt NullTime  `db:"email_verified_at" json:"email_verified_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt       NullTime  `db:"deleted_at" json:"deleted_at"`
	Token           Token
}

// Token is the data structure for any token in the database. Note that
// we do not send the TokenHash (a slice of bytes) in any exported JSON.
type Token struct {
	ID        int       `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	TokenHash []byte    `db:"-" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	ExpireAt  time.Time `db:"expire_at" json:"expire_at"`
}

func generateUUID() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	return uuid
}

// Insert inserts a new user into the database, and returns the ID of the
// newly inserted row
func (u *User) Insert(user User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var newID string

	uuid := generateUUID()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return "", err
	}

	stmt := `
		INSERT INTO
			users
			(
				user_id,
				first_name,
				last_name,
				email,
				password,
				created_at,
				updated_at
			)
			VALUES (
				?,
				?,
				?,
				?,
				?,
				?,
				?
			)
	`

	_, err = db.ExecContext(ctx, stmt,
		uuid,
		user.FirstName,
		user.LastName,
		user.Email,
		hashedPassword,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return "", err
	}

	newID = uuid

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
		UPDATE
			users
		SET
			password = $1
		WHERE
			id = $2
	`

	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the
// password and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// Index returns a slice of all users
func (u *User) Index() ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		select
			*
		from
			users
	`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return users, nil
}

// Show returns one user by id
func (u *User) ShowByID(userID string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		SELECT
			*
		FROM
			users
		WHERE
			user_id = ?
	`

	var user User
	row := db.QueryRowxContext(ctx, query, userID)

	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail returns one user by email
func (u *User) ShowByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		SELECT
			*
		FROM
			users
		WHERE
			email = ?
	`

	var user User
	row := db.QueryRowxContext(ctx, query, email)

	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	stmt := `
		UPDATE
			users
		SET
			email = $1,
			first_name = $2,
			last_name = $3,
			updated_at = $4
		WHERE
			id = $5
	`

	var user User
	_, err := db.ExecContext(ctx, stmt,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		time.Now(),
		&user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes one user from the database, by ID
func (u *User) Delete(userID string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	stmt := `
		DELETE from
			users
		WHERE
			user_id = ?
	`

	_, err := db.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetByToken takes a plain text token string, and looks up the full token
// from the database. It returns a pointer to the Token model.
func (t *Token) GetByToken(plainText string) (*Token, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		SELECT
			*
		FROM
			tokens
		WHERE
			token = ?
	`

	var token Token

	row := db.QueryRowxContext(ctx, query, plainText)

	err := row.StructScan(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetUserForToken takes a token parameter, and uses the UserID field from that
// parameter to look a user up by id. It returns a pointer to the user model.
func (t *Token) GetUserForToken(token Token) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
		SELECT
			*
		FROM
			users
		WHERE
			id = ?
	`

	var user User
	row := db.QueryRowxContext(ctx, query, token.UserID)

	err := row.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GenerateToken generates a secure token of exactly 26 characters in length and returns it
func (t *Token) GenerateToken(userID string, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID:   userID,
		ExpireAt: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]

	return token, nil
}

// AuthenticateToken takes the full http request, extracts the authorization header,
// takes the plain text token from that header and looks up the associated token entry
// in the database, and then finds the user associated with that token. If the token
// is valid and a user is found, the user is returned; otherwise, it returns an error.
func (t *Token) AuthenticateToken(r *http.Request) (*User, error) {
	// get the authorization header
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	// get the plain text token from the header
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no valid authorization header received")
	}

	// make sure the token is of the correct length
	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("token wrong size")
	}

	// get the token from the database, using the plain text token to find it
	tkn, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("no matching token found")
	}

	// make sure the token has not expired
	if tkn.ExpireAt.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	// get the user associated with the token
	user, err := t.GetUserForToken(*tkn)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

// Insert inserts a token into the database
func (t *Token) Insert(token Token, u User) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// delete any existing tokens
	stmt := `
		DELETE FROM
			tokens
		WHERE
			user_id = ?
	`
	_, err := db.ExecContext(ctx, stmt, token.UserID)
	if err != nil {
		return err
	}

	// we assign the email value, just to be safe, in case it was
	// not done in the handler that calls this function
	// token.Email = u.Email

	// insert the new token

	uuid := generateUUID()

	// fmt.Println(token)

	stmt = `
		INSERT INTO
			tokens (
				id,
				user_id,
				token,
				token_hash,
				created_at,
				updated_at,
				expire_at
			)
			VALUES (
				?,
				?,
				?,
				?,
				?,
				?,
				?
			)
	`

	// _, err = db.NamedExecContext(ctx, stmt, new)

	fmt.Println(uuid)
	fmt.Println(token.UserID)
	fmt.Println(token.Token)
	fmt.Println(token.TokenHash)
	fmt.Println(time.Now())
	fmt.Println(time.Now())
	fmt.Println(token.ExpireAt)

	_, err = db.ExecContext(ctx, stmt,
		uuid,
		token.UserID,
		token.Token,
		token.TokenHash,
		time.Now(),
		time.Now(),
		token.ExpireAt,
	)
	if err != nil {
		return err
	}

	return nil

	// ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	// defer cancel()

	// var newID string

	// uuid := generateUUID()

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	// if err != nil {
	// 	return "", err
	// }

	// stmt := `
	// 	INSERT INTO
	// 		users
	// 		(
	// 			user_id,
	// 			first_name,
	// 			last_name,
	// 			email,
	// 			password,
	// 			created_at,
	// 			updated_at
	// 		)
	// 		VALUES (
	// 			:user_id,
	// 			:first_name,
	// 			:last_name,
	// 			:email,
	// 			:password,
	// 			:created_at,
	// 			:updated_at
	// 		)
	// `

	// type newUser struct {
	// 	UserID    string    `db:"user_id" json:"user_id" validate:"required"`
	// 	FirstName string    `db:"first_name" json:"first_name,omitempty"`
	// 	LastName  string    `db:"last_name" json:"last_name,omitempty"`
	// 	Email     string    `db:"email" json:"email,omitempty" validate:"required"`
	// 	Password  []byte    `db:"password" json:"password,omitempty" validate:"required"`
	// 	CreatedAt time.Time `db:"created_at" json:"created_at"`
	// 	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	// }

	// new := newUser{
	// 	UserID:    uuid,
	// 	FirstName: user.FirstName,
	// 	LastName:  user.LastName,
	// 	Email:     user.Email,
	// 	Password:  hashedPassword,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }

	// if err != nil {
	// 	return "", err
	// }

	// newID = uuid

	// return newID, nil

}

// DeleteByToken deletes a token, by plain text token
func (t *Token) DeleteByToken(plainText string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	stmt := `
		DELETE FROM
			tokens
		WHERE
			token = ?
	`

	_, err := db.ExecContext(ctx, stmt, plainText)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteTokensForUser(id int) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	stmt := `
		DELETE FROM
			tokens
		WHERE
			user_id = ?
	`
	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// ValidToken makes certain that a given token is valid; in order to be valid,
// the token must exist in the database, the associated user must exist in the
// database, and the token must not have expired.
func (t *Token) ValidToken(plainText string) (bool, error) {
	token, err := t.GetByToken(plainText)
	if err != nil {
		return false, errors.New("no matching token found")
	}

	_, err = t.GetUserForToken(*token)
	if err != nil {
		return false, errors.New("no matching user found")
	}

	if token.ExpireAt.Before(time.Now()) {
		return false, errors.New("expired token")
	}

	return true, nil
}
