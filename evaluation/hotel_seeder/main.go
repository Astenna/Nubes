package main

import (
	"strconv"
	"time"

	"github.com/Astenna/Nubes/evaluation/hotel/types"
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/db"
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/Astenna/Nubes/lib"
	"github.com/google/uuid"
	"github.com/jftuga/geodist"
)

const UserCount = 50
const CitiesCount = 5
const HotelsPerCity = 20
const RoomsPerHotel = 5

const CityPrefix = "Milano_"
const HotelPrefix = "Bruschetti_"
const ReservationYear = 2023

func SeedUsers() {

	for i := 0; i < UserCount; i++ {
		suffix := strconv.Itoa(i)
		// baseline
		userb := models.User{
			Email:     "Email_" + suffix,
			FirstName: "Cornell_" + suffix,
			LastName:  "Baker_" + suffix,
			Password:  "Password_" + suffix,
		}
		insert(userb, db.UserTable)
		// nubes
		user := types.User{
			FirstName: "Cornell_" + suffix,
			LastName:  "Baker_" + suffix,
			Email:     "Email_" + suffix,
			Password:  "Password_" + suffix,
		}
		insert(user, user.GetTypeName())
	}
}

func SeedCities() {

	for i := 0; i < CitiesCount; i++ {
		suffix := strconv.Itoa(i)
		// baseline
		cityb := models.City{}
		cityb.CityName = CityPrefix + suffix
		cityb.HotelName = CityPrefix + suffix
		cityb.Region = "Lombardia" + suffix
		cityb.Description = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce scelerisque eu risus non lacinia. Nullam at ligula gravida, vehicula justo ac, feugiat est. Fusce hendrerit, orci sed fermentum molestie, odio felis laoreet tellus, non vulputate urna diam eu nibh. Etiam quis pharetra sem. Sed non lorem id lacus pellentesque egestas vel vitae metus. Quisque at magna massa. Praesent viverra velit dui, ac porta libero molestie sed. `
		insert(cityb, db.CityTable)
		// nubes
		city := types.City{}
		city.CityName = CityPrefix + suffix
		city.Region = "Lombardia" + suffix
		city.Description = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce scelerisque eu risus non lacinia. Nullam at ligula gravida, vehicula justo ac, feugiat est. Fusce hendrerit, orci sed fermentum molestie, odio felis laoreet tellus, non vulputate urna diam eu nibh. Etiam quis pharetra sem. Sed non lorem id lacus pellentesque egestas vel vitae metus. Quisque at magna massa. Praesent viverra velit dui, ac porta libero molestie sed. `
		insert(city, city.GetTypeName())
	}
}

func SeedHotels() {

	for i := 0; i < CitiesCount; i++ {
		citySuffix := strconv.Itoa(i)

		for j := 0; j < HotelsPerCity; j++ {
			hotelSuffix := strconv.Itoa(j)

			// baseline
			hotelb := models.Hotel{
				CityName:   CityPrefix + citySuffix,
				HotelName:  HotelPrefix + hotelSuffix,
				Street:     "AwesomeStreet" + hotelSuffix,
				PostalCode: hotelSuffix,
				Coordinates: geodist.Coord{
					Lat: float64(j%91) - 21.43,
					Lon: float64(j%181) - 12.45,
				},
				Rate: float32(j % 6),
			}
			insert(hotelb, db.HotelTable)
			// nubes
			hotel := types.Hotel{
				HName:      HotelPrefix + hotelSuffix,
				Street:     "AwesomeStreet" + hotelSuffix,
				PostalCode: hotelSuffix,
				Coordinates: geodist.Coord{
					Lat: float64(j%91) - 21.43,
					Lon: float64(j%181) - 12.45,
				},
				Rate: float32(j % 6),
				City: *lib.NewReference[types.City](CityPrefix + citySuffix),
			}
			insert(hotel, hotel.GetTypeName())
		}
	}
}

func SeedRoomsAndReservations() {
	const ReservationsPerRoom = 40

	for c := 0; c < CitiesCount; c++ {
		citySuffix := strconv.Itoa(c)
		for j := 0; j < HotelsPerCity; j++ {
			hotelSuffix := strconv.Itoa(j)

			for i := 0; i < RoomsPerHotel; i++ {
				roomSuffix := strconv.Itoa(i)

				// baseline
				roomb := models.Room{
					CityHotelName: CityPrefix + citySuffix + "_" + HotelPrefix + hotelSuffix,
					RoomId:        "Room_" + roomSuffix,
					Name:          "Room_" + roomSuffix,
					Description:   `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris mi, consequat quis dapibus eu, ullamcorper non metus. Suspendisse sit amet faucibus nisi. Nullam pharetra libero ut dui facilisis semper.`,
					Price:         float32(i) + 1,
				}

				for k := 0; k < ReservationsPerRoom; k++ {
					dateIn := time.Date(ReservationYear, 1, k*8, 0, 0, 0, 0, time.UTC)

					reservationb := models.Reservation{
						CityHotelRoomId: models.GetReservationPK(CityPrefix+citySuffix, HotelPrefix+hotelSuffix, "Room_"+roomSuffix),
						DateIn:          dateIn,
						DateOut:         dateIn.AddDate(0, 0, int(k%8)),
						UserEmail:       "Email_" + strconv.Itoa(int(k%UserCount)),
					}
					insert(reservationb, db.ReservationTable)
				}
				insert(roomb, db.RoomTable)

				// nubes
				room := types.Room{
					Id:           uuid.New().String(),
					Name:         "Room_" + roomSuffix,
					Description:  `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris mi, consequat quis dapibus eu, ullamcorper non metus. Suspendisse sit amet faucibus nisi. Nullam pharetra libero ut dui facilisis semper.`,
					Hotel:        lib.Reference[types.Hotel](HotelPrefix + hotelSuffix),
					Reservations: map[string][]types.ReservationInOut{},
					Price:        float32(i),
				}

				insert(room, room.GetTypeName())
				for k := 0; k < ReservationsPerRoom; k++ {
					dateIn := time.Date(ReservationYear, 1, k*8, 0, 0, 0, 0, time.UTC)

					param := types.ReserveParam{
						DateIn:                dateIn,
						DateOut:               dateIn.AddDate(0, 0, int(k%8)),
						User:                  lib.Reference[types.User]("Email_" + strconv.Itoa(int(k%UserCount))),
						RoomId:                room.Id,
						SkipAvailabilityCheck: true,
					}
					types.ExportReservation(param)
				}
			}
		}
	}
}

func main() {
	SeedUsers()
	SeedCities()
	SeedHotels()
	SeedRoomsAndReservations()
}
