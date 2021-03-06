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
	"time"

	mockdb "github.com/DarkHeros09/e-shop/v2/db/mock"
	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/token"
	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateOrderItemAPI(t *testing.T) {
	user, _ := randomOIUser(t)
	orderDetail := createRandomOrderDetail(t, user)
	orderItem := createRandomOrderItem(t, orderDetail)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"order_id":   orderItem.OrderID,
				"product_id": orderItem.ProductID,
				"quantity":   orderItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetOrderDetail(gomock.Any(), gomock.Eq(orderItem.OrderID)).
					Times(1).
					Return(orderDetail, nil)

				arg := db.CreateOrderItemParams{
					OrderID:   orderItem.OrderID,
					ProductID: orderItem.ProductID,
					Quantity:  orderItem.Quantity,
				}

				store.EXPECT().
					CreateOrderItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(orderItem, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrderItem(t, recorder.Body, orderItem)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"order_id":   orderItem.OrderID,
				"product_id": orderItem.ProductID,
				"quantity":   orderItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrderDetail(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateOrderItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"order_id":   orderItem.OrderID,
				"product_id": orderItem.ProductID,
				"quantity":   orderItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetOrderDetail(gomock.Any(), gomock.Eq(orderItem.OrderID)).
					Times(1).
					Return(orderDetail, nil)

				arg := db.CreateOrderItemParams{
					OrderID:   orderItem.OrderID,
					ProductID: orderItem.ProductID,
					Quantity:  orderItem.Quantity,
				}

				store.EXPECT().
					CreateOrderItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.OrderItem{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidOrderID",
			body: gin.H{
				"order_id":   0,
				"product_id": orderItem.ProductID,
				"quantity":   orderItem.Quantity,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrderDetail(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateOrderItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/order-items"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetOrderItemAPI(t *testing.T) {
	user, _ := randomOIUser(t)
	orderDetail := createRandomOrderDetail(t, user)
	orderItem := createRandomOrderItem(t, orderDetail)

	testCases := []struct {
		name          string
		ID            int64
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   orderItem.ID,
			body: gin.H{
				"id": orderItem.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.GetOrderItemParams{
					ID:     orderItem.ID,
					UserID: user.ID,
				}
				store.EXPECT().
					GetOrderItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(orderItem, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrderItem(t, recorder.Body, orderItem)
			},
		},
		{
			name: "NoAuthorization",
			ID:   orderItem.ID,
			body: gin.H{
				"id": orderItem.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrderItem(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   orderItem.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.GetOrderItemParams{
					ID:     orderItem.ID,
					UserID: user.ID,
				}
				store.EXPECT().
					GetOrderItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.OrderItem{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   orderItem.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.GetOrderItemParams{
					ID:     orderItem.ID,
					UserID: user.ID,
				}
				store.EXPECT().
					GetOrderItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.OrderItem{}, sql.ErrConnDone)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			ID:   0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetOrderItem(gomock.Any(), gomock.Any()).
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/order-items/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

func TestListOrderItemAPI(t *testing.T) {
	n := 5
	orderItems := make([]db.OrderItem, n)
	user, _ := randomOIUser(t)
	orderDetail1 := createRandomOrderDetail(t, user)
	orderDetail2 := createRandomOrderDetail(t, user)
	orderDetail3 := createRandomOrderDetail(t, user)
	orderItem1 := createRandomOrderItem(t, orderDetail1)
	orderItem2 := createRandomOrderItem(t, orderDetail2)
	orderItem3 := createRandomOrderItem(t, orderDetail3)

	orderItems = append(orderItems, orderItem1, orderItem2, orderItem3)

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListOrderItemsParams{
					UserID: user.ID,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListOrderItems(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(orderItems, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrderItems(t, recorder.Body, orderItems)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListOrderItems(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.OrderItem{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListOrderItems(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: 100000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListOrderItems(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/order-items"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomOIUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:        util.RandomMoney(),
		Username:  util.RandomUser(),
		Password:  hashedPassword,
		Telephone: int32(util.RandomInt(7, 7)),
		Email:     util.RandomEmail(),
	}
	return
}

func createRandomOrderItem(t *testing.T, orderDetail db.OrderDetail) (orderItem db.OrderItem) {
	orderItem = db.OrderItem{
		ID:        util.RandomMoney(),
		OrderID:   orderDetail.ID,
		ProductID: util.RandomMoney(),
		Quantity:  int32(util.RandomMoney()),
	}
	return
}

func requireBodyMatchOrderItem(t *testing.T, body *bytes.Buffer, orderItem db.OrderItem) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotOrderItem db.OrderItem
	err = json.Unmarshal(data, &gotOrderItem)

	require.NoError(t, err)
	require.Equal(t, orderItem.ID, gotOrderItem.ID)
	require.Equal(t, orderItem.ProductID, gotOrderItem.ProductID)
	require.Equal(t, orderItem.OrderID, gotOrderItem.OrderID)
	require.Equal(t, orderItem.Quantity, gotOrderItem.Quantity)
}

func requireBodyMatchOrderItems(t *testing.T, body *bytes.Buffer, orderItems []db.OrderItem) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotOrderItems []db.OrderItem
	err = json.Unmarshal(data, &gotOrderItems)
	require.NoError(t, err)
	require.Equal(t, orderItems, gotOrderItems)
}
