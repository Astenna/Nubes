package faas_lib_test

import (
	"testing"

	"github.com/Astenna/Nubes/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// LOAD
// - success if exists
// - fail if not
// - changes made through get/set are saved (retrieve from another instance)
// - changes made through state changing method exist

func TestLoadFailsIfIdDoesNotExist(t *testing.T) {
	// Arrange
	// Act
	result, err := lib.Load[types.User](uuid.New().String())
	// Assert
	if err != nil || result == nil {
		t.Error("export should failed, but found no error")
	}
}

func TestLoadReturnsInitializedInstance(t *testing.T) {
	// Arrange
	// Act
	result, err := lib.Load[types.User](existingUserId)
	// Assert
	require.Equal(t, nil, err, "error should be null")
	require.Equal(t, existingUserId, result.Email, "IDs should match")
}

func TestLoadGetSetShouldSaveChanges(t *testing.T) {
	// Arrange
	newLastName := uuid.New().String()
	// Act
	instance1, err1 := lib.Load[types.User](existingUserId)
	require.Equal(t, nil, err1, "error should be null", err1)
	instance2, err2 := lib.Load[types.User](existingUserId)
	require.Equal(t, nil, err2, "error should be null", err2)
	err3 := instance2.SetLastName(newLastName)
	require.Equal(t, nil, err3, "error should be null", err3)
	instance1LastName, err4 := instance1.GetLastName()
	require.Equal(t, nil, err4, "error should be null", err4)
	// Assert
	require.Equal(t, instance1, instance2, "IDs should match")
	require.Equal(t, instance1, instance2, "IDs should match")
	require.Equal(t, newLastName, instance1LastName, "changes made to the objects with the same IDs should be visible from different instances")
}

func TestLoadStateChangingMethodsShouldSaveChanges(t *testing.T) {
	// Arrange
	shop, _ := lib.Export[types.Shop](types.Shop{Name: "My first shop"})
	existingShopId := shop.Id
	initialQuantityAvailable := 100
	decreaseBy := 5
	product := types.Product{
		Name:              "Product1",
		QuantityAvailable: initialQuantityAvailable,
		Price:             88.88,
		SoldBy:            *lib.NewReference[types.Shop](existingShopId),
	}
	exportedProduct, exportError := lib.Export[types.Product](product)
	require.Equal(t, nil, exportError, "error occurred while exporting the product in arrange step", exportError)

	// Act
	loadedProduct, loadExistingProductError := lib.Load[types.Product](exportedProduct.Id)
	require.Equal(t, nil, loadExistingProductError, "error occurred while loading existing product", loadExistingProductError)
	methodInovcationError := loadedProduct.DecreaseAvailabilityBy(decreaseBy)
	require.Equal(t, nil, methodInovcationError, "error occurred while invoking method on product instances", methodInovcationError)
	modifiedQuantity, quantityRetrievalError := exportedProduct.GetQuantityAvailable()
	require.Equal(t, nil, quantityRetrievalError, "error occured while exucting GetQuantityAvailable", quantityRetrievalError)

	// Assert
	require.Equal(t, modifiedQuantity, initialQuantityAvailable-decreaseBy, "QuantityAvailable was not modified")
}
