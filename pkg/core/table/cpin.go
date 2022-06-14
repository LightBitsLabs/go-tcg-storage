// // Copyright (c) 2021 by library authors. All rights reserved.
// // Use of this source code is governed by a BSD-style
// // license that can be found in the LICENSE file.

// // Implements TCG Storage Core Table operations on Locking SP tables

package table

import (
	"github.com/open-source-firmware/go-tcg-storage/pkg/core"
	"github.com/open-source-firmware/go-tcg-storage/pkg/core/method"
	"github.com/open-source-firmware/go-tcg-storage/pkg/core/uid"
)

type CPINInfoRow struct {
	// 5.3.2.12.1 This is the unique identifier of this row in the C_PIN table.
	UID uid.RowUID
	// 5.3.2.12.2 This is the name of the C_PIN object.
	// For C_PIN objects that exist at issuance, this column SHALL NOT be modifiable by the host.
	Name *string
	// 5.3.2.12.3 This is a name that MAY be shared by multiple C_PIN objects.
	// For C_PIN objects that exist at issuance, this column SHALL NOT be modifiable by the host.
	CommonName *string
	// 5.3.2.12.4 This is the bytes value to which authentication attempts on this C_PIN object are matched.
	Password []byte
	// 5.3.2.12.5 This is a reference to the byte table that holds the character set used for TPer-generated PIN column
	// values created using the GenKey method.
	// If the value of this column is a NULL UID reference, then the default character set is used with the
	// GenKey method.
	CharSet []byte
	// 5.3.2.12.6 The value of this column is the maximum number of failed authentication attempts that are able to be
	// made using this C_PIN object.
	// The default value of the TryLimit column when a new C_PIN object is created is 0. The value 0 in this
	// column indicates that there is no limit on the number of authentication attempts for that object.
	TryLimit *uint32
	// 5.3.2.12.7 This column identifies the current number of failed authentication attempts using this C_PIN object.
	Tries *uint32
	// 5.3.2.12.8 The value of this column identifies if value of Tries column is persistent through power cycles.
	Persistence *bool
}

func CPINInfo(s *core.Session) (*CPINInfoRow, error) {
	rowUID := uid.RowUID{}
	copy(rowUID[:], uid.Admin_C_PIN_SIDRow[:])

	val, err := GetFullRow(s, rowUID)
	if err != nil {
		return nil, err
	}
	row := CPINInfoRow{}
	for col, val := range val {
		switch col {
		case "0", "UID":
			v, ok := val.([]byte)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			copy(row.UID[:], v[:8])
		case "1", "Name":
			v, ok := val.([]byte)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			vv := string(v)
			row.Name = &vv
		case "2", "CommonName":
			v, ok := val.([]byte)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			vv := string(v)
			row.CommonName = &vv
		case "3", "Password":
			v, ok := val.([]byte)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			vv := v
			row.Password = vv
		case "4", "CharSet":
			v, ok := val.([]uint8)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			vv := v
			row.CharSet = vv
		case "5", "TryLimit":
			v, ok := val.(uint)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			vv := uint32(v)
			row.TryLimit = &vv
		case "6", "Tries":
			v, ok := val.(uint)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			vv := uint32(v)
			row.Tries = &vv
		case "7", "Persistence":
			v, ok := val.(uint)
			if !ok {
				return nil, method.ErrMalformedMethodResponse
			}
			var vv bool
			if v > 0 {
				vv = true
			}
			row.Persistence = &vv
		}
	}
	return &row, nil
}
