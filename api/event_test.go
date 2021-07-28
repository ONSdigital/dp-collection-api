package api_test

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-collection-api/api"
	"github.com/ONSdigital/dp-collection-api/api/mock"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetEvents(t *testing.T) {

	Convey("Given a request to GET collection events", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("GET", "http://localhost:26000/collections/123/events", nil)
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.Router.ServeHTTP(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
			})

			Convey("Then the collection store is called to get collection data", func() {
				So(len(collectionStore.GetCollectionEventsCalls()), ShouldEqual, 1)
				getCollectionsCall := collectionStore.GetCollectionEventsCalls()[0]
				So(getCollectionsCall.QueryParams.Limit, ShouldEqual, limit)
				So(getCollectionsCall.QueryParams.Offset, ShouldEqual, offset)
				So(getCollectionsCall.QueryParams.CollectionID, ShouldEqual, "123")
			})

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("Then the response body should contain the collection events", func() {
				body, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				response := models.EventsResponse{}
				err = json.Unmarshal(body, &response)
				So(err, ShouldBeNil)
				So(response.TotalCount, ShouldEqual, totalCount)
				So(response.Count, ShouldEqual, len(response.Items))
				So(response.Offset, ShouldEqual, offset)
				So(response.Limit, ShouldEqual, limit)
				So(response.TotalCount, ShouldEqual, totalCount)
				So(response.Items[0].Type, ShouldEqual, "CREATED")
				So(response.Items[0].Email, ShouldEqual, "test@test.com")
			})
		})
	})
}

func TestGetEvents_paginationError(t *testing.T) {

	Convey("Given a paginator that returns an error", t, func() {

		paginator := &mock.PaginatorMock{
			ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
				return 1, 0, pagination.ErrInvalidLimitParameter
			},
		}
		collectionStore := &mock.CollectionStoreMock{}

		Convey("When the request is made to GET collection events", func() {

			r := httptest.NewRequest("GET", "http://localhost:26000/collections/123/events", nil)
			w := httptest.NewRecorder()

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.Router.ServeHTTP(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
			})

			Convey("Then the expected error code is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

func TestGetEvents_collectionStoreError(t *testing.T) {

	Convey("Given a collection store that returns an error", t, func() {

		paginator := mockPaginator()
		collectionStore := &mock.CollectionStoreMock{
			GetCollectionByIDFunc: func(ctx context.Context, id string, eTagSelector string) (*models.Collection, error) {
				return nil, nil
			},
			GetCollectionEventsFunc: func(ctx context.Context, queryParams collections.EventsQueryParams) ([]models.Event, int, error) {
				return nil, 0, errors.New("store error")
			},
		}

		Convey("When the request is made to GET collection events", func() {

			r := httptest.NewRequest("GET", "http://localhost:26000/collections/123/events", nil)
			w := httptest.NewRecorder()

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.Router.ServeHTTP(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
			})

			Convey("Then the expected error code is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestGetEvents_internalError(t *testing.T) {

	Convey("Given a paginator that returns an unrecognised error", t, func() {

		paginator := &mock.PaginatorMock{
			ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
				return 1, 0, errors.New("unrecognised error")
			},
		}
		collectionStore := &mock.CollectionStoreMock{}

		Convey("When the request is made to GET collection events", func() {

			r := httptest.NewRequest("GET", "http://localhost:26000/collections/123/events", nil)
			w := httptest.NewRecorder()

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.Router.ServeHTTP(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
			})

			Convey("Then an internal server error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
