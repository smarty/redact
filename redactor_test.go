package redact

import "testing"

func assertRedaction(t *testing.T, redaction *Redactor, input, expected string) {
	inputByte := []byte(input)
	actual := redaction.RedactAll(inputByte)
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
func BenchmarkThing(b *testing.B) {
	redaction := New()
	b.ReportAllocs()
	b.ResetTimer()
	input := []byte("+8014890464 Hello 6749-3-2345 (801)4890464there, my name is stuff. 1200 East 1200 North Mapleton " +
		"1(385)6668330. 18014890464 371449635398431numMayers 12, 1970ber is 647-48-6867. I " +
		"9/1/2020 to fill these wo04/16/1999rds in with other 371 449 635 398 1945/05/01 " +
		"431impoletsgitit@yahoo.com. 4111-111-111-111-111 pasting 02/14/1900the 647-21-12398 best of6011111111111117 the" +
		"valid, and Jan 32, 1990someMarch 12, 2020 are not. 647489009 This is a vDecember 111, 2000ey fun task to do. " +
		"647 40 4444 1+(801)4890464 647-48-9012")
	for n := 0; n < b.N; n++ {
		_ = redaction.RedactAll(input)
	}
}

func TestRedactCreditCard_Valid_Redaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"4111 1111 1111 1101 111 4111-1111-1111-1111. taco ",
		"*********************** *******************. taco ",
	)
	assertRedaction(t, redaction,
		"taco 6011-0009-9013-9424",
		"taco *******************",
	)
	assertRedaction(t, redaction,
		"6011-0009-9013-9424 ",
		"******************* ",
	)
	assertRedaction(t, redaction,
		" 4111111111111101111 ",
		" ******************* ",
	)
	assertRedaction(t, redaction,
		"4111111111111101111",
		"*******************",
	)
	assertRedaction(t, redaction,
		"taco 3714-496353-98431 - O ",
		"taco ***************** - O ",
	)
	assertRedaction(t, redaction,
		"6011-0009-9013-9424 taco.",
		"******************* taco.",
	)
	assertRedaction(t, redaction,
		"4111111111111101111TEST",
		"*******************TEST",
	)
	assertRedaction(t, redaction,
		"411 1111 1111 1110 1111ST",
		"***********************ST",
	)
}
func TestRedactCreditCard_Invalid_NoRedaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"6011 0009-9013-9424 taco.",
		"6011 0009-9013-9424 taco.",
	)
	assertRedaction(t, redaction,
		"601100 099 013 9424 taco.",
		"601100 099 013 9424 taco.",
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
		"+41 11 111 11 0",
		"+41 11 111 11 0",
	)
	assertRedaction(t, redaction,
		"93087097919637351852097735514114210893391460226901438681143583309305538438469497356679189567980884261"+
			"29971266084194311662842664602358893723650394247864792157944906026710951003240151287380372698948380361"+
			"55164632037694293000947090031773467719395901857854689140563439380055959763377841900038677981407681556"+
			"31549026189478389556566369157325037087613078605916540148639313073551273625059514900828293261488835150"+
			"67313603759318833137870438235772016686252244297993834557723091360234446940034078073980985453311613448"+
			"992144713703560551680141615000380747919129581233295746609790127688737740588379751",

		"93087097919637351852097735514114210893391460226901438681143583309305538438469497356679189567980884261"+
			"29971266084194311662842664602358893723650394247864792157944906026710951003240151287380372698948380361"+
			"55164632037694293000947090031773467719395901857854689140563439380055959763377841900038677981407681556"+
			"31549026189478389556566369157325037087613078605916540148639313073551273625059514900828293261488835150"+
			"67313603759318833137870438235772016686252244297993834557723091360234446940034078073980985453311613448"+
			"992144713703560551680141615000380747919129581233295746609790127688737740588379751",
	)
}

func TestRedactEmail_Valid_Redaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"test@test.com",
		"****@test.com",
	)
	assertRedaction(t, redaction,
		"Blah test.test@gmail.com, our employee's email is test@gmail. and we have one more which may or not be an email test@test taco",
		"Blah *********@gmail.com, our employee's email is ****@gmail. and we have one more which may or not be an email ****@test taco",
	)
}
func TestRedactEmail_Invalid_NoRedaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"Blah test.gmail.com",
		"Blah test.gmail.com",
	)
}

func TestRedactPhone_Valid_Redaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"801-111-1111 +1(801)111-1111 taco",
		"************ +************** taco",
	)
	assertRedaction(t, redaction,
		"+1(801)111-1111 taco",
		"+************** taco",
	)
	assertRedaction(t, redaction,
		"801-111-1111 and (801) 111-1111 +1(801)111-1111 taco",
		"************ and (801) 111-1111 +************** taco",
	)
	assertRedaction(t, redaction,
		"Blah 801-111-1111 and (801) 111-1111 +1(801)111-1111 taco",
		"Blah ************ and (801) 111-1111 +************** taco",
	)

}
func TestRedactPhone_Invalid_NoRedaction(t *testing.T) {
	t.Parallel()
	redaction := New()
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
}

func TestRedactSSN_Valid_Redaction(t *testing.T) {
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
}
func TestRedactSSN_Invalid_NoRedaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		" 123-121234 taco",
		" 123-121234 taco",
	)
	assertRedaction(t, redaction,
		"450 900 100",
		"450 900 100",
	)

}

func TestRedactDOB_Valid_Redaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		"stuff February 1, 2020 ",
		"stuff **************** ",
	)
	assertRedaction(t, redaction,
		"Feb 01, 2012",
		"************",
	)
	assertRedaction(t, redaction,
		" APRIL 3, 2019",
		" *************",
	)

	assertRedaction(t, redaction,
		"Blah 12-01-1998 and 12/01/1998 ",
		"Blah ********** and ********** ",
	)
	assertRedaction(t, redaction,
		"Blah 12-12-1998 and 01/01/1998 ",
		"Blah ********** and ********** ",
	)
	assertRedaction(t, redaction,
		"1982/11/8",
		"*********",
	)
	assertRedaction(t, redaction,
		"Jan 1, 2021 ",
		"*********** ",
	)
}
func TestRedactDOB_Invalid_NoRedaction(t *testing.T) {
	t.Parallel()
	redaction := New()
	assertRedaction(t, redaction,
		" Apr 39, 2021 ",
		" Apr 39, 2021 ",
	)
	assertRedaction(t, redaction,
		"April 21, 2025",
		"April 21, 2025",
	)
	assertRedaction(t, redaction,
		" 7/13/2023",
		" 7/13/2023",
	)
	assertRedaction(t, redaction,
		"30-12-12",
		"30-12-12",
	)
	assertRedaction(t, redaction,
		"1/12/2123",
		"1/12/2123",
	)
	assertRedaction(t, redaction,
		"[5-4-212/80]",
		"[5-4-212/80]",
	)
}
