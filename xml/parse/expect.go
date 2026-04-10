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

package parse

import (
	"fmt"
	"io"
	"testing"

	"github.com/katydid/parser-go/cast"
	"github.com/katydid/parser-go/parse"
)

func expectHint(t *testing.T, p Parser, want Hint) {
	t.Helper()
	got, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectFalse(t *testing.T, tzer Parser) {
	t.Helper()
	tokenKind, _, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.FalseKind {
		t.Fatalf("expected false, but got %v", tokenKind)
	}
}

func expectTrue(t *testing.T, tzer Parser) {
	t.Helper()
	tokenKind, _, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.TrueKind {
		t.Fatalf("expected true, but got %v", tokenKind)
	}
}

func expectInt(t *testing.T, tzer Parser, want int64) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.Int64Kind {
		t.Fatalf("expected int64, but got %v", tokenKind)
	}
	got := cast.ToInt64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectFloat(t *testing.T, tzer Parser, want float64) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.Float64Kind {
		t.Fatalf("expected float64, but got %v", tokenKind)
	}
	got := cast.ToFloat64(gotb)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectString(t *testing.T, tzer Parser, want string) {
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

func expectTag(t *testing.T, tzer Parser, want string) {
	t.Helper()
	tokenKind, gotb, err := tzer.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tokenKind != parse.TagKind {
		t.Fatalf("expected tag, but got %v", tokenKind)
	}
	gotf := string(gotb)
	got := fmt.Sprintf("%v", gotf)
	if got != want {
		t.Fatalf("want %v, but got %v", want, got)
	}
}

func expectNoErr(t *testing.T, f func() error) {
	t.Helper()
	err := f()
	if err != nil {
		t.Fatal(err)
	}
}

func expectErr[A any](t *testing.T, f func() (A, error)) {
	t.Helper()
	got, err := f()
	if err == nil {
		t.Fatalf("expected error, but got %v", got)
	}
}

func expectEOF(t *testing.T, p Parser) {
	t.Helper()
	if h, err := p.Next(); err != io.EOF {
		t.Fatalf("expected EOF, but got %v with hint %v", err, h)
	}
}
