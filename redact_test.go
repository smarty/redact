package redact

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestSanitizeFixture(t *testing.T) {
	gunit.Run(new(SanitizeFixture), t, gunit.Options.AllSequential())
}

type SanitizeFixture struct {
	*gunit.Fixture
	redaction *Redaction
}

func (this *SanitizeFixture) Setup() {
	this.redaction = New()
}

func (this *SanitizeFixture) FocusTestRedactCreditCard() {
	input := "Blah 5500-0000-0000-0004. CC number is 36551639043330 and 4111 1111 1111 1101 111 4556-7375-8689-9855. taco "
	expected := "Blah *******************. CC number is ************** and *********************** *******************. taco "

	actual := this.redaction.All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestRedactEmail() {
	input := "Blah test@gmail.com, our employee's email is test@gmail. and we have one more which may or not be an email " +
		"test@test "
	expected := "Blah **********.com, our employee's email is **********. and we have one more which may or not be an email " +
		"********* "
	// TODO: Test with taco
	// TODO: No space

	actual := this.redaction.All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestRedactPhoneNum() {
	input := "Blah 801-111-1111 and 801 111 1111 and (801) 111-1111 +1(801)111-1111 taco"

	expected := "Blah ************ and ************ and ************** +1************* taco"

	actual := this.redaction.All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestRedactSSN() {
	input := "Blah 123-12-1234 and 123121234 or 123 12 1234 taco"

	expected := "Blah *********** and ********* or *********** taco"

	actual := this.redaction.All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestRedactDOB() {
	input := "Blah 12-01-1998 and 12/01/1998 or 1 3 98 and March 09, 1997 and 09 May 1900 taco"

	expected := "Blah ********** and ********** or ****** and ********, 1997 and 09 ******00 taco"

	actual := this.redaction.All(input)

	this.So(actual, should.Equal, expected)
}
