package api_test

import (
	"context"
	"github.com/ONSdigital/dp-collection-api/api"
	"github.com/ONSdigital/dp-collection-api/api/mock"
	"github.com/ONSdigital/dp-collection-api/collections"
	"github.com/ONSdigital/dp-collection-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {

	paginator := &mock.PaginatorMock{
		ReadPaginationParametersFunc: func(r *http.Request) (int, int, error) {
			return 1, 0, nil
		},
	}

	collectionStore := &mock.CollectionStoreMock{
		GetCollectionsFunc: func(ctx context.Context, offset int, limit int, orderBy collections.OrderBy) ([]models.Collection, int, error) {
			return []models.Collection{}, 0, nil
		},
	}

	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		api := api.Setup(ctx, r, paginator, collectionStore)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/collections", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
