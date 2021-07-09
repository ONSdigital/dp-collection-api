package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-collection-api/api/mock"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"github.com/ONSdigital/dp-collection-api/pagination"
	"github.com/ONSdigital/dp-mongodb/v2/pkg/mongodb"
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

var (
	offset     = 0
	limit      = 1
	totalCount = 3
)

func TestGetCollections(t *testing.T) {

	Convey("Given a request to GET collections", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

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
				getCollectionsCall := collectionStore.GetCollectionsCalls()[0]
				So(getCollectionsCall.QueryParams.Limit, ShouldEqual, limit)
				So(getCollectionsCall.QueryParams.Offset, ShouldEqual, offset)
				So(getCollectionsCall.QueryParams.OrderBy, ShouldEqual, collections.OrderByDefault)
				So(getCollectionsCall.QueryParams.NameSearch, ShouldEqual, "")
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

func TestGetCollections_orderByPublishDate(t *testing.T) {

	Convey("Given a request to GET collections", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("GET", "http://localhost:26000/collections?order_by=publish_date", nil)
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the collection store is called with the expected orderBy value", func() {
				So(len(collectionStore.GetCollectionsCalls()), ShouldEqual, 1)
				getCollectionsCall := collectionStore.GetCollectionsCalls()[0]
				So(getCollectionsCall.QueryParams.Limit, ShouldEqual, limit)
				So(getCollectionsCall.QueryParams.Offset, ShouldEqual, offset)
				So(getCollectionsCall.QueryParams.OrderBy, ShouldEqual, collections.OrderByPublishDate)
				So(getCollectionsCall.QueryParams.NameSearch, ShouldEqual, "")
			})
		})
	})
}

func TestGetCollections_nameSearch(t *testing.T) {

	Convey("Given a request to GET collections with a name search value", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("GET", "http://localhost:26000/collections?name=LMSV3", nil)
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the collection store is called with the expected orderBy value", func() {
				So(len(collectionStore.GetCollectionsCalls()), ShouldEqual, 1)
				getCollectionsCall := collectionStore.GetCollectionsCalls()[0]
				So(getCollectionsCall.QueryParams.Limit, ShouldEqual, limit)
				So(getCollectionsCall.QueryParams.Offset, ShouldEqual, offset)
				So(getCollectionsCall.QueryParams.OrderBy, ShouldEqual, collections.OrderByDefault)
				So(getCollectionsCall.QueryParams.NameSearch, ShouldEqual, "LMSV3")
			})
		})
	})
}

func TestGetCollections_nameSearchTooLong(t *testing.T) {

	Convey("Given a request to GET collections with an empty order_by value", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("GET", "http://localhost:26000/collections?name=1234567890123456789012345678901234567890123456789012345678901234567890", nil)
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the expected error code is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then the response body should contain the expected error response", func() {
				body, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				response := models.ErrorsResponse{}
				err = json.Unmarshal(body, &response)
				So(err, ShouldBeNil)
				So(len(response.Errors), ShouldEqual, 1)
				So(response.Errors[0].Message, ShouldEqual, collections.ErrNameSearchTooLong.Error())
			})
		})
	})
}

func TestGetCollections_emptyOrderBy(t *testing.T) {

	Convey("Given a request to GET collections with an empty order_by value", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("GET", "http://localhost:26000/collections?order_by=", nil)
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.GetCollectionsHandler(w, r)

			Convey("Then the collection store is called with the expected orderBy value", func() {
				So(len(collectionStore.GetCollectionsCalls()), ShouldEqual, 1)
				getCollectionsCall := collectionStore.GetCollectionsCalls()[0]
				So(getCollectionsCall.QueryParams.OrderBy, ShouldEqual, collections.OrderByDefault)
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

		paginator := mockPaginator()
		collectionStore := &mock.CollectionStoreMock{
			GetCollectionsFunc: func(ctx context.Context, queryParams collections.QueryParams) ([]models.Collection, int, error) {
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

func TestGetCollections_invalidOrderByError(t *testing.T) {

	Convey("Given an invalid orderBy value", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		Convey("When the request is made to GET collections", func() {

			r := httptest.NewRequest("GET", "http://localhost:26000/collections?order_by=fubar", nil)
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

			Convey("Then the response body should contain the expected error response", func() {
				body, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				response := models.ErrorsResponse{}
				err = json.Unmarshal(body, &response)
				So(err, ShouldBeNil)
				So(len(response.Errors), ShouldEqual, 1)
				So(response.Errors[0].Message, ShouldEqual, collections.ErrInvalidOrderBy.Error())
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

func TestPostCollection(t *testing.T) {

	newCollectionJson := `{
		"name": "Coronavirus key indicators",
		"publish_date": "2020-05-05T14:58:29.317Z"
	}`
	expectedName := "Coronavirus key indicators"
	expectedID := "12345"
	api.NewID = func() (string, error) {
		return expectedID, nil
	}

	Convey("Given a request to POST a collection", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("POST", "http://localhost:26000/collections", bytes.NewBufferString(newCollectionJson))
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.AddCollectionHandler(w, r)

			Convey("Then the collection store is called", func() {

				So(len(collectionStore.GetCollectionByNameCalls()), ShouldEqual, 1)
				So(collectionStore.GetCollectionByNameCalls()[0].Name, ShouldEqual, expectedName)

				So(len(collectionStore.UpsertCollectionCalls()), ShouldEqual, 1)
				getCollectionsCall := collectionStore.UpsertCollectionCalls()[0]
				So(getCollectionsCall.Collection.ID, ShouldEqual, expectedID)
				So(getCollectionsCall.Collection.Name, ShouldEqual, expectedName)
				So(getCollectionsCall.Collection.PublishDate.String(), ShouldEqual, "2020-05-05 14:58:29.317 +0000 UTC")
			})

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusCreated)
			})

			Convey("Then the response body should contain the collections", func() {
				body, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				response := models.Collection{}
				err = json.Unmarshal(body, &response)
				So(err, ShouldBeNil)
				So(response.Name, ShouldEqual, expectedName)
			})
		})
	})
}

func TestPostCollection_CollectionNameAlreadyExists(t *testing.T) {

	newCollectionJson := `{
		"name": "Coronavirus key indicators",
		"publish_date": "2020-05-05T14:58:29.317Z"
	}`
	expectedName := "Coronavirus key indicators"
	expectedID := "12345"
	api.NewID = func() (string, error) {
		return expectedID, nil
	}

	Convey("Given a request to POST a collection with a name that already exists", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()
		collectionStore.GetCollectionByNameFunc = func(ctx context.Context, name string) (*models.Collection, error) {
			return &models.Collection{}, nil
		}

		r := httptest.NewRequest("POST", "http://localhost:26000/collections", bytes.NewBufferString(newCollectionJson))
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.AddCollectionHandler(w, r)

			Convey("Then the collection store is called", func() {
				So(len(collectionStore.GetCollectionByNameCalls()), ShouldEqual, 1)
				So(collectionStore.GetCollectionByNameCalls()[0].Name, ShouldEqual, expectedName)
			})

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusConflict)
			})
		})
	})
}

func TestPostCollection_EmptyCollectionNameError(t *testing.T) {

	newCollectionJson := `{
		"name": ""
	}`
	expectedID := "12345"
	api.NewID = func() (string, error) {
		return expectedID, nil
	}

	Convey("Given a request to POST a collection with an empty collection name", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("POST", "http://localhost:26000/collections", bytes.NewBufferString(newCollectionJson))
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.AddCollectionHandler(w, r)

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

func TestPostCollection_CollectionNameLookupError(t *testing.T) {

	newCollectionJson := `{
		"name": "Coronavirus key indicators",
		"publish_date": "2020-05-05T14:58:29.317Z"
	}`
	expectedName := "Coronavirus key indicators"
	expectedID := "12345"
	api.NewID = func() (string, error) {
		return expectedID, nil
	}

	Convey("Given a request to POST a collection with a failed collection name lookup", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()
		collectionStore.GetCollectionByNameFunc = func(ctx context.Context, name string) (*models.Collection, error) {
			return nil, errors.New("well that was unexpected")
		}

		r := httptest.NewRequest("POST", "http://localhost:26000/collections", bytes.NewBufferString(newCollectionJson))
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.AddCollectionHandler(w, r)

			Convey("Then the collection store is called", func() {
				So(len(collectionStore.GetCollectionByNameCalls()), ShouldEqual, 1)
				So(collectionStore.GetCollectionByNameCalls()[0].Name, ShouldEqual, expectedName)
			})

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestPostCollection_storeError(t *testing.T) {

	newCollectionJson := `{
		"name": "Coronavirus key indicators",
		"publish_date": "2020-05-05T14:58:29.317Z"
	}`

	Convey("Given a request to POST a collection", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		collectionStore.UpsertCollectionFunc = func(ctx context.Context, collection *models.Collection) error {
			return errors.New("db is broken")
		}

		r := httptest.NewRequest("POST", "http://localhost:26000/collections", bytes.NewBufferString(newCollectionJson))
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API and an error is returned from the DB", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.AddCollectionHandler(w, r)

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestPostCollection_invalidRequestBody(t *testing.T) {

	newCollectionJson := `{`

	Convey("Given a request to POST a collection with an invalid request body", t, func() {

		paginator := mockPaginator()
		collectionStore := mockCollectionStore()

		r := httptest.NewRequest("POST", "http://localhost:26000/collections", bytes.NewBufferString(newCollectionJson))
		w := httptest.NewRecorder()

		Convey("When the request is sent to the API", func() {

			api := api.Setup(context.Background(), mux.NewRouter(), paginator, collectionStore)
			api.AddCollectionHandler(w, r)

			Convey("Then the response has the expected status code", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

func mockPaginator() *mock.PaginatorMock {

	paginator := &mock.PaginatorMock{
		ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
			return offset, limit, nil
		},
	}
	return paginator
}

func mockCollectionStore() *mock.CollectionStoreMock {

	collectionStore := &mock.CollectionStoreMock{
		GetCollectionsFunc: func(ctx context.Context, queryParams collections.QueryParams) ([]models.Collection, int, error) {
			return []models.Collection{{
				ID:          "123",
				Name:        "collection 1",
				PublishDate: &time.Time{},
				LastUpdated: time.Time{},
			}}, totalCount, nil
		},
		UpsertCollectionFunc: func(ctx context.Context, collection *models.Collection) error {
			return nil
		},
		GetCollectionByNameFunc: func(ctx context.Context, name string) (*models.Collection, error) {
			return nil, &mongodb.ErrNoDocumentFound{}
		},
		GetCollectionEventsFunc: func(ctx context.Context, queryParams collections.EventsQueryParams) ([]models.Event, int, error) {
			return []models.Event{{
				ID:           "321",
				Type:         "CREATED",
				Email:        "test@test.com",
				Date:         time.Time{},
				CollectionID: "123",
			}}, totalCount, nil
		},
		GetCollectionByIDFunc: func(ctx context.Context, id string) (*models.Collection, error) {
			return nil, nil
		},
	}

	return collectionStore
}
