package models

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCollectionHash(t *testing.T) {

	testCollection := func() Collection {
		return Collection{
			ID:          "1234",
			Name:        "collection name",
			PublishDate: &time.Time{},
			LastUpdated: time.Time{},
		}
	}

	Convey("Given a Collection with some data", t, func() {
		Collection := testCollection()

		Convey("We can generate a valid hash", func() {
			h, err := Collection.Hash(nil)
			So(err, ShouldBeNil)
			So(len(h), ShouldEqual, 40)

			Convey("Then hashing it twice, produces the same result", func() {
				hash, err := Collection.Hash(nil)
				So(err, ShouldBeNil)
				So(hash, ShouldEqual, h)
			})

			Convey("Then storing the hash as its ETag value and hashing it again, produces the same result (field is ignored) and ETag field is preserved", func() {
				Collection.ETag = h
				hash, err := Collection.Hash(nil)
				So(err, ShouldBeNil)
				So(hash, ShouldEqual, h)
				So(Collection.ETag, ShouldEqual, h)
			})

			Convey("Then another Collection with exactly the same data will resolve to the same hash", func() {
				Collection2 := testCollection()
				hash, err := Collection2.Hash(nil)
				So(err, ShouldBeNil)
				So(hash, ShouldEqual, h)
			})

			Convey("Then if a Collection value is modified, its hash changes", func() {
				Collection.Name = "new collection name"
				hash, err := Collection.Hash(nil)
				So(err, ShouldBeNil)
				So(hash, ShouldNotEqual, h)
			})
		})
	})
}

func TestGetNewETagForUpdate(t *testing.T) {

	Convey("Given a filer that we want to update", t, func() {

		newCollection := func() *Collection {
			return &Collection{
				ID:          "1234",
				Name:        "collection name",
				PublishDate: &time.Time{},
				LastUpdated: time.Time{},
			}
		}
		collection := newCollection()

		update := &Collection{
			Name: "updatedCollectionName",
		}

		Convey("getNewETagForUpdate returns an eTag that is different from the original collection ETag", func() {
			eTag1, err := collection.NewETagForUpdate(update)
			So(err, ShouldBeNil)
			So(eTag1, ShouldNotEqual, collection.ETag)

			Convey("Applying the same update to a different collection results in a different ETag", func() {
				collection2 := newCollection()
				collection2.ID = "someOtherCollectionID"
				eTag2, err := collection2.NewETagForUpdate(update)
				So(err, ShouldBeNil)
				So(eTag2, ShouldNotEqual, eTag1)
			})

			Convey("Applying a different update to the same collection results in a different ETag", func() {
				updatedTime := collection.PublishDate.Add(time.Second)
				update2 := &Collection{
					PublishDate: &updatedTime,
				}
				eTag3, err := collection.NewETagForUpdate(update2)
				So(err, ShouldBeNil)
				So(eTag3, ShouldNotEqual, eTag1)
			})
		})
	})
}
