// Copyright (c) 2017-2020 VMware, Inc. or its affiliates
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"reflect"
	"testing"
)

func TestParsePorts(t *testing.T) {
	cases := []struct {
		input    string
		expected []uint32
	}{
		{"", []uint32(nil)},
		{"1", []uint32{1}},
		{"1,3,5", []uint32{1, 3, 5}},
		/* ranges */
		{"1-5", []uint32{1, 2, 3, 4, 5}},
		{"1-5,6-10", []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"1-5,10,12,15-15", []uint32{1, 2, 3, 4, 5, 10, 12, 15}},
	}

	for _, c := range cases {
		actual, err := parsePorts(c.input)
		if err != nil {
			t.Errorf("parsePorts(%q) returned error %#v", c.input, err)
		}
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("parsePorts(%q) returned %v, want %v", c.input, actual, c.expected)
		}
	}

	errorCases := []string{
		"1, 3, 5",
		"sdklfjds",
		"-1",
		"5-1",
		"1--5",
		"1-3-5",
		"1,,2",
		"1,a",
		"1-a",
		"a-1",
		"900000",
		"1-900000",
		"900000-1000000",
		",1",
	}

	for _, c := range errorCases {
		actual, err := parsePorts(c)
		if err == nil {
			t.Errorf("parsePorts(%q) returned %v instead of an error", c, actual)
		}
	}
}

func TestIsLinkMode(t *testing.T) {
	cases := []struct {
		name     string
		mode     string
		expected bool
	}{
		{
			name:     "parses copy",
			mode:     "copy",
			expected: false,
		},
		{
			name:     "parses link",
			mode:     "link",
			expected: true,
		},
		{
			name:     "parses capitalizations",
			mode:     "LiNk",
			expected: true,
		},
		{
			name:     "trims spaces",
			mode:     " link  \t",
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			linkMode, err := isLinkMode(c.mode)
			if err != nil {
				t.Errorf("unexpected error %#v", err)
			}

			if linkMode != c.expected {
				t.Errorf("got %t want %t", linkMode, c.expected)
			}
		})
	}

	errCases := []struct {
		name string
		mode string
	}{
		{
			name: "empty string",
			mode: "",
		},
		{
			name: "invalid mode",
			mode: "depeche",
		},
		{
			name: "errors on numbers",
			mode: "1",
		},
	}

	for _, c := range errCases {
		t.Run(c.name, func(t *testing.T) {
			linkMode, err := isLinkMode(c.mode)
			if err == nil {
				t.Errorf("isLinkMode(%q) returned %v instead of an error", c.mode, err)
			}

			if linkMode != false {
				t.Errorf("got linkMode %t want %t", linkMode, false)
			}
		})
	}
}
