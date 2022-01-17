package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/DarkHeros09/e-shop/v2/db/mock"
	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetProductAPI(t *testing.T) {
	product := randomProduct()

	testCases := []struct {
		name          string
		productID     int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			productID: product.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(product, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatchProduct(t, recorder.Body, product)
			},
		},
		{
			name:      "NotFound",
			productID: product.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(db.Product{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			productID: product.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Eq(product.ID)).
					Times(1).
					Return(db.Product{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			productID: 0,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetProduct(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) // no need to call defer ctrl.finish() in 1.6V

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStub(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/products/%d", tc.productID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

// func TestUpdateProductAPI(t *testing.T) {
// 	product := randomProduct()
// 	arg := updateProductRequest{
// 		ID:         product.ID,
// 		Price:      product.Price,
// 		Active:     true,
// 		DiscountID: product.DiscountID,
// 	}

// 	send, _ := json.Marshal(arg)
// 	dd := string(send)

// 	product2 := db.Product{
// 		ID:          product.ID,
// 		Name:        product.Name,
// 		Description: product.Description,
// 		Sku:         product.Sku,
// 		CategoryID:  product.CategoryID,
// 		InventoryID: product.InventoryID,
// 		Price:       product.Price,
// 		Active:      arg.Active,
// 		DiscountID:  product.DiscountID,
// 	}

// 	testCases := []struct {
// 		name          string
// 		arg           int64
// 		buildStub     func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			arg:  arg.ID,
// 			buildStub: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					UpdateProduct(gomock.Any(), gomock.Eq(string(dd))).
// 					Times(1).
// 					Return(product2, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 			},
// 		},
// {
// 	name: "NotFound",
// 	arg:  arg,
// 	buildStub: func(store *mockdb.MockStore) {
// 		store.EXPECT().
// 			UpdateProduct(gomock.Any(), gomock.Eq(arg)).
// 			Times(1).
// 			Return(db.Product{}, sql.ErrNoRows)
// 	},
// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		require.Equal(t, http.StatusNotFound, recorder.Code)
// 	},
// },
// {
// 	name: "InternalError",
// 	arg:  arg,
// 	buildStub: func(store *mockdb.MockStore) {
// 		store.EXPECT().
// 			UpdateProduct(gomock.Any(), gomock.Eq(arg)).
// 			Times(1).
// 			Return(db.Product{}, sql.ErrConnDone)
// 	},
// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 	},
// },
// {
// 	name: "InvalidID",
// 	arg:  arg,
// 	buildStub: func(store *mockdb.MockStore) {
// 		store.EXPECT().
// 			UpdateProduct(gomock.Any(), gomock.Any()).
// 			Times(0)

// 	},
// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
// 	},
// },
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t) // no need to call defer ctrl.finish() in 1.6V

// 			store := mockdb.NewMockStore(ctrl)

// 			// build stubs
// 			tc.buildStub(store)

// 			// start test server and send request
// 			server := NewServer(store)
// 			recorder := httptest.NewRecorder()

// 			url := fmt.Sprintf("/products/%d", tc.arg)
// 			request, err := http.NewRequest(http.MethodPut, url, nil)
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)

// 			//check response
// 			tc.checkResponse(t, recorder)
// 		})

// 	}
// }

func randomProduct() db.Product {
	return db.Product{
		ID:          util.RandomInt(1, 1000),
		Name:        util.RandomUser(),
		Description: util.RandomUser(),
		Sku:         util.RandomUser(),
		CategoryID:  util.RandomInt(1, 500),
		InventoryID: util.RandomInt(1, 500),
		Price:       fmt.Sprint(util.RandomMoney()),
		Active:      util.RandomBool(),
		DiscountID:  util.RandomInt(1, 500),
	}
}

func requireBodyMatchProduct(t *testing.T, body *bytes.Buffer, product db.Product) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotProduct db.Product
	err = json.Unmarshal(data, &gotProduct)
	require.NoError(t, err)
	require.Equal(t, product, gotProduct)
}
