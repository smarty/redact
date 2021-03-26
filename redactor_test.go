package redact

import "testing"

func assertRedaction(t *testing.T, redaction *Redactor, input, expected string) {
	inputByte := []byte(input)
	actual := redaction.All(inputByte)
	if string(actual) == expected {
		return
	}
	t.Helper()
	t.Errorf("\n"+
		"Expected: %s\n"+
		"Actual:   %s",
		expected,
		actual,
	)
}

func TestRedactCreditCard(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"",
		"",
	)
	assertRedaction(t, redaction,
		"523d0555760656D3FC1D315E8",
		"523d0555760656D3FC1D315E8",
	)
	assertRedaction(t, redaction,
		"4111 1111 1111 11 10 111",
		"4111 1111 1111 11 10 111",
	)
	assertRedaction(t, redaction,
		"4111 1111 1111 1110-111",
		"4111 1111 1111 1110-111",
	)
	assertRedaction(t, redaction,
		"4111 1111 1111 1101 111 4111-1111-1111-1111. taco ",
		"*********************** *******************. taco ",
	)
	assertRedaction(t, redaction,
		"4111111111111101111",
		"*******************",
	)
	assertRedaction(t, redaction,
		"6011-0009-9013-9424 ",
		"******************* ",
	)
	assertRedaction(t, redaction,
		"1234-1234-1243-1234 ",
		"1234-1234-1243-1234 ",
	)
	assertRedaction(t, redaction,
		"4011111111111101111",
		"4011111111111101111",
	)
	assertRedaction(t, redaction,
		"6011-0o09-9013-9424 ",
		"6011-0o09-9013-9424 ",
	)
	assertRedaction(t, redaction,
		"3714-496353-98431 - O ",
		"***************** - O ",
	)
	assertRedaction(t, redaction,
		"6011-0009-9013-9424 taco.",
		"******************* taco.",
	)
	assertRedaction(t, redaction,
		"taco 6011-0009-9013-9424",
		"taco *******************",
	)
	assertRedaction(t, redaction,
		"6011 0009 9013  9424 taco.",
		"6011 0009 9013  9424 taco.",
	)
	assertRedaction(t, redaction,
		"423432343 111110101     ",
		"423432343 111110101     ",
	)
	assertRedaction(t, redaction,
		"4111111111111101111TEST",
		"*******************TEST",
	)
	assertRedaction(t, redaction,
		"411 1111 1111 1110 1111ST",
		"***********************ST",
	)
	assertRedaction(t, redaction,
		"+41 11 111 11 0",
		"+41 11 111 11 0",
	)
}
func TestRedactEmail(t *testing.T) {
	t.Parallel()

	redaction := New()

	assertRedaction(t, redaction,
		"Blah test@gmail.com, our employee's email is test@gmail. and we have one more which may or not be an email test@test taco",
		"Blah ****@gmail.com, our employee's email is ****@gmail. and we have one more which may or not be an email ****@test taco",
	)
}
func TestRedactPhone(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"801-111-1111 and (801) 111-1111 +1(801)111-1111 taco",
		"************ and (801) 111-1111 +1************* taco",
	)
	assertRedaction(t, redaction,
		"Blah 801-111-1111 and (801) 111-1111 +1(801)111-1111 taco",
		"Blah ************ and (801) 111-1111 +1************* taco",
	)
	assertRedaction(t, redaction,
		"40512-4618",
		"40512-4618",
	)
	assertRedaction(t, redaction,
		"405-124618",
		"405-124618",
	)
	assertRedaction(t, redaction,
		"This is not valid: 801 111 1111",
		"This is not valid: 801 111 1111",
	)
	assertRedaction(t, redaction,
		"801-111-1111 +1(801)111-1111 taco",
		"************ +1************* taco",
	)
}
func TestRedactSSN(t *testing.T) {
	t.Parallel()

	redaction := New()

	assertRedaction(t, redaction,
		"Blah 123-12-1234.",
		"Blah ***********.",
	)
	assertRedaction(t, redaction,
		"123 12 1234 taco",
		"*********** taco",
	)
	assertRedaction(t, redaction,
		" 123-121234 taco",
		" 123-121234 taco",
	)
	assertRedaction(t, redaction,
		"450 900 100",
		"450 900 100",
	)

}
func TestRedactDOB(t *testing.T) {
	t.Parallel()

	redaction := New()
	assertRedaction(t, redaction,
		" Apr 39 ",
		" Apr 39 ",
	)
	assertRedaction(t, redaction,
		"APRIL 3, 2019",
		"******** 2019",
	)
	assertRedaction(t, redaction,
		" 7/13/2023",
		" 7/13/2023",
	)

	assertRedaction(t, redaction,
		"1982/11/8",
		"*********",
	)
	assertRedaction(t, redaction,
		"Blah 12-01-1998 and 12/01/1998 ",
		"Blah ********** and ********** ",
	)
	assertRedaction(t, redaction,
		"Jan 1, 2021",
		"****** 2021",
	)
	assertRedaction(t, redaction,
		" February 1, 2020",
		" *********** 2020",
	)
	assertRedaction(t, redaction,
		"30-12-12",
		"30-12-12",
	)
	assertRedaction(t, redaction,
		"1/12/21",
		"1/12/21",
	)
	assertRedaction(t, redaction,
		"[5-4-212/80]",
		"[5-4-212/80]",
	)
}
