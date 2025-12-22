// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// MarshalJSON implements json.Marshaler for Hex.
func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ConciseString())
}

// UnmarshalJSON implements json.Unmarshaler for Hex.
func (h *Hex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := scanConciseString(s)
	if err != nil {
		return err
	}
	*h = parsed
	return nil
}

// Value implements driver.Valuer for database storage.
func (h Hex) Value() (driver.Value, error) {
	return h.ConciseString(), nil
}

// Scan implements sql.Scanner for database retrieval.
func (h *Hex) Scan(src any) (err error) {
	if s, ok := src.(string); ok {
		*h, err = scanConciseString(s)
		return err
	}
	return fmt.Errorf("unknown type")
}

func scanConciseString(s string) (h Hex, err error) {
	qStart, rStart, sStart := -1, -1, -1
	for i := 0; i < len(s); i = i + 1 {
		if s[i] == '-' || s[i] == '+' {
			if qStart == -1 {
				qStart = i
			} else if rStart == -1 {
				rStart = i
			} else {
				sStart = i
			}
		}
	}
	if sStart == -1 {
		return Hex{}, fmt.Errorf("invalid hexg")
	}
	h.q, err = strconv.Atoi(s[qStart:rStart])
	if err != nil {
		return Hex{}, err
	}
	h.r, err = strconv.Atoi(s[rStart:sStart])
	if err != nil {
		return Hex{}, err
	}
	h.s, err = strconv.Atoi(s[sStart:])
	if err != nil {
		return Hex{}, err
	}
	if h.q+h.r+h.s != 0 {
		return Hex{}, fmt.Errorf("invalid hexg")
	}
	return h, nil
}
