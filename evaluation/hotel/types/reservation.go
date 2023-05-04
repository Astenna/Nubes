package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/Astenna/Nubes/lib"
)

type Reservation struct {
	Id              string
	Room            lib.Reference[Room]
	User            lib.Reference[User] `dynamodbav:",omitempty"`
	DateIn          time.Time
	DateOut         time.Time
	isInitialized   bool
	invocationDepth int
}

func (o Reservation) GetTypeName() string {
	return "Reservation"
}

type ReserveParam struct {
	DateIn                time.Time
	DateOut               time.Time
	User                  lib.Reference[User]
	RoomId                string
	SkipAvailabilityCheck bool
}

func ExportReservation(param ReserveParam) (string, error) {
	if param.DateOut.Before(param.DateIn) {
		return "", fmt.Errorf("dateOut can not be before DateIn")
	}

	// STEP 1: ensure the room is available
	room, err := lib.Load[Room](param.RoomId)
	if err != nil {
		return "", err
	}
	roomReservations, err := room.GetReservations()
	if err != nil {
		return "", err
	}
	if roomReservations == nil {
		roomReservations = make(map[string][]ReservationInOut)
	}

	if !param.SkipAvailabilityCheck {

		isAlreadyBooked := false
		for tempStart := param.DateIn; tempStart.Before(param.DateOut) && !isAlreadyBooked; {
			inKey := getYearAndMonth(tempStart)
			reservations := roomReservations[inKey]
			tempOutDay := param.DateOut.Day()
			if tempStart.Month() != param.DateOut.Month() {
				tempOutDay = param.DateOut.AddDate(0, 1, -tempOutDay).Day() + tempOutDay
			}

			for _, r := range reservations {
				if r.In < tempOutDay && r.Out > tempStart.Day() {
					isAlreadyBooked = true
					break
				}
			}

			tempStart = tempStart.AddDate(0, 1, -tempStart.Day()+1)
		}

		if isAlreadyBooked {
			return "", errors.New("room not available in selected time")
		}
	}

	// STEP 2: insert reservation info to the auxiliary struct
	for tempStart := param.DateIn; tempStart.Before(param.DateOut); {
		inKey := getYearAndMonth(tempStart)
		tempOutDay := param.DateOut.Day()
		if tempStart.Month() != param.DateOut.Month() {
			tempOutDay = param.DateOut.AddDate(0, 1, -tempOutDay).Day() + tempOutDay
		}

		if roomReservations[inKey] == nil {
			roomReservations[inKey] = []ReservationInOut{}
		}
		roomReservations[inKey] = append(roomReservations[inKey], ReservationInOut{In: tempStart.Day(), Out: tempOutDay})
		tempStart = tempStart.AddDate(0, 1, -tempStart.Day()+1)
	}
	err = room.SetReservations(roomReservations)
	if err != nil {
		return "", err
	}

	// STEP 3: insert new reservation
	res, err := lib.Export[Reservation](Reservation{
		Room:    lib.Reference[Room](room.Id),
		User:    param.User,
		DateIn:  param.DateIn,
		DateOut: param.DateOut,
	})
	if err != nil {
		return "", err
	}
	return res.Id, err
}

func getYearAndMonth(date time.Time) string {
	return fmt.Sprintf("%d-%02d", date.Year(), date.Month())
}

func (receiver *Reservation) Init() {
	receiver.isInitialized = true
}
func (receiver *Reservation) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}
