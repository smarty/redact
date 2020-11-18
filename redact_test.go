package redact

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestSanitizeFixture(t *testing.T) {
	gunit.Run(new(SanitizeFixture), t)
}

type SanitizeFixture struct {
	*gunit.Fixture
}

//func (this *SanitizeFixture) TestRedactDOB() {
//	input := "Hello my name is John, my date of birth is 11/1/2000 and my employee's date of birth is 01-01-2001, oh also November 1, 2000, May 23, 2019, 23 June 1989, Sept 4, 2010."
//	expectedOutput := "Hello my name is John, my date of birth is [DOB REDACTED] and my employee's date of birth is [DOB REDACTED], oh also [DOB REDACTED], [DOB REDACTED], [DOB REDACTED], [DOB REDACTED]."
//
//	output := DateOfBirth(input) // TODO: change to ALL
//
//	this.So(output, should.Resemble, expectedOutput)
//}
//
//func (this *SanitizeFixture) TestRedactEmail() {
//	input := "Hello my name is John, my email address is john@test.com and my employee's email is jake@test.com and Jake Smith <jake@smith.com>."
//	expectedOutput := "Hello my name is John, my email address is [EMAIL REDACTED] and my employee's email is [EMAIL REDACTED] and Jake Smith <[EMAIL REDACTED]>."
//
//	output := Email(input)
//
//	this.So(output, should.Resemble, expectedOutput)
//}
//
//func (this *SanitizeFixture) SkipTestRedactCreditCard() {
//	input := "Hello my name is John, my Credit card number is: 1111-1111-1111-1111. My employees CC number is 1111111111111111 and 1111 1111 1111 1111 plus 1111111111111."
//	expectedOutput := "Hello my name is John, my Credit card number is: [CARD 1111****1111]. My employees CC number is [CARD 1111****1111] and [CARD 1111****1111] plus [CARD 1111****1111]."
//
//	output := CreditCard(input)
//
//	this.So(output, should.Resemble, expectedOutput)
//}

func (this *SanitizeFixture) TestMatchCreditCard() {
	input := "Blah 4556-7375-8689-9855. CC number is 36551639043330 and 4556 3172 3465 5089 670 6011-7674-3539-9843"
	this.So(matchCreditCard(input), should.Resemble, []match{{
		InputIndex: 5,
		Length:     19,
	}, {
		InputIndex: 39,
		Length:     14,
	}, {
		InputIndex: 58,
		Length:     19,
	}, {
		InputIndex: 78,
		Length:     19,
	}})
}

func (this *SanitizeFixture) TestRedactCreditCard() {
	input := "Blah 4556-7375-8689-9855. CC number is 36551639043330 and 4556 3172 3465 5089 670 6011-7674-3539-9843"
	expected := "Blah *******************. CC number is ************** and ******************* *******************"

	actual := All(input)

	this.So(actual, should.Equal, expected)
}

func (this *SanitizeFixture) TestMatchEmail() {
	input := "Blah test@gmail.com, our employee's email is test@gmail. and we have one more which may or not be an email " +
		"test@test."
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
		"test@test."
	expected := "Blah **********.com, our employee's email is **********. and we have one more which may or not be an email " +
		"*********."

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
	input := "Blah 801-111-1111 and 801 111 1111 and (801) 111-1111 +1(801)111-1111"

	expected := "Blah ************ and ************ and ************** +1*************"

	actual := All(input)

	this.So(actual, should.Equal, expected)
}

