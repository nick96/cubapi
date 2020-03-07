package user

// User is a representation of a user entity.
type User struct {
	Id        int64  `json:"id,omitempty" db:"id"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"firstName" db:"firstname"`
	LastName  string `json:"lastName" db:"lastname"`
	Password  string `json:"-" password:"password"`
	Salt      string `json:"-" salt:"salt"`
}
