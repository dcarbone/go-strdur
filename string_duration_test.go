package strdur_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/dcarbone/go-strdur"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func TestStringDuration(t *testing.T) {
	type (
		// represents a set of inputs and their expected outputs
		testValue struct {
			iString string
			iText   []byte
			iJSON   []byte
			iBinary []byte
			iHCL    []byte

			oString string
			oText   []byte
			oJSON   []byte
			oBinary []byte
		}

		testRun struct {
			name string
			fn   func(*testing.T, *strdur.StringDuration, testValue) error
		}
	)

	var testValues = []testValue{
		{
			iString: "",
			iText:   []byte(""),
			iJSON:   []byte("\"\""),
			iBinary: []byte{0, 0, 0, 0, 0, 0, 0, 0},
			iHCL:    []byte("sdv = \"\""),

			oString: "0s",
			oText:   []byte("0s"),
			oJSON:   []byte("\"0s\""),
			oBinary: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			iString: "24h",
			iText:   []byte("24h"),
			iJSON:   []byte("\"24h\""),
			iBinary: []byte{0, 0, 79, 145, 148, 78, 0, 0},
			iHCL:    []byte("sdv = \"24h\""),

			oString: "24h0m0s",
			oText:   []byte("24h0m0s"),
			oJSON:   []byte("\"24h0m0s\""),
			oBinary: []byte{0, 0, 79, 145, 148, 78, 0, 0},
		},
		{
			iString: "2562047h47m16s854775807ns",
			iText:   []byte("2562047h47m16s854775807ns"),
			iJSON:   []byte("\"2562047h47m16s854775807ns\""),
			iBinary: []byte{255, 255, 255, 255, 255, 255, 255, 127},
			iHCL:    []byte("sdv = \"2562047h47m16s854775807ns\""),

			oString: "2562047h47m16.854775807s",
			oText:   []byte("2562047h47m16.854775807s"),
			oJSON:   []byte("\"2562047h47m16.854775807s\""),
			oBinary: []byte{255, 255, 255, 255, 255, 255, 255, 127},
		},
	}

	var testDefinitions = []testRun{
		{
			name: "flag",
			fn: func(t *testing.T, sdv *strdur.StringDuration, tv testValue) error {
				fs := flag.NewFlagSet("StringDurationTest", flag.ContinueOnError)
				fs.Var(sdv, "sdv", "")
				if err := fs.Parse([]string{"-sdv", tv.iString}); err != nil {
					return fmt.Errorf("error parsing %q as %T: %w", tv.iString, sdv, err)
				}
				return nil
			},
		},
		{
			name: "text",
			fn: func(t *testing.T, sdv *strdur.StringDuration, tv testValue) error {
				if err := sdv.UnmarshalText(tv.iText); err != nil {
					return fmt.Errorf("error paring %q as %T: %w", tv.iText, sdv, err)
				}
				return nil
			},
		},
		{
			name: "json",
			fn: func(t *testing.T, sdv *strdur.StringDuration, tv testValue) error {
				if err := json.Unmarshal(tv.iJSON, sdv); err != nil {
					return fmt.Errorf("error parsing %q: %v", string(tv.iJSON), err)
				}
				return nil
			},
		},
		{
			name: "binary",
			fn: func(t *testing.T, sdv *strdur.StringDuration, tv testValue) error {
				if err := sdv.UnmarshalBinary(tv.iBinary); err != nil {
					return fmt.Errorf("error parsing %q as %T: %w", tv.iBinary, sdv, err)
				}
				return nil
			},
		},
		{
			name: "hcl",
			fn: func(t *testing.T, sdv *strdur.StringDuration, tv testValue) error {
				type confT struct {
					SDV *strdur.StringDuration `hcl:"sdv"`
				}
				if err := hclsimple.Decode("example.hcl", tv.iHCL, nil, &confT{SDV: sdv}); err != nil {
					return fmt.Errorf("error parsing %q as %T: %w", tv.iHCL, sdv, err)
				}
				return nil
			},
		},
	}

	for _, testValue := range testValues {
		localValue := testValue
		for _, testDefinition := range testDefinitions {
			localTestDefinition := testDefinition
			t.Run(testDefinition.name, func(t *testing.T) {
				t.Parallel()
				var sdv strdur.StringDuration
				if err := localTestDefinition.fn(t, &sdv, localValue); err != nil {
					t.Errorf(err.Error())
					t.Fail()
					return
				}
				if str := sdv.String(); str != localValue.oString {
					t.Errorf("String value mismatch: expected=%q; actual=%q", localValue.oString, str)
					t.Fail()
				}
				if txt, err := sdv.MarshalText(); err != nil {
					t.Errorf("Text marshalling error: %v", err)
					t.Fail()
				} else if !bytes.Equal(localValue.oText, txt) {
					t.Errorf("Text value mismatch: expected=%v; actual=%v", localValue.oText, txt)
					t.Fail()
				}
				if jsn, err := sdv.MarshalJSON(); err != nil {
					t.Errorf("JSON marshalling error: %v", err)
					t.Fail()
				} else if !bytes.Equal(localValue.oJSON, jsn) {
					t.Errorf("JSON value mismatch: expected=%v; actual=%v", localValue.oJSON, jsn)
					t.Fail()
				}
				if bn, err := sdv.MarshalBinary(); err != nil {
					t.Errorf("Binary marshalling error: %v", err)
					t.Fail()
				} else if !bytes.Equal(localValue.oBinary, bn) {
					t.Errorf("Binary value mismatch: expected=%v; actual=%v", localValue.oBinary, bn)
					t.Fail()
				}
			})
		}
	}
}
