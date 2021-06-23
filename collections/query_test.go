package collections

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateNameSearchInput(t *testing.T) {

	Convey("ValidateNameSearchInput returns nil error for valid values", t, func() {
		So(ValidateNameSearchInput(""), ShouldBeNil)
		So(ValidateNameSearchInput("collection123"), ShouldBeNil)
	})

	Convey("ValidateNameSearchInput returns an error for a value over 64 characters", t, func() {
		tooLongInput := "1234567890123456789012345678901234567890123456789012345678901234567890"
		So(ValidateNameSearchInput(tooLongInput), ShouldEqual, ErrNameSearchTooLong)
	})
}
