package user

import (
	"log/slog"
	"math"
	"math/rand/v2"

	"github.com/brianvoe/gofakeit/v7"
)

const (
	HttpUserCreatePath = "/users/create"
)

// User represents a speeddate user
type User struct {
	Id             int    `json:"id,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password,omitempty"`
	Name           string `json:"name,omitempty"`
	Gender         string `json:"gender,omitempty"`
	Age            int    `json:"age,omitempty"`
	longitude      float64
	latitude       float64
	DistanceFromMe int `json:"distanceFromMe"`
}

// NewRandomUser generates a random user
func NewRandomUser() User {
	return User{
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, true, false, 10),
		Gender:    gofakeit.Gender(),
		Age:       rand.IntN(100-18) + 18,
		longitude: gofakeit.Longitude(),
		latitude:  gofakeit.Latitude(),
	}
}

// CalculateDistance uses the pythagorean formula to calculate the distance between some coordinates and a user
func (u *User) CalculateDistance(fromLat, fromLong float64) {
	longDistance := math.Max(fromLong, u.longitude) - math.Min(fromLong, u.longitude)*54.6 // 54.6 miles in a degree of longitude
	latDistance := math.Max(fromLat, u.latitude) - math.Min(fromLat, u.latitude)*69.0      // 69 miles in a degree of latitude, Nice!

	distance := math.Sqrt((longDistance * longDistance) + (latDistance * latDistance))
	slog.Info("distance calculated", "distance", int(distance))

	u.DistanceFromMe = int(distance)
}

// ByDistance is a utility type to satisfy the sort.Interface
type ByDistance []User

func (d ByDistance) Len() int           { return len(d) }
func (d ByDistance) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByDistance) Less(i, j int) bool { return d[i].DistanceFromMe < d[j].DistanceFromMe }
