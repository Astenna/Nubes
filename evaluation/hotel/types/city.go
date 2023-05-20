package types

import (
	"github.com/Astenna/Nubes/lib"
	"github.com/jftuga/geodist"
)

type City struct {
	CityName        string `nubes:"Id" dynamodbav:"Id"`
	Region          string
	Description     string
	Hotels          lib.ReferenceNavigationList[Hotel] `nubes:"hasOne-City" dynamodbav:"-"`
	isInitialized   bool
	invocationDepth int
}

func (o City) GetTypeName() string {
	return "City"
}

type CloseToParams struct {
	Longitude float64
	Latitude  float64
	Count     int
}

func (c City) GetHotelsCloseTo(param CloseToParams) ([]Hotel, error) {
	c.invocationDepth++
	if c.isInitialized && c.invocationDepth == 1 {
		_libError := lib.GetStub(c.CityName, &c)
		if _libError != nil {
			c.invocationDepth--
			return *new([]Hotel), _libError
		}
	}
	hotels, err := c.Hotels.GetStubs()
	result := make([]Hotel, param.Count)

	if err != nil {
		c.invocationDepth--
		return nil, err
	}

	if len(hotels) <= param.Count {
		c.invocationDepth--
		return hotels, err
	}
	hotelDists := make([]hotelDist, len(hotels))
	from := geodist.Coord{Lat: param.Latitude, Lon: param.Longitude}
	for i := range hotels {

		to := geodist.Coord{
			Lat: hotels[i].Coordinates.Lat,
			Lon: hotels[i].Coordinates.Lon,
		}

		_, km := geodist.HaversineDistance(from, to)
		hotelDists[i] = hotelDist{hotel: &hotels[i], distance: km}
	}

	quickSortHotelDist(hotelDists, 0, len(hotelDists))

	for i := 0; i < param.Count; i++ {
		result[i] = *hotelDists[i].hotel
	}
	c.invocationDepth--
	return result, nil
}

func (c City) GetHotelsWithBestRates(count int) ([]Hotel, error) {
	c.invocationDepth++
	if c.isInitialized && c.invocationDepth == 1 {
		_libError := lib.GetStub(c.CityName, &c)
		if _libError != nil {
			c.invocationDepth--
			return *new([]Hotel), _libError
		}
	}
	hotels, err := c.Hotels.GetStubs()
	result := make([]Hotel, count)

	if err != nil {
		c.invocationDepth--
		return nil, err
	}

	if len(hotels) <= count {
		c.invocationDepth--
		return result, err
	}

	quickSortRate(hotels, 0, len(hotels))

	for i := 0; i < count; i++ {
		result[i] = hotels[len(hotels)-i-1]
	}
	c.invocationDepth--

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
func (receiver City) GetId() string {
	return receiver.CityName
}
func (receiver *City) Init() {
	receiver.isInitialized = true
	receiver.Hotels = *lib.NewReferenceNavigationList[Hotel](lib.ReferenceNavigationListParam{OwnerId: receiver.CityName, OwnerTypeName: receiver.GetTypeName(), OtherTypeName: (*new(Hotel)).GetTypeName(), ReferringFieldName: "City", IsManyToMany: false})
}
func (receiver *City) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.CityName)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}
