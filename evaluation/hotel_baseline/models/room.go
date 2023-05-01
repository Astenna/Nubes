package models

import (
	"errors"
	"time"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/db"
)

type Room struct {
	CityHotelName string // CityName_HotelName
	RoomId        string
	Name          string
	Description   string
	Price         float32
}

type ReserveParam struct {
	DateIn    time.Time
	DateOut   time.Time
	UserEmail string
	HotelName string
	CityName  string
	RoomId    string
}

func ReserveRoom(param ReserveParam) error {
	// STEP 1: ensure the room is available
	dbParam := db.GetItemBySortKey{
		PkName:    "RoomId",
		PkValue:   GetReservationPK(param.CityName, param.HotelName, param.RoomId),
		SkName:    "DateIn",
		SkValue:   param.DateIn,
		TableName: db.ReservationTable,
	}
	resBefore, err := db.GetItemBeforeSortKey[Reservation](dbParam)
	if err != nil {
		return err
	}
	resAfter, err := db.GetItemAfterSortKey[Reservation](dbParam)
	if err != nil {
		return err
	}

	r1 := param.DateIn.Before(resBefore.DateOut)
	r2 := param.DateIn.Format("20060102") != resBefore.DateOut.Format("20060102")
	r3 := param.DateOut.After(resAfter.DateIn)
	r4 := param.DateOut.Format("20060102") != resAfter.DateIn.Format("20060102")
	if (r1 && r2) || (r3 || r4) {
		return errors.New("room already booked")
	}

	// STEP 2: create reservation
	err = db.Insert(Reservation{
		CityHotelNameRoomId: GetReservationPK(param.CityName, param.HotelName, param.RoomId),
		DateIn:              param.DateIn,
		DateOut:             param.DateOut,
		UserEmail:           param.UserEmail,
	}, db.ReservationTable)
	return err
}

func GetReservationPK(cityName, hotelName, roomId string) string {
	return cityName + "_" + hotelName + "_" + roomId
}
