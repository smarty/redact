package sanitize

import (
    "github.com/smartystreets/assertions/should"
    "github.com/smartystreets/gunit"
    "testing"
)

const TEST_TEXT1 = "Hello my name is John, my email address is john@test.com and my birthday is January 1, 2000."

func TestSanitizeFixture(t *testing.T) {
    gunit.Run(new(SanitizeFixture), t)
}
type SanitizeFixture struct {
    *gunit.Fixture
}

func (this *SanitizeFixture) Setup() {
}

func (this *SanitizeFixture) TestRedactDOB() {
    input := "Hello my name is John, my date of birth is 11/1/2000 and my employee's date of birth is 01-01-2001, oh also November 1, 2000, May 23, 2019, 23 June 1989."
    expectedOutput := "Hello my name is John, my date of birth is [DOB REDACTED] and my employee's date of birth is [DOB REDACTED], oh also [DOB REDACTED], [DOB REDACTED], [DOB REDACTED]."

    output := RedactDateOfBirth(input)

    this.So(output, should.Resemble, expectedOutput)
}

func (this *SanitizeFixture) TestRedactEmail() {
    input := "Hello my name is John, my email address is john@test.com and my employee's email is jake@test.com."
    expectedOutput := "Hello my name is John, my email address is [EMAIL REDACTED] and my employee's email is [EMAIL REDACTED]."

    output := RedactEmail(input)

    this.So(output, should.Resemble, expectedOutput)
}

func (this *SanitizeFixture) TestRedactCC(){
    input := "Hello my name is John, my Credit card number is: 1111-1111-1111-1111. My employees CC number is 1111111111111111."
    expectedOutput := "Hello my name is John, my Credit card number is: [MASKED 1111****1111]. My employees CC number is [MASKED 1111****1111]."

    output := RedactCreditCard(input)
    this.So(output, should.Resemble, expectedOutput)
}

func (this *SanitizeFixture) TestSanitizePhone() {
    sanitized := RedactPhone("111-111-1111")
    this.So(sanitized, should.Equal, "[TEL REDACTED]")
}

func (this *SanitizeFixture) TestSanitizeSSN() {
    sanitized := RedactSSN("111-11-1111")
    this.So(sanitized, should.Equal, "[SSN REDACTED]")
}

func (this *SanitizeFixture) SkipTestSanitizeCreditCardBasic() {
    sanitized := RedactCreditCard("1111111111111111")
    this.So(sanitized, should.Equal, "[MASKED 1111****1111]")
}

func (this *SanitizeFixture) SkipTestSanitizeCreditCardDashes() {
    sanitized := RedactCreditCard("1111-1111-1111-1111")
    this.So(sanitized, should.Equal, "[MASKED 1111****1111]")
}

func (this *SanitizeFixture) SkipTestSanitizeCreditCardSpaces() {
    sanitized := RedactCreditCard("1111 1111 1111 1111")
    this.So(sanitized, should.Equal, "[MASKED 1111****1111]")
}





