package collections

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseOrderBy(t *testing.T) {

	Convey("ParseOrderBy returns default for an empty value", t, func() {
		orderBy, err := ParseOrderBy("")
		So(orderBy, ShouldEqual, OrderByDefault)
		So(err, ShouldBeNil)
	})

	Convey("ParseOrderBy parses valid values ", t, func() {
		orderBy, err := ParseOrderBy("publish_date")
		So(orderBy, ShouldEqual, OrderByPublishDate)
		So(err, ShouldBeNil)
	})

	Convey("ParseOrderBy returns an error for an unrecognised value", t, func() {
		orderBy, err := ParseOrderBy("unrecognised")
		So(orderBy, ShouldEqual, OrderByDefault)
		So(err, ShouldEqual, ErrInvalidOrderBy)
	})
}
