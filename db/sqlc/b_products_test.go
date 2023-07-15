package sqlc

import (
	"context"
	"simple_shop/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T) Product {
	username := createRandomBankAccount(t)

	arg := CreateProductParams{
		ProductName:   util.RandomString(8),
		KindOfProduct: util.RandomString(8),
		Owner:         username.Username,
		Currency:      "USD",
		Price:         int32(util.RandomIntNumber(50, 100)),
		Quantity:      int32(util.RandomIntNumber(20, 50)),
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, arg.ProductName, product.ProductName)
	require.Equal(t, arg.KindOfProduct, product.KindOfProduct)
	require.Equal(t, arg.Owner, product.Owner)
	require.Equal(t, arg.Currency, product.Currency)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.Quantity, product.Quantity)

	return product
}

func TestCreateProduct(t *testing.T) {
	createRandomProduct(t)
}
