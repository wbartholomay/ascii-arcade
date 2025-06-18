package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
	input    string
	expected []string
}{
	{
		input:    "  hello  world  ",
		expected: []string{"hello", "world"},
	},
	{
		input:	  "",
		expected: []string{},
	},
	{
		input:	  "Charmander Bulbasaur PIKACHU",
		expected: []string{"charmander", "bulbasaur", "pikachu"},
	},
}

for _, c := range cases {
	actual := cleanInput(c.input)

	if len(c.expected) != len(actual) {
		t.Errorf("Length of output slice not the expected size. Expected Length: %v   Actual Length: %v. Actual Values: %v",
		len(c.expected), len(actual), actual)
		continue
	}

	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]

		if word != expectedWord {
			t.Errorf("Word is incorrect. Expected Word: %v   Actual Word: %v", expectedWord, word)
		}
	}
}
}