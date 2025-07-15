package auth

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	PassHash []byte `db:"pass_hash"`
}
