package models

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline_simple/db"
	"github.com/jftuga/geodist"
)

type Hotel struct {
	HotelName   string
	CityName    string
	Street      string
	PostalCode  string
	Coordinates geodist.Coord `nubes:"readonly"`
	Rate        float32
}

func GetHotelsInCity(city string) ([]Hotel, error) {
	hotelIds, err := db.GetHotelIdsInCityByIndex(city)

	if err != nil {
		return nil, err
	}
	return db.GetItemsByPartitonKeys[Hotel](db.HotelTable, "CityName", hotelIds)
}

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

func RecommendHotelsLocation(city string, coordinates Coordinates, count int) ([]Hotel, error) {
	hotels, err := GetHotelsInCity(city)
	if err != nil {
		return nil, err
	}

	result := make([]Hotel, count)

	if err != nil {
		return nil, err
	}

	if len(hotels) <= count {
		return hotels, err
	}
	hotelDists := make([]hotelDist, len(hotels))
	from := geodist.Coord{Lat: coordinates.Latitude, Lon: coordinates.Longitude}
	for i := range hotels {

		to := geodist.Coord{
			Lat: hotels[i].Coordinates.Lat,
			Lon: hotels[i].Coordinates.Lon,
		}

		_, km := geodist.HaversineDistance(from, to)
		hotelDists[i] = hotelDist{hotel: &hotels[i], distance: km}
	}

	quickSortHotelDist(hotelDists, 0, len(hotelDists))

	for i := 0; i < count; i++ {
		result[i] = *hotelDists[i].hotel
	}
	return result, nil
}

func RecommendHotelsRate(city string, count int) ([]Hotel, error) {
	hotels, err := GetHotelsInCity(city)
	if err != nil {
		return nil, err
	}

	if len(hotels) <= count {
		return hotels, err
	}

	quickSortRate(hotels, 0, len(hotels))

	result := make([]Hotel, count)
	for i := 0; i < count; i++ {
		result[i] = hotels[len(hotels)-1-i]
	}
	return result, nil
}

type hotelDist struct {
	hotel    *Hotel
	distance float64
}

func quickSortHotelDist(arr []hotelDist, from int, to int) {

	if from < 0 || from >= to {
		return
	}
	pivot := partitionHotelDist(arr, from, to)

	quickSortHotelDist(arr, from, pivot)
	quickSortHotelDist(arr, pivot+1, to)
}

func partitionHotelDist(arr []hotelDist, from int, to int) int {

	pivot := arr[from]
	pivotPos := from

	for j := from + 1; j < to; j++ {
		if arr[j].distance < pivot.distance {
			pivotPos++

			temp := arr[pivotPos]
			arr[pivotPos] = arr[j]
			arr[j] = temp
		}

	}

	arr[from] = arr[pivotPos]
	arr[pivotPos] = pivot
	return pivotPos
}

func quickSortRate(arr []Hotel, from int, to int) {

	if from < 0 || from >= to {
		return
	}
	pivot := partitionRate(arr, from, to)

	quickSortRate(arr, from, pivot)
	quickSortRate(arr, pivot+1, to)
}

func partitionRate(arr []Hotel, from int, to int) int {

	pivot := arr[from]
	pivotPos := from

	for j := from + 1; j < to; j++ {
		if arr[j].Rate < pivot.Rate {
			pivotPos++

			temp := arr[pivotPos]
			arr[pivotPos] = arr[j]
			arr[j] = temp
		}

	}

	arr[from] = arr[pivotPos]
	arr[pivotPos] = pivot
	return pivotPos
}
