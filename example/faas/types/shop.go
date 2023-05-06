package types

import (
	"fmt"
	"math"

	"github.com/Astenna/Nubes/lib"
	"github.com/jftuga/geodist"
)

type Shop struct {
	Id              string
	Name            string
	Owners          lib.ReferenceNavigationList[User]    `nubes:"hasMany-Shops" dynamodbav:"-"`
	Products        lib.ReferenceNavigationList[Product] `nubes:"hasOne-SoldBy,readonly" dynamodbav:"-"`
	isInitialized   bool
	invocationDepth int
}

func (Shop) GetTypeName() string {
	return "Shop"
}

func (s Shop) GetOwners() ([]string, error) {
	if !s.isInitialized {
		return nil, fmt.Errorf(`fields of type ReferenceNavigationList can be used only after instance initialization. 
			Use lib.Load or lib.Export from the Nubes library to create initialized instances`)
	}
	return s.Owners.GetIds()
}

// Example of a method returning a Nobject
func (s Shop) GetNearestOwnerCopy(point Coordinates) (User, error) {
	s.invocationDepth++
	if s.isInitialized && s.invocationDepth == 1 {
		_libError := lib.GetStub(s.Id, &s)
		if _libError != nil {
			s.invocationDepth--
			return *new(User), _libError
		}
	}

	owners, err := s.Owners.GetStubs()
	if err != nil {
		s.invocationDepth--
		return *new(User), err
	}

	var closestOwner User
	from := geodist.Coord{Lat: point.Latitude, Lon: point.Longitude}
	min := math.MaxFloat32
	for _, owner := range owners {

		to := geodist.Coord{
			Lat: owner.AddressCoordinates.Latitude,
			Lon: owner.AddressCoordinates.Longitude,
		}

		_, km := geodist.HaversineDistance(from, to)

		if km < min {
			min = km
			closestOwner = owner
		}
	}
	s.invocationDepth--

	return closestOwner, nil
}

// Example of a method returning a Nobject's reference
func (s Shop) GetNearestOwnerReference(point Coordinates) (lib.Reference[User], error) {
	s.invocationDepth++
	if s.isInitialized && s.invocationDepth == 1 {
		_libError := lib.GetStub(s.Id, &s)
		if _libError != nil {
			s.invocationDepth--
			return *new(lib.Reference[User]), _libError
		}
	}
	owners, err := s.Owners.GetStubs()
	if err != nil {
		s.invocationDepth--
		return *new(lib.Reference[User]), err
	}

	var closestOwner User
	from := geodist.Coord{Lat: point.Latitude, Lon: point.Longitude}
	min := math.MaxFloat32
	for _, owner := range owners {

		to := geodist.Coord{
			Lat: owner.AddressCoordinates.Latitude,
			Lon: owner.AddressCoordinates.Longitude,
		}

		_, km := geodist.HaversineDistance(from, to)

		if km < min {
			min = km
			closestOwner = owner
		}
	}
	s.invocationDepth--

	return *lib.NewReference[User](closestOwner.Email), nil
}

// DeleteShop is an example of custom delete implementation
// Note that, the invocation of lib.Delete must be added.
func DeleteShop(id string) error {
	shopToBeDeleted, err := lib.Load[Shop](id)
	if err != nil {
		return err
	}

	shopProducts, err := shopToBeDeleted.Products.GetIds()
	if err != nil {
		return err
	}

	for _, id := range shopProducts {
		err = lib.Delete[Product](id)
		if err != nil {
			return err
		}
	}

	return lib.Delete[Shop](id)
}

func (receiver *Shop) Init() {
	receiver.isInitialized = true
	receiver.Products = *lib.NewReferenceNavigationList[Product](receiver.Id, receiver.GetTypeName(), "SoldBy", false)
	receiver.Owners = *lib.NewReferenceNavigationList[User](receiver.Id, receiver.GetTypeName(), "", true)
}
func (receiver *Shop) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}
