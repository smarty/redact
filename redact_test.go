package redact

import "testing"

func assertRedaction(t *testing.T, redaction *Redaction, input, expected string) {
	actual := redaction.All(input)
	if actual == expected {
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
		"Blank 5500-0000-0000-0004.",
		"Blank *******************.",
	)
	assertRedaction(t, redaction,
		"36551639043330",
		"**************",
	)
	assertRedaction(t, redaction,
		"4111 1111 1111 1101 111 4556-7375-8689-9855. taco ",
		"*********************** *******************. taco ",
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
		"Blah 801-111-1111 and 801 111 1111 and (801) 111-1111 +1(801)111-1111 taco",
		"Blah ************ and ************ and ************** +1************* taco",
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
		"123121234",
		"*********",
	)
	assertRedaction(t, redaction,
		"123 12 1234 taco",
		"*********** taco",
	)
}
func TestRedactDOB(t *testing.T) {
	t.Parallel()

	redaction := New()

	assertRedaction(t, redaction,
		"Blah 12-01-1998 and 12/01/1998 ",
		"Blah ********** and ********** ",
	)
	assertRedaction(t, redaction,
		"1 3 98",
		"******",
	)
	assertRedaction(t, redaction,
		" March 09, 1997 and 09 May 1900 taco",
		" ********, 1997 and 09 ******00 taco",
	)
	assertRedaction(t, redaction,
		"1234    ",
		"1234    ",
	)
}
