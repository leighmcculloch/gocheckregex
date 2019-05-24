package main

import (
	"testing"
)

func TestCheckNoGlobals(t *testing.T) {
	cases := []struct {
		path         string
		wantMessages []string
	}{
		{
			path:         "testdata/0",
			wantMessages: nil,
		},
		{
			path:         "testdata/0/code.go",
			wantMessages: nil,
		},
		{
			path:         "testdata/1",
			wantMessages: nil,
		},
		{
			path: "testdata/2",
			wantMessages: []string{
				"testdata/2/code.go:5 error parsing regexp: missing closing ]: `[0-9`",
			},
		},
		{
			path: "testdata/3",
			wantMessages: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			messages, err := checkRegex(c.path)
			if err != nil {
				t.Fatalf("got error %#v", err)
			}
			if !stringSlicesEqual(messages, c.wantMessages) {
				t.Errorf("got %#v, want %#v", messages, c.wantMessages)
			}
		})
	}
}

func stringSlicesEqual(s1, s2 []string) bool {
	diff := map[string]int{}
	for _, s := range s1 {
		diff[s]++
	}
	for _, s := range s2 {
		diff[s]--
		if diff[s] == 0 {
			delete(diff, s)
		}
	}
	return len(diff) == 0
}
