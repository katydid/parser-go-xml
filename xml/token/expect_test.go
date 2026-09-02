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

package token

import (
	"fmt"
	"io"
	"testing"

	"katydid.org.za/go/parser-go/cp"
	"katydid.org.za/go/parser-go/parse"
)

func expect[A comparable](t *testing.T, f func() (A, error), want A) {
	t.Helper()
	got, err := f()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectFalse(t *testing.T, tzer Tokenizer) {
	t.Helper()
	tokenKind, _, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.FalseKind {
		t.Fatalf("expected false, but got %v", tokenKind)
	}
}

func expectTrue(t *testing.T, tzer Tokenizer) {
	t.Helper()
	tokenKind, _, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.TrueKind {
		t.Fatalf("expected true, but got %v", tokenKind)
	}
}

func expectInt(t *testing.T, tzer Tokenizer, want int64) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.Int64Kind {
		t.Fatalf("expected int64, but got %v", tokenKind)
	}
	got := cp.ToInt64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectFloat(t *testing.T, tzer Tokenizer, want float64) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.Float64Kind {
		t.Fatalf("expected float64, but got %v", tokenKind)
	}
	got := cp.ToFloat64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectStr(t *testing.T, tzer Tokenizer, want string) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.StringKind {
		t.Fatalf("expected string, but got %v", tokenKind)
	}
	gotf := string(gotb)
	got := fmt.Sprintf("%v", gotf)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectErr[A any](t *testing.T, f func() (A, error)) {
	t.Helper()
	got, err := f()
	if err == nil {
		t.Fatalf("expected error, but got %v", got)
	}
}

func expectEOF[A any](t *testing.T, f func() (A, error)) {
	t.Helper()
	_, err := f()
	if err != io.EOF {
		t.Fatalf("expected EOF, but got err = %v", err)
	}
}
