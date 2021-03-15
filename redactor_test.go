package redact

import "testing"

//TODO: test why it doesn't catch the space in front? <- thought we fixed that guess not.
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
		"taco 6556-7375-8689-9835",
		"taco *******************",
	)
	assertRedaction(t, redaction,
		"",
		"",
	)
	assertRedaction(t, redaction, //Invalid credit card not separate from letters
		"52353330555760656D3FC1D315E80069",
		"52353330555760656D3FC1D315E80069",
		//TODO: Looks like a valid card from the right but not all the way- I.e. flip this test
	)
	assertRedaction(t, redaction, //Invalid length (>=20) for credit card.
		"41111111111111011101",
		"41111111111111011101",
	)
	assertRedaction(t, redaction, //Invalid number of spaces (==0 || >1 || <5)
		"411111111111110 1111",
		"411111111111110 1111",
	)
	assertRedaction(t, redaction, ///Invalid number of spaces (==0 || >1 || <5)
		"4111 1111 1111 11 10 111",
		"4111 1111 1111 11 10 111",
	)
	assertRedaction(t, redaction, //Invalid two different types of spaces
		"4111 1111 1111 1110-111",
		"4111 1111 1111 1110-111",
	)
	assertRedaction(t, redaction,
		"4111 1111 1111 1101 111 4556-7375-8689-9855. taco ",
		"*********************** *******************. taco ",
	)
	assertRedaction(t, redaction,
		"4111111111111101111",
		"*******************",
	)
	assertRedaction(t, redaction,
		"4556-7375-8689-9855 ",
		"******************* ",
	)
	assertRedaction(t, redaction,
		"2556-7375-8689-9855 ",
		"2556-7375-8689-9855 ",
	)
	assertRedaction(t, redaction,
		"6556-7375-8689-9835 ",
		"******************* ",
	)
	//Test for something that looks like a credit card but isn't
	assertRedaction(t, redaction,
		"4011111111111101111",
		"4011111111111101111",
	)
	assertRedaction(t, redaction,
		"6556-7371-8689-9835 ",
		"6556-7371-8689-9835 ",
	)
	//Something that has a random letter mixed in
	assertRedaction(t, redaction,
		"6556-7a75-8689-9835 ",
		"6556-7a75-8689-9835 ",
	)
	assertRedaction(t, redaction,
		"6y56737586899835 ",
		"6y56737586899835 ",
	)
	//Beginning of a string, stuff after,
	assertRedaction(t, redaction,
		"6556-7375-8689-9835 theoretically this would be an address.",
		"******************* theoretically this would be an address.",
	)

	//TODO: add tests, this isn't enough.
	//Different lengths of credit cards
	//Check for EVERYTHING.
	//Check coverage for the lund check,
	//Test for double spaces
	//Check and fix test for space in front of card
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
		"[329993740 873518800     ]",
		"[329993740 873518800     ]",
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
