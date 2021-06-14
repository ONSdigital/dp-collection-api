package api_test

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-collection-api/api/mock"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-collection-api/api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetCollections(t *testing.T) {

	offset := 0
	limit := 1
	totalCount := 3

	Convey("Given a request to GET collections", t, func() {

		paginator := &mock.PaginatorMock{
			ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
				return offset, limit, nil
			},
		}

		collectionStore := &mock.CollectionStoreMock{
			GetCollectionsFunc: func(ctx context.Context, offset int, limit int) ([]models.Collection, int, error) {
				return []models.Collection{{
					ID:          "123",
					Name:        "collection 1",
					PublishDate: time.Time{},
					LastUpdated: time.Time{},
				}}, totalCount, nil
			},
		}

		r := httptest.NewRequest("GET", "http://localhost:26000/collections", nil)
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
				So(paginator.ReadPaginationParametersCalls()[0].R, ShouldEqual, r)
			})

			Convey("Then the collection store is called to get collection data", func() {
				So(len(collectionStore.GetCollectionsCalls()), ShouldEqual, 1)
				So(collectionStore.GetCollectionsCalls()[0].Limit, ShouldEqual, limit)
				So(collectionStore.GetCollectionsCalls()[0].Offset, ShouldEqual, offset)
			})

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("Then the response body should contain the collections", func() {
				body, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				response := models.CollectionsResponse{}
				err = json.Unmarshal(body, &response)
				So(err, ShouldBeNil)
				So(response.TotalCount, ShouldEqual, totalCount)
				So(response.Count, ShouldEqual, len(response.Items))
				So(response.Offset, ShouldEqual, offset)
				So(response.Limit, ShouldEqual, limit)
				So(response.TotalCount, ShouldEqual, totalCount)
			})
		})
	})
}

func TestGetCollections_paginationError(t *testing.T) {

	Convey("Given a paginator that returns an error", t, func() {

		paginator := &mock.PaginatorMock{
			ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
				return 1, 0, pagination.ErrInvalidLimitParameter
			},
		}

		collectionStore := &mock.CollectionStoreMock{}

		Convey("When the request is made to GET collections", func() {

			r := httptest.NewRequest("GET", "http://localhost:26000/collections", nil)
			w := httptest.NewRecorder()

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
				So(paginator.ReadPaginationParametersCalls()[0].R, ShouldEqual, r)
			})

			Convey("Then the expected error code is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

func TestGetCollections_collectionStoreError(t *testing.T) {

	Convey("Given a collection store that returns an error", t, func() {

		paginator := &mock.PaginatorMock{
			ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
				return 0, 0, nil
			},
		}

		collectionStore := &mock.CollectionStoreMock{
			GetCollectionsFunc: func(ctx context.Context, offset int, limit int) ([]models.Collection, int, error) {
				return nil, 0, errors.New("store error")
			},
		}

		Convey("When the request is made to GET collections", func() {

			r := httptest.NewRequest("GET", "http://localhost:26000/collections", nil)
			w := httptest.NewRecorder()

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
				So(paginator.ReadPaginationParametersCalls()[0].R, ShouldEqual, r)
			})

			Convey("Then the expected error code is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestGetCollections_internalError(t *testing.T) {

	Convey("Given a paginator that returns an unrecognised error", t, func() {

		paginator := &mock.PaginatorMock{
			ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
				return 1, 0, errors.New("unrecognised error")
			},
		}

		collectionStore := &mock.CollectionStoreMock{}

		Convey("When the request is made to GET collections", func() {

			r := httptest.NewRequest("GET", "http://localhost:26000/collections", nil)
			w := httptest.NewRecorder()

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the paginator is called to extract pagination parameters", func() {
				So(len(paginator.ReadPaginationParametersCalls()), ShouldEqual, 1)
				So(paginator.ReadPaginationParametersCalls()[0].R, ShouldEqual, r)
			})

			Convey("Then an internal server error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
