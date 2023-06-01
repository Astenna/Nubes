package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Astenna/Nubes/evaluation/hotel/types"
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/db"
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	db_simple "github.com/Astenna/Nubes/evaluation/hotel_baseline_simple/db"
	models_simple "github.com/Astenna/Nubes/evaluation/hotel_baseline_simple/models"
	"github.com/Astenna/Nubes/lib"
	"github.com/jftuga/geodist"
	"golang.org/x/sync/semaphore"
)

const UserCount = 50000
const CitiesCount = 5
const HotelsPerCity = 100
const RoomsPerHotel = 25
const ReservationsPerRoom = 20

const CityPrefix = "Milano"
const HotelPrefix = "Bruschetti"
const ReservationYear = 2023

const K = 32

func SeedUsers() {

	fmt.Println("Seeding USERS")
	var wg sync.WaitGroup

	for j := 0; j < K; j++ {
		wg.Add(1)
		i_start := j
		go func() {
			defer wg.Done()
			for i := i_start; i < UserCount; i += K {
				suffix := strconv.Itoa(i)
				// baseline
				userb := models.User{
					Email:     "Email" + suffix,
					FirstName: "Cornell" + suffix,
					LastName:  "Baker" + suffix,
					Password:  "Password" + suffix,
				}
				insert(userb, db.UserTable)
				// nubes
				user := types.User{
					FirstName: "Cornell" + suffix,
					LastName:  "Baker" + suffix,
					Email:     "Email" + suffix,
					Password:  "Password" + suffix,
				}
				insert(user, user.GetTypeName())
			}
		}()
	}

	wg.Wait()
}

func SeedCities() {

	fmt.Println("Seeding CITIES")

	for i := 0; i < CitiesCount; i++ {
		suffix := strconv.Itoa(i)
		// baseline
		cityb := models.City{}
		cityb.CityName = CityPrefix + suffix
		cityb.HotelName = CityPrefix + suffix
		cityb.Region = "Lombardia" + suffix
		cityb.Description = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce scelerisque eu risus non lacinia. Nullam at ligula gravida, vehicula justo ac, feugiat est. Fusce hendrerit, orci sed fermentum molestie, odio felis laoreet tellus, non vulputate urna diam eu nibh. Etiam quis pharetra sem. Sed non lorem id lacus pellentesque egestas vel vitae metus. Quisque at magna massa. Praesent viverra velit dui, ac porta libero molestie sed. `
		insert(cityb, db.CityTable)
		// simple baseline
		citybs := models_simple.City{}
		citybs.CityName = CityPrefix + suffix
		citybs.Region = "Lombardia" + suffix
		citybs.Description = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce scelerisque eu risus non lacinia. Nullam at ligula gravida, vehicula justo ac, feugiat est. Fusce hendrerit, orci sed fermentum molestie, odio felis laoreet tellus, non vulputate urna diam eu nibh. Etiam quis pharetra sem. Sed non lorem id lacus pellentesque egestas vel vitae metus. Quisque at magna massa. Praesent viverra velit dui, ac porta libero molestie sed. `
		insert(citybs, db_simple.CityTable)
		// nubes
		city := types.City{}
		city.CityName = CityPrefix + suffix
		city.Region = "Lombardia" + suffix
		city.Description = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce scelerisque eu risus non lacinia. Nullam at ligula gravida, vehicula justo ac, feugiat est. Fusce hendrerit, orci sed fermentum molestie, odio felis laoreet tellus, non vulputate urna diam eu nibh. Etiam quis pharetra sem. Sed non lorem id lacus pellentesque egestas vel vitae metus. Quisque at magna massa. Praesent viverra velit dui, ac porta libero molestie sed. `
		insert(city, city.GetTypeName())
	}
}

func SeedHotels() {

	fmt.Println("Seeding HOTELS")

	for i := 0; i < CitiesCount; i++ {
		citySuffix := strconv.Itoa(i)

		var wg sync.WaitGroup

		for j := 0; j < HotelsPerCity; j++ {
			jj := j

			wg.Add(1)
			go func() {
				defer wg.Done()
				hotelSuffix := strconv.Itoa(jj)

				// baseline
				hotelb := models.Hotel{
					CityName:   CityPrefix + citySuffix,
					HotelName:  HotelPrefix + hotelSuffix,
					Street:     "AwesomeStreet" + hotelSuffix,
					PostalCode: hotelSuffix,
					Coordinates: geodist.Coord{
						Lat: float64(jj%91) - 21.43,
						Lon: float64(jj%181) - 12.45,
					},
					Rate: float32(jj % 6),
				}
				insert(hotelb, db.HotelTable)
				// simple baseline
				hotelbs := models_simple.Hotel{
					CityName:   CityPrefix + citySuffix,
					HotelName:  CityPrefix + citySuffix + "_" + HotelPrefix + hotelSuffix,
					Street:     "AwesomeStreet" + hotelSuffix,
					PostalCode: hotelSuffix,
					Coordinates: geodist.Coord{
						Lat: float64(j%91) - 21.43,
						Lon: float64(j%181) - 12.45,
					},
					Rate: float32(j % 6),
				}
				insert(hotelbs, db_simple.HotelTable)
				// nubes
				hotel := types.Hotel{
					HName:      CityPrefix + citySuffix + "_" + HotelPrefix + hotelSuffix,
					Street:     "AwesomeStreet" + hotelSuffix,
					PostalCode: hotelSuffix,
					Coordinates: geodist.Coord{
						Lat: float64(jj%91) - 21.43,
						Lon: float64(jj%181) - 12.45,
					},
					Rate: float32(jj % 6),
					City: *lib.NewReference[types.City](CityPrefix + citySuffix),
				}
				insert(hotel, hotel.GetTypeName())
			}()
		}

		wg.Wait()
	}
}

func nextUser(q int, mod int) int {
	return (q*364526735 + 23562367) % mod
}

func SeedRoomsAndReservations() {
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(K)
	ctx := context.TODO()

	for c := 0; c < CitiesCount; c++ {
		citySuffix := strconv.Itoa(c)
		fmt.Println("Seeding ROOMS and RESERVATIONS for city " + citySuffix + "out of " + strconv.Itoa(CitiesCount))

		cc := c
		for j := 0; j < HotelsPerCity; j++ {
			hotelSuffix := strconv.Itoa(j)

			jj := j
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer sem.Release(1)
				sem.Acquire(ctx, 1)
				fmt.Println("------------------------------ in hotel " + hotelSuffix + "out of " + strconv.Itoa(HotelsPerCity))

				for i := 0; i < RoomsPerHotel; i++ {
					roomSuffix := strconv.Itoa(i)

					// baseline
					next := nextUser(cc*101+jj*19+i*7, UserCount)
					roomb := models.Room{
						CityHotelName: CityPrefix + citySuffix + "_" + HotelPrefix + hotelSuffix,
						RoomId:        "Room" + roomSuffix,
						Name:          "Room" + roomSuffix,
						Description:   `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris mi, consequat quis dapibus eu, ullamcorper non metus. Suspendisse sit amet faucibus nisi. Nullam pharetra libero ut dui facilisis semper.`,
						Price:         float32(i) + 1,
					}
					insert(roomb, db.RoomTable)

					for k := 0; k < ReservationsPerRoom; k++ {
						dateIn := time.Date(ReservationYear, 1, k*8, 0, 0, 0, 0, time.UTC)

						reservationb := models.Reservation{
							CityHotelRoomId: models.GetReservationPK(CityPrefix+citySuffix, HotelPrefix+hotelSuffix, "Room"+roomSuffix),
							DateIn:          dateIn,
							DateOut:         dateIn.AddDate(0, 0, int(k%8)),
						}
						next = nextUser(next, UserCount)

						insert(reservationb, db.ReservationTable)
						insert(db.UserReservationsJoinTableEntry{
							UserId:          "Email_" + strconv.Itoa(int(k%UserCount)),
							CityHotelRoomId: reservationb.CityHotelRoomId,
							DateIn:          dateIn,
						}, db.UserResevationsJoinTable)
					}

					// nubes
					next = nextUser(cc*101+jj*19+i*7, UserCount)
					room := types.Room{
						Id:           CityPrefix + citySuffix + "_" + HotelPrefix + hotelSuffix + "_Room" + roomSuffix,
						Name:         "Room" + roomSuffix,
						Description:  `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris mi, consequat quis dapibus eu, ullamcorper non metus. Suspendisse sit amet faucibus nisi. Nullam pharetra libero ut dui facilisis semper.`,
						Hotel:        lib.Reference[types.Hotel](HotelPrefix + hotelSuffix),
						Reservations: map[string][]types.ReservationInOut{},
						Price:        float32(i),
					}

					insert(room, room.GetTypeName())
					for k := 0; k < ReservationsPerRoom; k++ {
						dateIn := time.Date(ReservationYear, 1, k*8, 0, 0, 0, 0, time.UTC)

						param := types.ReserveParam{
							DateIn:                dateIn.Format("2006-01-02"),
							DateOut:               dateIn.AddDate(0, 0, int(k%8)).Format("2006-01-02"),
							User:                  lib.Reference[types.User]("Email" + strconv.Itoa(int(next))),
							RoomId:                room.Id,
							SkipAvailabilityCheck: true,
						}
						next = nextUser(next, UserCount)

						types.ExportReservation(param)
						// fmt.Print(",")
					}
				}
			}()
		}
	}
	wg.Wait()
}

func main() {
	SeedUsers()
	SeedCities()
	SeedHotels()
	SeedRoomsAndReservations()
}
