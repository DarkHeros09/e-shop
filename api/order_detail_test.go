package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

func TestCreateOrderDetailAPI(t *testing.T) {
	user, _ := randomODUser(t)
	orderDetail := createRandomOrderDetail(t, user)

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
				"user_id":    orderDetail.UserID,
				"total":      orderDetail.Total,
				"payment_id": orderDetail.PaymentID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.CreateOrderDetailParams{
					UserID:    orderDetail.UserID,
					Total:     orderDetail.Total,
					PaymentID: orderDetail.PaymentID,
				}

				store.EXPECT().
					CreateOrderDetail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(orderDetail, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrderDetail(t, recorder.Body, orderDetail)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"user_id":    orderDetail.UserID,
				"total":      orderDetail.Total,
				"payment_id": orderDetail.PaymentID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateOrderDetail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"user_id":    orderDetail.UserID,
				"total":      orderDetail.Total,
				"payment_id": orderDetail.PaymentID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateOrderDetailParams{
					UserID:    orderDetail.UserID,
					Total:     orderDetail.Total,
					PaymentID: orderDetail.PaymentID,
				}

				store.EXPECT().
					CreateOrderDetail(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.OrderDetail{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUserID",
			body: gin.H{
				"user_id":    0,
				"total":      orderDetail.Total,
				"payment_id": orderDetail.PaymentID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 0, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateOrderDetail(gomock.Any(), gomock.Any()).
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

			url := "/orderdetails"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomODUser(t *testing.T) (user db.User, password string) {
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

func createRandomOrderDetail(t *testing.T, user db.User) (orderDetail db.OrderDetail) {
	orderDetail = db.OrderDetail{
		ID:        util.RandomMoney(),
		UserID:    user.ID,
		Total:     util.RandomDecimal(1, 100),
		PaymentID: util.RandomMoney(),
	}
	return
}

func requireBodyMatchOrderDetail(t *testing.T, body *bytes.Buffer, orderDetail db.OrderDetail) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotOrderDetail db.OrderDetail
	err = json.Unmarshal(data, &gotOrderDetail)

	require.NoError(t, err)
	require.Equal(t, orderDetail.ID, gotOrderDetail.ID)
	require.Equal(t, orderDetail.Total, gotOrderDetail.Total)
	require.Equal(t, orderDetail.PaymentID, gotOrderDetail.PaymentID)
}
