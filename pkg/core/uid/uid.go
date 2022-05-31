// Copyright (c) 2021 by library authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uid

import (
	"fmt"
)

// UID is a general type which all UID shall be based upon.
// Specified in TCG Storage Architecture Core Specification Version 2.01 - Rev 1.0
type UID [8]byte

type RowUID UID

type InvokingID UID

type SPID UID

type AuthorityObjectUID UID

var (
	InvokeIDNull   = InvokingID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	InvokeIDThisSP = InvokingID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	InvokeIDSMU    = InvokingID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF}
)

var (
	LockingAuthorityBandMaster0 = AuthorityObjectUID{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x80, 0x01}
	LockingAuthorityAdmin1      = AuthorityObjectUID{0x00, 0x00, 0x00, 0x09, 0x00, 0x01, 0x00, 0x01}
	LockingAuthorityAdmin2      = AuthorityObjectUID{0x00, 0x00, 0x00, 0x09, 0x00, 0x01, 0x00, 0x02}
	AuthorityAnybody            = AuthorityObjectUID{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x01}
	AuthoritySID                = AuthorityObjectUID{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x06}
	AuthorityPSID               = AuthorityObjectUID{0x00, 0x00, 0x00, 0x09, 0x00, 0x01, 0xFF, 0x01} // Opal Feature Set: PSID
)

var (
	GlobalRangeRowUID RowUID = [8]byte{0x00, 0x00, 0x08, 0x02, 0x00, 0x00, 0x00, 0x01}
)

var (
	AdminSP             = SPID{0x00, 0x00, 0x02, 0x05, 0x00, 0x00, 0x00, 0x01}
	LockingSP           = SPID{0x00, 0x00, 0x02, 0x05, 0x00, 0x00, 0x00, 0x02}
	EnterpriseLockingSP = SPID{0x00, 0x00, 0x02, 0x05, 0x00, 0x01, 0x00, 0x01} // Enterprise SSC
)

var (
	// 5.3.2.12 Credential Table Group - C_PIN (Object Table)
	CPIN_SID  InvokingID = [8]byte{0x00, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x00, 0x01}
	CPIN_MSID InvokingID = [8]byte{0x00, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x54, 0x02}
)

func (iid InvokingID) String() string {
	switch iid {
	case InvokeIDNull:
		return "InvokeIDNull"
	case InvokeIDThisSP:
		return "InvokeIDThisSP"
	case InvokeIDSMU:
		return "InvokeIDSMU"
	case InvokingID(AdminSP):
		return "AdminSP"
	case InvokingID(LockingSP):
		return "LockingSP"
	case InvokingID(EnterpriseLockingSP):
		return "EnterpriseLockingSP"
	case InvokingID(AuthorityPSID):
		return "AuthorityPSID"
	case InvokingID(AuthorityAnybody):
		return "AuthorityAnybody"
	case InvokingID(AuthoritySID):
		return "AuthoritySID"
	case InvokingID(CPIN_SID):
		return "CPIN_SID"
	case InvokingID(CPIN_MSID):
		return "CPIN_MSID"
	}
	return fmt.Sprintf("%v", iid)
}
