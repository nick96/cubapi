package attendance

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

// AttendanceStorer is the interface an object that stores attendance must
// implement.
type AttendanceStorer interface {
	// GetAll should get all the attendances store in the store.
	GetAll() ([]Attendance, error)
	// GetByID should get the attendance corresponding to the given ID.
	GetByID(attendanceID uint64) (Attendance, error)
	// GetByCub should get all the attendances corresponding to the
	// given cub ID.
	GetByCub(cubID uint64) ([]Attendance, error)
	// GetByDate should get all the attendances on a given date.
	GetByDate(date time.Time) ([]Attendance, error)

	// Insert should insert an attendance into the store and return its ID.
	Insert(attendance Attendance) (uint64, error)

	// Update should update the attendance..
	Update(attendance Attendance) error

	// Delete should delete the attendance.
	Delete(attendance Attendance) error
}

// AttendanceStore is a concrete implementation of `AttendanceStorer` that
// accesses a real database.
type AttendanceStore struct {
	db *sqlx.DB
}

func NewAttendanceStore(db *sqlx.DB) AttendanceStore {
	return AttendanceStore{db}
}

func (s AttendanceStore) GetAll() ([]Attendance, error) {
	return []Attendance{}, nil
}

func (s AttendanceStore) GetByID(attendanceID uint64) (Attendance, error) {
	return Attendance{}, nil
}

func (s AttendanceStore) GetByCub(cubID uint64) ([]Attendance, error) {
	return []Attendance{}, nil
}

func (s AttendanceStore) GetByDate(date time.Time) ([]Attendance, error) {
	return []Attendance{}, nil
}

func (s AttendanceStore) Insert(attendance Attendance) (uint64, error) {
	return 0, nil
}

func (s AttendanceStore) Update(attendance Attendance) error {
	return nil
}

func (s AttendanceStore) Delete(attendance Attendance) error {
	return nil
}

// Attendance represents a single cubs attendance on a given date.
type Attendance struct {
	Model

	// Date of attendance.
	Date time.Time `json:"date"`
	// Cub is the cub the attendance entity applies to.
	Cub Cub
}

type Handler struct {
	cubStore        CubStore
	attendanceStore AttendanceStore
}

func NewHandler(cubStore CubStore, attendanceStore AttendanceStore) Handler {
	return Handler{cubStore, attendanceStore}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
