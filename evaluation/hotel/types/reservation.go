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
	Users           lib.ReferenceNavigationList[User] `nubes:"hasMany-Reservations" dynamodbav:"-"`
	DateIn          time.Time
	DateOut         time.Time
	isInitialized   bool
	invocationDepth int
}

func (o Reservation) GetTypeName() string {
	return "Reservation"
}

type ReserveParam struct {
	// date in format YYYY-MM-DD
	DateIn                string
	DateOut               string
	User                  lib.Reference[User]
	RoomId                string
	SkipAvailabilityCheck bool
}

func ExportReservation(param ReserveParam) (string, error) {
	dateOut, err1 := time.Parse("2006-01-02", param.DateOut)
	dateIn, err2 := time.Parse("2006-01-02", param.DateIn)
	if err1 != nil || err2 != nil {
		return "", errors.New("failed to parse DateIn or DateOut, the dates must be specified as strings in format YYYY-MM-DD")
	}

	if dateOut.Before(dateIn) {
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
		for tempStart := dateIn; tempStart.Before(dateOut) && !isAlreadyBooked; {
			inKey := getYearAndMonth(tempStart)
			reservations := roomReservations[inKey]
			tempOutDay := dateOut.Day()
			if tempStart.Month() != dateOut.Month() {
				tempOutDay = dateOut.AddDate(0, 1, -tempOutDay).Day() + tempOutDay
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
	for tempStart := dateIn; tempStart.Before(dateOut); {
		inKey := getYearAndMonth(tempStart)
		tempOutDay := dateOut.Day()
		if tempStart.Month() != dateOut.Month() {
			tempOutDay = dateOut.AddDate(0, 1, -tempOutDay).Day() + tempOutDay
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
		DateIn:  dateIn,
		DateOut: dateOut,
	})
	if err != nil {
		return "", err
	}

	err = res.Users.AddToManyToMany(param.User.Id())
	if err != nil {
		fmt.Println(err)
	}
	return res.Id, err
}

func getYearAndMonth(date time.Time) string {
	return fmt.Sprintf("%d-%02d", date.Year(), date.Month())
}
func (receiver *Reservation) Init() {
	receiver.isInitialized = true
	receiver.Users = *lib.NewReferenceNavigationList[User](lib.ReferenceNavigationListParam{OwnerId: receiver.Id, OwnerTypeName: receiver.GetTypeName(), OtherTypeName: (*new(User)).GetTypeName(), ReferringFieldName: "Users", IsManyToMany: true})
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
