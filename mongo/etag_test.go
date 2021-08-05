package mongo

import (
	"github.com/ONSdigital/dp-collection-api/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGetNewETagForUpdate(t *testing.T) {

	Convey("Given a filer that we want to update", t, func() {

		newCollection := func() *models.Collection {
			return &models.Collection{
				ID:          "1234",
				Name:        "collection name",
				PublishDate: &time.Time{},
				LastUpdated: time.Time{},
			}
		}
		collection := newCollection()

		update := &models.Collection{
			Name: "updatedCollectionName",
		}

		Convey("getNewETagForUpdate returns an eTag that is different from the original collection ETag", func() {
			eTag1, err := newETagForUpdate(collection, update)
			So(err, ShouldBeNil)
			So(eTag1, ShouldNotEqual, collection.ETag)

			Convey("Applying the same update to a different collection results in a different ETag", func() {
				collection2 := newCollection()
				collection2.ID = "someOtherCollectionID"
				eTag2, err := newETagForUpdate(collection2, update)
				So(err, ShouldBeNil)
				So(eTag2, ShouldNotEqual, eTag1)
			})

			Convey("Applying a different update to the same collection results in a different ETag", func() {
				updatedTime := collection.PublishDate.Add(time.Second)
				update2 := &models.Collection{
					PublishDate: &updatedTime,
				}
				eTag3, err := newETagForUpdate(collection, update2)
				So(err, ShouldBeNil)
				So(eTag3, ShouldNotEqual, eTag1)
			})
		})
	})
}
