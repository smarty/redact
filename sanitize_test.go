package sanitize

import (
    "github.com/smartystreets/assertions/should"
    "github.com/smartystreets/gunit"
    "testing"
)

const TEST_TEXT1 = "Hello my name is John, my email address is john@test.com and my birthday is January 1, 2000."
const TEST_TEXT2 = "Hello my name is John, my email address is john@test.com. My birthday is January 1, 2000. My "
func TestSanitizeFixture(t *testing.T) {
    gunit.Run(new(SanitizeFixture), t)
}
type SanitizeFixture struct {
    *gunit.Fixture
}

func (this *SanitizeFixture) Setup() {
}

func (this *SanitizeFixture) TestSanitizeDOB() {
    sanitized := RedactDateOfBirth("November 1, 2000")
    this.So(sanitized, should.Equal, "[DOB REDACTED]")
}

func (this *SanitizeFixture) TestSanitizeEmail() {
    sanitized := RedactEmail("user@test.com")
    this.So(sanitized, should.Equal, "[EMAIL REDACTED]")
}
func (this *SanitizeFixture)TestFindEmail() {
    emails := FindEmails(TEST_TEXT2)
    this.So(emails[0], should.Equal, "john@test.com")
}

func (this *SanitizeFixture) TestSanitizePhone() {
    sanitized := RedactPhone("111-111-1111")
    this.So(sanitized, should.Equal, "[TEL REDACTED]")
}

func (this *SanitizeFixture) TestSanitizeSSN() {
    sanitized := RedactSSN("111-11-1111")
    this.So(sanitized, should.Equal, "[SSN REDACTED]")
}

func (this *SanitizeFixture) TestSanitizeCreditCardBasic() {
    sanitized := RedactCreditCard("1111111111111111")
    this.So(sanitized, should.Equal, "[MASKED 1111****1111]")
}

func (this *SanitizeFixture) TestSanitizeCreditCardDashes() {
    sanitized := RedactCreditCard("1111-1111-1111-1111")
    this.So(sanitized, should.Equal, "[MASKED 1111****1111]")
}

func (this *SanitizeFixture) TestSanitizeCreditCardSpaces() {
    sanitized := RedactCreditCard("1111 1111 1111 1111")
    this.So(sanitized, should.Equal, "[MASKED 1111****1111]")
}





