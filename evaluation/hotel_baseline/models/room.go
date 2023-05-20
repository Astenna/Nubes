package models

import (
	"errors"
	"fmt"
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
	DateIn    string
	DateOut   string
	UserEmail string
	HotelName string
	CityName  string
	RoomId    string
}

func ReserveRoom(param ReserveParam) error {
	dateOut, err1 := time.Parse("2006-01-02", param.DateOut)
	dateIn, err2 := time.Parse("2006-01-02", param.DateIn)
	if err1 != nil || err2 != nil {
		return errors.New("failed to parse DateIn or DateOut, the dates must be specified as strings in format YYYY-MM-DD")
	}

	if dateOut.Before(dateIn) {
		return fmt.Errorf("dateOut can not be before DateIn")
	}

	// STEP 1: ensure the room is available
	dbParam := db.GetItemBySortKey{
		PkName:    "CityHotelRoomId",
		PkValue:   GetReservationPK(param.CityName, param.HotelName, param.RoomId),
		SkName:    "DateIn",
		SkValue:   dateIn,
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

	var r1, r2, r3, r4 bool

	beforeExists := true
	if resBefore == nil {
		beforeExists = false
	} else {
		r1 = dateIn.Before(resBefore.DateOut)
		r2 = dateIn.Format("20060102") != resBefore.DateOut.Format("20060102")
	}

	afterExists := true
	if resAfter == nil {
		afterExists = false
	} else {
		r3 = dateOut.After(resAfter.DateIn)
		r4 = dateOut.Format("20060102") != resAfter.DateIn.Format("20060102")
	}

	if (beforeExists && (r1 && r2)) || (afterExists && (r3 || r4)) {
		return errors.New("room already booked")
	}

	// STEP 2: create reservation
	reservationPK := GetReservationPK(param.CityName, param.HotelName, param.RoomId)
	err = db.Insert(Reservation{
		CityHotelRoomId: reservationPK,
		DateIn:          dateIn,
		DateOut:         dateOut,
	}, db.ReservationTable)

	if err != nil {
		return err
	}

	err = db.Insert(db.UserReservationsJoinTableEntry{
		UserId:          param.UserEmail,
		CityHotelRoomId: reservationPK,
		DateIn:          dateIn,
	}, db.UserResevationsJoinTable)
	return err
}

func GetReservationPK(cityName, hotelName, roomId string) string {
	return cityName + "_" + hotelName + "_" + roomId
}
