package json

import (
	"encoding/json"
	"strconv"
	"strings"
)

// The JSON Feed spec asks readers to be liberal about a few field types. An id
// "presented as a number or other type" must be coerced to a string, and real
// feeds also send booleans and numeric fields as strings or floats. These
// Unmarshalers keep one off-spec field from failing the whole feed.
//
// Each uses the standard alias trick: a shadow type without methods decodes
// every field normally, while the loose fields are pulled out as raw JSON and
// coerced by hand.

func (i *Item) UnmarshalJSON(data []byte) error {
	type alias Item
	aux := &struct {
		ID json.RawMessage `json:"id"`
		*alias
	}{alias: (*alias)(i)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	i.ID = coerceString(aux.ID)
	return nil
}

func (f *Feed) UnmarshalJSON(data []byte) error {
	type alias Feed
	aux := &struct {
		Expired json.RawMessage `json:"expired"`
		*alias
	}{alias: (*alias)(f)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	f.Expired = coerceBool(aux.Expired)
	return nil
}

func (a *Attachments) UnmarshalJSON(data []byte) error {
	type alias Attachments
	aux := &struct {
		SizeInBytes       json.RawMessage `json:"size_in_bytes"`
		DurationInSeconds json.RawMessage `json:"duration_in_seconds"`
		*alias
	}{alias: (*alias)(a)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.SizeInBytes = coerceInt64(aux.SizeInBytes)
	a.DurationInSeconds = coerceInt64(aux.DurationInSeconds)
	return nil
}

// coerceString accepts a JSON string or number and returns it as a string.
func coerceString(raw json.RawMessage) string {
	if isEmptyJSON(raw) {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		return n.String()
	}
	return strings.TrimSpace(string(raw))
}

// coerceBool accepts a JSON boolean or a "true"/"false" string.
func coerceBool(raw json.RawMessage) bool {
	if isEmptyJSON(raw) {
		return false
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return b
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		b, _ := strconv.ParseBool(strings.TrimSpace(s))
		return b
	}
	return false
}

// coerceInt64 accepts a JSON number (integer or float) or a numeric string.
func coerceInt64(raw json.RawMessage) int64 {
	if isEmptyJSON(raw) {
		return 0
	}
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		if i, err := n.Int64(); err == nil {
			return i
		}
		if fl, err := n.Float64(); err == nil {
			return int64(fl)
		}
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if fl, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
			return int64(fl)
		}
	}
	return 0
}

func isEmptyJSON(raw json.RawMessage) bool {
	s := strings.TrimSpace(string(raw))
	return s == "" || s == "null"
}
