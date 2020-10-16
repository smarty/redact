package sanitize

import (
    "github.com/smartystreets/assertions/should"
    "github.com/smartystreets/gunit"
    "testing"
)

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
    input := "Hello my name is John, my Credit card number is: 1111-1111-1111-1111. My employees CC number is 1111111111111111 and 1111 1111 1111 1111 plus 1111111111111."
    expectedOutput := "Hello my name is John, my Credit card number is: [MASKED 1111****1111]. My employees CC number is [MASKED 1111****1111] and [MASKED 1111****1111] plus [MASKED 1111****1]."

    output := RedactCreditCard(input)

    this.So(output, should.Resemble, expectedOutput)
}

func (this *SanitizeFixture) TestRedactSSN(){
    input := "Hello my name is John, my SSN is: 111-11-1111 my employees SSN is 111111111 and 111 11 1111."
    expectedOutput := "Hello my name is John, my SSN is: [SSN REDACTED] my employees SSN is [SSN REDACTED] and [SSN REDACTED]."

    output := RedactSSN(input)

    this.So(output, should.Resemble, expectedOutput)
}

func (this *SanitizeFixture) TestRedactTel(){
    input := "Hello my name is John, my number is: 1(801) 111-1111 and (111)111 1111 also 1111111111 one more 1-801-111-1111."
    expectedOutput := "Hello my name is John, my number is: [TEL REDACTED] and [TEL REDACTED] also [TEL REDACTED] one more [TEL REDACTED]."

    output := RedactTel(input)

    this.So(output, should.Resemble, expectedOutput)
}

func (this *SanitizeFixture) TestRedactAll(){
   input := "Hello my name is John, my email address is john@test.com. My phone-number is 1-111-111-1111, my birthday is 1/11/1111, and my CC is 1111111111111."
   expectedOutput := "Hello my name is John, my email address is [EMAIL REDACTED]. My phone-number is [TEL REDACTED], my birthday is [DOB REDACTED], and my CC is [MASKED 1111****1]."

   output := RedactAll(input)

   this.So(output, should.Resemble, expectedOutput)
}





