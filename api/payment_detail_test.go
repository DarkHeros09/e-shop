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

func TestGetPaymentDetailAPI(t *testing.T) {
	user, _ := randomPDUser(t)
	orderDetail := createRandomOrderDetail(t, user)
	paymentDetail := createRandomPaymentDetail(t, orderDetail)

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
			ID:   paymentDetail.ID,
			body: gin.H{
				"id": paymentDetail.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.GetPaymentDetailParams{
					ID:     paymentDetail.ID,
					UserID: user.ID,
				}
				store.EXPECT().
					GetPaymentDetail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(paymentDetail, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPaymentDetail(t, recorder.Body, paymentDetail)
			},
		},
		{
			name: "NoAuthorization",
			ID:   paymentDetail.ID,
			body: gin.H{
				"id": paymentDetail.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPaymentDetail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			ID:   paymentDetail.ID,
			body: gin.H{
				"id": paymentDetail.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, "unauthorizedUser", time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPaymentDetail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.PaymentDetail{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "NotFound",
			ID:   paymentDetail.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {

				arg := db.GetPaymentDetailParams{
					ID:     paymentDetail.ID,
					UserID: user.ID,
				}
				store.EXPECT().
					GetPaymentDetail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.PaymentDetail{}, sql.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			ID:   paymentDetail.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStub: func(store *mockdb.MockStore) {
				arg := db.GetPaymentDetailParams{
					ID:     paymentDetail.ID,
					UserID: user.ID,
				}
				store.EXPECT().
					GetPaymentDetail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.PaymentDetail{}, sql.ErrConnDone)

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
					GetPaymentDetail(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/payment-details/%d", tc.ID)
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			//check response
			tc.checkResponse(t, recorder)
		})

	}

}

func TestListPaymentDetailsAPI(t *testing.T) {
	n := 5
	paymentDetails := []db.ListPaymentDetailsRow{}
	user, _ := randomPDUser(t)
	orderDetail1 := createRandomOrderDetail(t, user)
	paymentDetail1 := createRandomPaymentDetailForList(t, orderDetail1)
	orderDetail2 := createRandomOrderDetail(t, user)
	paymentDetail2 := createRandomPaymentDetailForList(t, orderDetail2)
	orderDetail3 := createRandomOrderDetail(t, user)
	paymentDetail3 := createRandomPaymentDetailForList(t, orderDetail3)

	paymentDetails = append(paymentDetails, paymentDetail1, paymentDetail2, paymentDetail3)

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
				arg := db.ListPaymentDetailsParams{
					UserID: user.ID,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListPaymentDetails(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(paymentDetails, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPaymentDetails(t, recorder.Body, paymentDetails)
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
					ListPaymentDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListPaymentDetailsRow{}, sql.ErrConnDone)
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
					ListPaymentDetails(gomock.Any(), gomock.Any()).
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
					ListPaymentDetails(gomock.Any(), gomock.Any()).
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

			url := "/payment-details"
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

func TestUpdatePaymentDetailAPI(t *testing.T) {
	user, _ := randomPDUser(t)
	orderDetail := createRandomOrderDetail(t, user)
	paymentDetail := createRandomPaymentDetail(t, orderDetail)

	testCases := []struct {
		name            string
		PaymentDetailID int64
		body            gin.H
		setupAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:            "OK",
			PaymentDetailID: paymentDetail.ID,
			body: gin.H{
				"id":       paymentDetail.ID,
				"order_id": paymentDetail.OrderID,
				"amount":   99,
				"provider": "done",
				"status":   "done",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentDetailParams{
					ID:       paymentDetail.ID,
					UserID:   orderDetail.UserID,
					OrderID:  paymentDetail.OrderID,
					Amount:   99,
					Provider: "done",
					Status:   "done",
				}
				store.EXPECT().
					UpdatePaymentDetail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(paymentDetail, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:            "Unauthorized",
			PaymentDetailID: paymentDetail.ID,
			body: gin.H{
				"id":       paymentDetail.ID,
				"order_id": paymentDetail.OrderID,
				"amount":   99,
				"provider": "done",
				"status":   "done",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, "unauthorized", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdatePaymentDetail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.PaymentDetail{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:            "NoAuthorization",
			PaymentDetailID: paymentDetail.ID,
			body: gin.H{
				"id":       paymentDetail.ID,
				"order_id": paymentDetail.OrderID,
				"amount":   99,
				"provider": "done",
				"status":   "done",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdatePaymentDetail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:            "InternalError",
			PaymentDetailID: paymentDetail.ID,
			body: gin.H{
				"id":       paymentDetail.ID,
				"order_id": paymentDetail.OrderID,
				"amount":   99,
				"provider": "done",
				"status":   "done",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdatePaymentDetailParams{
					ID:       paymentDetail.ID,
					UserID:   orderDetail.UserID,
					OrderID:  paymentDetail.OrderID,
					Amount:   99,
					Provider: "done",
					Status:   "done",
				}
				store.EXPECT().
					UpdatePaymentDetail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.PaymentDetail{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:            "InvalidID",
			PaymentDetailID: 0,
			body: gin.H{
				"id":       0,
				"order_id": paymentDetail.OrderID,
				"amount":   99,
				"provider": "done",
				"status":   "done",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdatePaymentDetail(gomock.Any(), gomock.Any()).
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
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/payment-details/%d", tc.PaymentDetailID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomPDUser(t *testing.T) (user db.User, password string) {
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

func createRandomPaymentDetail(t *testing.T, orderDetail db.OrderDetail) (paymentDetail db.PaymentDetail) {
	paymentDetail = db.PaymentDetail{
		ID:       util.RandomMoney(),
		OrderID:  orderDetail.ID,
		Amount:   int32(util.RandomMoney()),
		Provider: util.RandomUser(),
		Status:   util.RandomUser(),
	}
	return
}

func createRandomPaymentDetailForList(t *testing.T, orderDetail db.OrderDetail) (paymentDetails db.ListPaymentDetailsRow) {
	paymentDetails = db.ListPaymentDetailsRow{
		ID:       util.RandomMoney(),
		OrderID:  orderDetail.ID,
		Amount:   int32(util.RandomMoney()),
		Provider: util.RandomUser(),
		Status:   util.RandomUser(),
	}
	return
}

func requireBodyMatchPaymentDetail(t *testing.T, body *bytes.Buffer, paymentDetail db.PaymentDetail) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotPaymentDetail db.PaymentDetail
	err = json.Unmarshal(data, &gotPaymentDetail)

	require.NoError(t, err)
	require.Equal(t, paymentDetail.ID, gotPaymentDetail.ID)
	require.Equal(t, paymentDetail.OrderID, gotPaymentDetail.OrderID)
	require.Equal(t, paymentDetail.Amount, gotPaymentDetail.Amount)
	require.Equal(t, paymentDetail.Status, gotPaymentDetail.Status)
	require.Equal(t, paymentDetail.Provider, gotPaymentDetail.Provider)
}

func requireBodyMatchPaymentDetails(t *testing.T, body *bytes.Buffer, paymentDetails []db.ListPaymentDetailsRow) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotPaymentDetails []db.ListPaymentDetailsRow
	err = json.Unmarshal(data, &gotPaymentDetails)
	require.NoError(t, err)
	require.Equal(t, paymentDetails, gotPaymentDetails)
}
