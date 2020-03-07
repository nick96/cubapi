package user

import (
	"log"
	"os"

	"github.com/nick96/cubapi/db"
	"go.uber.org/zap"
)

var authService AuthService
var store UserStorer

type mockUserStore map[int64]User

func newMockUserStore() mockUserStore {
	return make(map[int64]User)
}

func (s mockUserStore) FindByEmail(email string) (User, bool, error) {
	for _, user := range s {
		if user.Email == email {
			return user, true, nil
		}
	}
	return User{}, false, nil
}

func (s mockUserStore) AddUser(user User) (int64, error) {
	nextID := int64(1)
	for _, user := range s {
		if user.Id > nextID {
			nextID = user.Id + 1
		}
	}
	user.Id = nextID
	s[nextID] = user
	return nextID, nil
}

func getStore() UserStorer {
	return store
}

func getAuthService() AuthService {
	return authService
}

func withMockUserStorer() {
	store = mockUserStore(make(map[int64]User))
	authService = AuthService{store}
}

func withUserStore() {
	dbHandle, err := db.NewConn(
		zap.NewNop(),
		os.Getenv("USER_DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("USER_DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_SSL_MODE"),
	)
	if err != nil {
		log.Fatal(err)
	}
	store = UserStore{dbHandle}
	authService = AuthService{store}
}

func cleanStore() {
	if _, ok := store.(UserStore); ok {
		log.Printf("Deleting all rows in users table")
		store.(UserStore).db.MustExec(`DELETE FROM users;`)
	} else {
		store = mockUserStore(make(map[int64]User))
		authService = AuthService{store}
	}
}
