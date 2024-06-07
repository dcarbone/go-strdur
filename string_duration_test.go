package strdur_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/mitchellh/mapstructure"

	"github.com/dcarbone/go-strdur/v2"
)

func TestStringDuration(t *testing.T) {
	type (
		// represents a set of inputs and their expected outputs
		testValue struct {
			iString       string
			iText         []byte
			iJSON         []byte
			iBinary       []byte
			iHCL          []byte
			iMapStructure map[string]any

			oString string
			oText   []byte
			oJSON   []byte
			oBinary []byte
		}

		testRun struct {
			name string
			op   func(*testing.T, testValue) (strdur.StringDuration, error)
		}
	)

	var testValues = []testValue{
		{
			iString:       "",
			iText:         []byte(""),
			iJSON:         []byte(`""`),
			iBinary:       []byte{0, 0, 0, 0, 0, 0, 0, 0},
			iHCL:          []byte(`sdv = ""`),
			iMapStructure: map[string]any{"sdv": ""},

			oString: "0s",
			oText:   []byte("0s"),
			oJSON:   []byte(`"0s"`),
			oBinary: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			iString:       "24h",
			iText:         []byte("24h"),
			iJSON:         []byte(`"24h"`),
			iBinary:       []byte{0, 0, 79, 145, 148, 78, 0, 0},
			iHCL:          []byte(`sdv = "24h"`),
			iMapStructure: map[string]any{"sdv": "24h"},

			oString: "24h0m0s",
			oText:   []byte("24h0m0s"),
			oJSON:   []byte(`"24h0m0s"`),
			oBinary: []byte{0, 0, 79, 145, 148, 78, 0, 0},
		},
		{
			iString:       "2562047h47m16s854775807ns",
			iText:         []byte("2562047h47m16s854775807ns"),
			iJSON:         []byte(`"2562047h47m16s854775807ns"`),
			iBinary:       []byte{255, 255, 255, 255, 255, 255, 255, 127},
			iHCL:          []byte(`sdv = "2562047h47m16s854775807ns"`),
			iMapStructure: map[string]any{"sdv": "2562047h47m16s854775807ns"},

			oString: "2562047h47m16.854775807s",
			oText:   []byte("2562047h47m16.854775807s"),
			oJSON:   []byte(`"2562047h47m16.854775807s"`),
			oBinary: []byte{255, 255, 255, 255, 255, 255, 255, 127},
		},
	}

	var testRuns = []testRun{
		{
			name: "flag",
			op: func(t *testing.T, tv testValue) (strdur.StringDuration, error) {
				var sdv strdur.StringDuration
				fs := flag.NewFlagSet("StringDurationTest", flag.ContinueOnError)
				fs.Var(&sdv, "sdv", "")
				if err := fs.Parse([]string{"-sdv", tv.iString}); err != nil {
					return sdv, fmt.Errorf("error parsing %q as %T: %w", tv.iString, sdv, err)
				}
				return sdv, nil
			},
		},
		{
			name: "text",
			op: func(t *testing.T, tv testValue) (strdur.StringDuration, error) {
				var sdv strdur.StringDuration
				if err := sdv.UnmarshalText(tv.iText); err != nil {
					return sdv, fmt.Errorf("error paring %q as %T: %w", tv.iText, sdv, err)
				}
				return sdv, nil
			},
		},
		{
			name: "json",
			op: func(t *testing.T, tv testValue) (strdur.StringDuration, error) {
				var sdv strdur.StringDuration
				if err := json.Unmarshal(tv.iJSON, &sdv); err != nil {
					return sdv, fmt.Errorf("error parsing %q: %v", string(tv.iJSON), err)
				}
				return sdv, nil
			},
		},
		{
			name: "json-field",
			op: func(t *testing.T, tv testValue) (strdur.StringDuration, error) {
				type testT struct {
					SDV strdur.StringDuration `json:"strdur"`
				}
				var (
					tt testT

					tj = []byte(fmt.Sprintf(`{"strdur": %s}`, string(tv.iJSON)))
				)
				if err := json.Unmarshal(tj, &tt); err != nil {
					return tt.SDV, fmt.Errorf("error parsing %q: %v", string(tv.iJSON), err)
				}
				return tt.SDV, nil
			},
		},
		{
			name: "binary",
			op: func(t *testing.T, tv testValue) (strdur.StringDuration, error) {
				var sdv strdur.StringDuration
				if err := sdv.UnmarshalBinary(tv.iBinary); err != nil {
					return sdv, fmt.Errorf("error parsing %q as %T: %w", tv.iBinary, sdv, err)
				}
				return sdv, nil
			},
		},
		{
			name: "hcl",
			op: func(t *testing.T, tv testValue) (strdur.StringDuration, error) {
				type testT struct {
					SDV strdur.StringDuration `hcl:"sdv"`
				}
				var tt testT
				if err := hclsimple.Decode("example.hcl", tv.iHCL, nil, &tt); err != nil {
					return tt.SDV, fmt.Errorf("error parsing %q as %T: %w", tv.iHCL, tt.SDV, err)
				}
				return tt.SDV, nil
			},
		},
		{
			name: "mapstructure",
			op: func(t *testing.T, tv testValue) (strdur.StringDuration, error) {
				type testT struct {
					SDV strdur.StringDuration `mapstructure:"sdv"`
				}
				var tt testT
				if err := mapstructure.Decode(tv.iMapStructure, &tt); err != nil {
					return tt.SDV, fmt.Errorf("error parsing %q as %T: %w", tv.iMapStructure, tt.SDV, err)
				}
				return tt.SDV, nil
			},
		},
	}

	for _, value := range testValues {
		value := value
		for _, testDef := range testRuns {
			testDef := testDef
			t.Run(testDef.name, func(t *testing.T) {
				t.Parallel()
				var (
					sdv strdur.StringDuration
					err error
				)
				if sdv, err = testDef.op(t, value); err != nil {
					t.Errorf(err.Error())
					t.Fail()
					return
				}
				if str := sdv.String(); str != value.oString {
					t.Errorf("String value mismatch: expected=%q; actual=%q", value.oString, str)
					t.Fail()
				}
				if txt, err := sdv.MarshalText(); err != nil {
					t.Errorf("Text marshalling error: %v", err)
					t.Fail()
				} else if !bytes.Equal(value.oText, txt) {
					t.Errorf("Text value mismatch: expected=%v; actual=%v", value.oText, txt)
					t.Fail()
				}
				if jsn, err := sdv.MarshalJSON(); err != nil {
					t.Errorf("JSON marshalling error: %v", err)
					t.Fail()
				} else if !bytes.Equal(value.oJSON, jsn) {
					t.Errorf("JSON value mismatch: expected=%v; actual=%v", value.oJSON, jsn)
					t.Fail()
				}
				if bn, err := sdv.MarshalBinary(); err != nil {
					t.Errorf("Binary marshalling error: %v", err)
					t.Fail()
				} else if !bytes.Equal(value.oBinary, bn) {
					t.Errorf("Binary value mismatch: expected=%v; actual=%v", value.oBinary, bn)
					t.Fail()
				}
			})
		}
	}
}
