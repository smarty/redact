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

func (this *SanitizeFixture) TestRedactCreditCard() {
	this.testCC("Blank 5500-0000-0000-0004.", "Blank *******************.")
	this.testCC("36551639043330", "**************")
	this.testCC("4111 1111 1111 1101 111 4556-7375-8689-9855. taco ", "*********************** *******************. taco ")
}
func (this *SanitizeFixture) testCC(input, expected string) {
	this.So(this.redaction.All(input), should.Equal, expected)
}

func (this *SanitizeFixture) TestRedactEmail() {
	input := "Blah test@gmail.com, our employee's email is test@gmail. and we have one more which may or not be an email " +
		"test@test taco"
	expected := "Blah ****@gmail.com, our employee's email is ****@gmail. and we have one more which may or not be an email " +
		"****@test taco"

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
	this.testSSN("Blah 123-12-1234.", "Blah ***********.")
	this.testSSN("123121234", "*********")
	this.testSSN("123 12 1234 taco", "*********** taco")
}
func (this *SanitizeFixture) testSSN(input, expected string) {
	this.So(this.redaction.All(input), should.Equal, expected)
}

func (this *SanitizeFixture) TestRedactDOB() {
	this.testDOB("Blah 12-01-1998 and 12/01/1998 ", "Blah ********** and ********** ")
	this.testDOB("1 3 98", "******")
	this.testDOB(" March 09, 1997 and 09 May 1900 taco", " ********, 1997 and 09 ******00 taco")
}
func (this *SanitizeFixture) testDOB(input, expected string) {
	this.So(this.redaction.All(input), should.Equal, expected)
}
