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
	input := "Blah 1111-1111-1111-1111. CC number is 1111111111111111 and 1111 1111 1111 1111 2222-2222-2222-2222"
	//My employees CC number is 1111111111111111 and 1111 1111 1111 1111 plus 1111111111111."
	this.So(matchCreditCard(input),should.Resemble, []match{{
		InputIndex: 5,
		Length:     19,
	}, {
		InputIndex: 39,
		Length:     16,
	}, {
		InputIndex: 60,
		Length:     19,
	}, {
		InputIndex: 80,
		Length:     19,
	}})
}
//
//
//func (this *SanitizeFixture) TestRedactSSN() {
//	input := "Hello my name is John, my SSN is: 111-11-1111 my employees SSN is 111111111 and 111 11 1111."
//	expectedOutput := "Hello my name is John, my SSN is: [SSN REDACTED] my employees SSN is [SSN REDACTED] and [SSN REDACTED]."
//
//	output := SSN(input)
//
//	this.So(output, should.Resemble, expectedOutput)
//}
//
//func (this *SanitizeFixture) TestRedactTelephone() {
//	input := "Hello my name is John, my number is: 1(801) 111-1111 and (111)111 1111 also 1111111111 one more 1-801-111-1111."
//	expectedOutput := "Hello my name is John, my number is: [PHONE REDACTED] and [PHONE REDACTED] also [PHONE REDACTED] one more [PHONE REDACTED]."
//
//	output := Phone(input)
//
//	this.So(output, should.Resemble, expectedOutput)
//}
//
//func (this *SanitizeFixture) TestRedactAll() {
//	input := "Hello my name is John, my email address is john@test.com. My phone-number is 1-111-111-1111, my birthday is 1/11/1111, and my CC is 1111111111111."
//	expectedOutput := "Hello my name is John, my email address is [EMAIL REDACTED]. My phone-number is [PHONE REDACTED], my birthday is [DOB REDACTED], and my CC is [CARD 1111****1111]."
//
//	output := All(input)
//
//	this.So(output, should.Resemble, expectedOutput)

