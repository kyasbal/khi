// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parserutil

import "strings"

var ANSIBeginSequences []string = []string{
	"\\x1b[",
	"\\033[",
	"\\u001B[",
}

// SpecialSequenceStripper removes specific sequences from string. (e.g ASCII escape characters, ANSI color characters)
type SpecialSequenceStripper interface {
	// Strip receives the original string and returns stripped string.
	Strip(s string) string
}

// StripSpecialSequences returns the stripped string with applying provided strippers to the original string.
func StripSpecialSequences(original string, stripper ...SpecialSequenceStripper) string {
	for _, s := range stripper {
		original = s.Strip(original)
	}
	return original
}

type SequenceStripper struct {
	StripTarget []string
}

func NewSequenceStripper(targets ...string) *SequenceStripper {
	return &SequenceStripper{
		StripTarget: targets,
	}
}

func (c *SequenceStripper) Strip(s string) string {
	for _, target := range c.StripTarget {
		s = strings.ReplaceAll(s, target, "")
	}
	return s
}

var _ SpecialSequenceStripper = (*SequenceStripper)(nil)

type ANSIEscapeSequenceStripper struct {
}

// Strip implements SpecialSequenceStripper.
func (a *ANSIEscapeSequenceStripper) Strip(s string) string {
	builder := strings.Builder{}
	for i := 0; i < len(s); i++ {
		ansiFound := false
		nextFound := len(s)
		for _, beginSequence := range ANSIBeginSequences {
			nextFoundForSequence := strings.Index(s[i:], beginSequence)
			if nextFoundForSequence != -1 && nextFoundForSequence < nextFound {
				nextFound = nextFoundForSequence
			}
		}
		if nextFound != len(s) {
			ansiFound = true
			foundSuffix := false
			builder.WriteString(s[i : i+nextFound])
			i += nextFound
			for j := i; j < len(s); j++ {
				if s[j] == 'm' {
					i = j
					foundSuffix = true
					break
				}
			}
			if !foundSuffix { // This ANSI sequence is not complete. Write it as is.
				builder.WriteString(s[i:])
				break
			}
		}
		if !ansiFound {
			builder.WriteString(s[i:])
			break
		}
	}
	return builder.String()
}

func NewANSIEscapeSequenceStripper() *ANSIEscapeSequenceStripper {
	return &ANSIEscapeSequenceStripper{}
}

var _ SpecialSequenceStripper = (*ANSIEscapeSequenceStripper)(nil)
