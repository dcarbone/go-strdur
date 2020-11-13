package strdur

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"strings"
	"time"
)

// StringDuration is a quick hack to let us use time.Duration strings as config values in hcl files
type StringDuration string

func (sd StringDuration) String() string {
	// todo: this...is awful
	return sd.Duration().String()
}

func (sd *StringDuration) Set(v string) error {
	if v == "" {
		*sd = "0s"
		return nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return err
	}
	*sd = StringDuration(d.String())
	return nil
}

func (sd StringDuration) MarshalBinary() ([]byte, error) {
	b := make([]byte, 8, 8)
	binary.LittleEndian.PutUint64(b, uint64(sd.Duration()))
	return b, nil
}

func (sd *StringDuration) UnmarshalBinary(b []byte) error {
	if l := len(b); l != 8 {
		return fmt.Errorf("expected 8 bytes, saw %d", l)
	}
	uv := binary.LittleEndian.Uint64(b)
	if uv > math.MaxInt64 {
		return fmt.Errorf("int64 overflow: %d", uv)
	}
	*sd = StringDuration(time.Duration(uv).String())
	return nil
}

func (sd StringDuration) GobEncode() ([]byte, error) {
	return sd.MarshalBinary()
}

func (sd *StringDuration) GobDecode(b []byte) error {
	return sd.UnmarshalBinary(b)
}

func (sd StringDuration) MarshalText() ([]byte, error) {
	return []byte(sd.String()), nil
}

func (sd *StringDuration) UnmarshalText(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	return sd.Set(string(b))
}

func (sd StringDuration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + sd.String() + "\""), nil
}

func (sd *StringDuration) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	clean := strings.Trim(string(b), "\"")
	if clean == "" {
		*sd = "0s"
		return nil
	}
	return sd.Set(clean)
}

func (sd StringDuration) Duration() time.Duration {
	td, _ := time.ParseDuration(string(sd))
	return td
}

func (sd *StringDuration) FromDuration(td time.Duration) {
	*sd = StringDuration(td.String())
}

func ConfinatorFlagVarTypeFunc(fs *flag.FlagSet, varPtr interface{}, name, usage string) {
	fs.Var(varPtr.(*StringDuration), name, usage)
}
