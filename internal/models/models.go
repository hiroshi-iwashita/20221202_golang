package models

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package.
// It returns the type Model, which embeds all of the types we want to
// be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
		// Token: Token{},
	}
}

// Models is the type for this package. Note that any model that is
// included as a member in this type is available to us throughout the
// application, anywhere that the app variable is used, provided that the
// model is also added in the New function.
type Models struct {
	User User
	// Token Token
}

// User is the structure which holds one user from the database. Note
// that it embeds a token type.
type User struct {
	ID              int       `json:"id"`
	UserID          string    `json:"user_id" validate:"required"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email" validate:"required"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	// TokenID   Token     `json:"token_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Token is the data structure for any token in the database. Note that
// we do not send the TokenHash (a slice of bytes) in any exported JSON.
// type Token struct {
// 	ID        int       `json:"id"`
// 	UserID    int       `json:"user_id"`
// 	Email     string    `json:"email"`
// 	Token     string    `json:"token"`
// 	TokenHash []byte    `json:"-"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// 	Expiry    time.Time `json:"expiry"`
// }

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select
			*
		from
			users
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.EmailVerifiedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
