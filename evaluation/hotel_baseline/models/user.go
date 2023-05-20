package models

import (
	"errors"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/db"
)

type User struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}

func Login(email, password string) error {
	user, err := db.GetSingleItemByPartitonKey[User](db.UserTable, "Email", email)

	if err != nil {
		return err
	}

	if user.Password == password {
		return nil
	}

	return errors.New("incorrect password")
}

func DeleteUser(email, password string) error {
	user, err := db.GetSingleItemByPartitonKey[User](db.UserTable, "Email", email)
	if err != nil {
		return err
	}

	if user.Password == password {
		return db.DeleteByPartitionKey(email, "Email", db.UserTable)
	}

	return errors.New("invalid password")
}

func RegisterUser(user User) error {
	return db.Insert(user, db.UserTable)
}

func GetUserReservations(userEmail string) ([]Reservation, error) {
	reservationIds, err := db.GetUserReservationsCompositeKeys(userEmail)
	if err != nil {
		return nil, err
	}

	if len(reservationIds) > 0 {
		return db.GetBatchItemsUsingCompositeKeys[Reservation](reservationIds, db.ReservationTable, "CityHotelRoomId", "DateIn")
	}

	return []Reservation{}, nil
}
