package stringer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var stringerStudentTests = []struct {
	in   Student
	want string
}{
	{Student{
		ID:        42,
		FirstName: "John",
		LastName:  "Doe",
		Age:       25,
	}, "Student ID: 42. Name: Doe, John. Age: 25."},
	{Student{
		ID:        1490,
		FirstName: "Tormod",
		LastName:  "Lea",
		Age:       30,
	}, "Student ID: 1490. Name: Lea, Tormod. Age: 30."},
}

func TestStringerStudent(t *testing.T) {
	for _, test := range stringerStudentTests {
		if diff := cmp.Diff(test.want, test.in.String()); diff != "" {
			t.Errorf("String(%q): (-want +got):\n%s", test.in, diff)
		}
	}
}
