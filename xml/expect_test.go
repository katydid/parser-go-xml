//  Copyright 2025 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package xml

import (
	"io"
	"testing"
)

func expectNoErr(t *testing.T, f func() error) {
	t.Helper()
	err := f()
	if err != nil {
		t.Fatal(err)
	}
}

func expect[A comparable](t *testing.T, f func() (A, error), want A) {
	t.Helper()
	got, err := f()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("want %#v, but got %#v", want, got)
	}
}

func expectEOF(t *testing.T, f func() error) {
	t.Helper()
	err := f()
	if err != io.EOF {
		t.Fatalf("expected EOF, but got err = %v", err)
	}
}
