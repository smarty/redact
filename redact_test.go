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
}

func (this *SanitizeFixture) Teardown() {
	used = make(map[int]struct{})
}


func (this *SanitizeFixture) TestMatchCreditCard() {
	input := "Blah 4556-7375-8689-9855. CC number is 36551639043330 and 4556 3172 3465 5089 670 4556-7375-8689-9855 taco"
	this.So(matchCreditCard(input), should.Resemble, []match{{
		InputIndex: 5,
		Length:     19,
	}, {
		InputIndex: 39,
		Length:     14,
	}, {
		InputIndex: 58,
		Length:     23,
	}, {
		InputIndex: 82,
		Length:     19,
	}})
}

func (this *SanitizeFixture) TestRedactCreditCard() {
	input := "Blah 4556-7375-8689-9855. CC number is 36551639043330 and 4556 3172 3465 5089 670 4556-7375-8689-9855 taco "
	expected := "Blah *******************. CC number is ************** and *********************** ******************* taco "

	actual := All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestMatchEmail() {
	input := "Blah test@gmail.com, our employee's email is test@gmail. and we have one more which may or not be an email " +
		"test@test. taco"
	this.So(matchEmail(input), should.Resemble, []match{{
		InputIndex: 5,
		Length:     10,
	}, {
		InputIndex: 45,
		Length:     10,
	}, {
		InputIndex: 107,
		Length:     9,
	}})
}

func (this *SanitizeFixture) TestRedactEmail() {
	input := "Blah test@gmail.com, our employee's email is test@gmail. and we have one more which may or not be an email " +
		"test@test. taco"
	expected := "Blah **********.com, our employee's email is **********. and we have one more which may or not be an email " +
		"*********. taco"

	actual := All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestMatchPhoneNum() {
	input := "Blah 801-111-1111 and 801 111 1111 and (801) 111-1111 +1(801)111-1111"
	this.So(matchPhoneNum(input), should.Resemble, []match{
		{
			InputIndex: 5,
			Length:     12,
		},
		{
			InputIndex: 22,
			Length:     12,
		},
		{
			InputIndex: 39,
			Length:     14,
		},
		{
			InputIndex: 56,
			Length:     13,
		},
	})
}
func (this *SanitizeFixture) TestRedactPhoneNum() {
	input := "Blah 801-111-1111 and 801 111 1111 and (801) 111-1111 +1(801)111-1111 taco"

	expected := "Blah ************ and ************ and ************** +1************* taco"

	actual := All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestMatchSSN() {
	input := "Blah 123-12-1234 and 123121234 or 123 12 1234 taco"

	this.So(matchSSN(input), should.Resemble, []match{
		{
			InputIndex: 5,
			Length:     11,
		},
		{
			InputIndex: 21,
			Length:     9,
		},
		{
			InputIndex: 34,
			Length:     11,
		},
	})
}
func (this *SanitizeFixture) TestRedactSSN() {
	input := "Blah 123-12-1234 and 123121234 or 123 12 1234 taco"

	expected := "Blah *********** and ********* or *********** taco"

	actual := All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestMatchDOB() {
	input := "Blah 12-01-1998 and 12/01/1998 or 1 3 98 and March 09, 1997 and 09 May 1900 taco"

	this.So(matchDOB(input), should.Resemble, []match{
		{
			InputIndex: 5,
			Length:     10,
		},
		{
			InputIndex: 20,
			Length:     10,
		},
		{
			InputIndex: 34,
			Length:     6,
		},
		{
			InputIndex: 45,
			Length:     8,
		},
		{
			InputIndex: 67,
			Length:     6,
		},
	})
}

func (this *SanitizeFixture) TestRedactDOB() {
	input := "Blah 12-01-1998 and 12/01/1998 or 1 3 98 and March 09, 1997 and 09 May 1900 taco"

	expected := "Blah ********** and ********** or ****** and ********, 1997 and 09 ******00 taco"

	actual := All(input)

	this.So(actual, should.Equal, expected)
}
