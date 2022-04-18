// Copyright (c) 2021 by library authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Implements TCG Storage Core Method calling

package method

import (
	"bytes"
	"errors"

	//"fmt"
	//"log"

	"github.com/open-source-firmware/go-tcg-storage/pkg/core/stream"
	"github.com/open-source-firmware/go-tcg-storage/pkg/core/uid"
)

type MethodFlag int

const (
	MethodFlagOptionalAsName MethodFlag = 1
)

var (
	ErrMalformedMethodResponse    = errors.New("method response was malformed")
	ErrEmptyMethodResponse        = errors.New("method response was empty")
	ErrMethodListUnbalanced       = errors.New("method argument list is unbalanced")
	ErrTPerClosedSession          = errors.New("TPer forcefully closed our session")
	ErrReceivedUnexpectedResponse = errors.New("method response was unexpected")
	ErrMethodTimeout              = errors.New("method call timed out waiting for a response")

	MethodStatusSuccess uint = 0x00
	MethodStatusCodeMap      = map[uint]error{
		0x00: errors.New("method returned status SUCCESS"),
		0x01: errors.New("method returned status NOT_AUTHORIZED"),
		0x02: errors.New("method returned status OBSOLETE"),
		0x03: errors.New("method returned status SP_BUSY"),
		0x04: errors.New("method returned status SP_FAILED"),
		0x05: errors.New("method returned status SP_DISABLED"),
		0x06: errors.New("method returned status SP_FROZEN"),
		0x07: errors.New("method returned status NO_SESSIONS_AVAILABLE"),
		0x08: errors.New("method returned status UNIQUENESS_CONFLICT"),
		0x09: errors.New("method returned status INSUFFICIENT_SPACE"),
		0x0A: errors.New("method returned status INSUFFICIENT_ROWS"),
		0x0B: errors.New("method returned status INVALID_COMMAND"), /* from Core Revision 0.9 Draft */
		0x0C: errors.New("method returned status INVALID_PARAMETER"),
		0x0D: errors.New("method returned status INVALID_REFERENCE"),         /* from Core Revision 0.9 Draft */
		0x0E: errors.New("method returned status INVALID_SECMSG_PROPERTIES"), /* from Core Revision 0.9 Draft */
		0x0F: errors.New("method returned status TPER_MALFUNCTION"),
		0x10: errors.New("method returned status TRANSACTION_FAILURE"),
		0x11: errors.New("method returned status RESPONSE_OVERFLOW"),
		0x12: errors.New("method returned status AUTHORITY_LOCKED_OUT"),
		0x3F: errors.New("method returned status FAIL"),
	}

	ErrMethodStatusNotAuthorized       = MethodStatusCodeMap[0x01]
	ErrMethodStatusSPBusy              = MethodStatusCodeMap[0x03]
	ErrMethodStatusNoSessionsAvailable = MethodStatusCodeMap[0x07]
	ErrMethodStatusInvalidParameter    = MethodStatusCodeMap[0x0C]
	ErrMethodStatusAuthorityLockedOut  = MethodStatusCodeMap[0x12]
)

type Call interface {
	MarshalBinary() ([]byte, error)
	IsEOS() bool
}

type MethodCall struct {
	iid uid.InvokingID
	mid uid.MethodID
	buf bytes.Buffer
	// Used to verify detect programming errors
	depth int
	flags MethodFlag
}

// Prepare a new method call
func NewMethodCall(iid uid.InvokingID, mid uid.MethodID, flags MethodFlag) *MethodCall {
	m := &MethodCall{iid, mid, bytes.Buffer{}, 0, flags}
	m.buf.Write(stream.Token(stream.Call))
	m.Bytes(iid[:])
	m.Bytes(mid[:])
	// Start argument list
	m.StartList()
	return m
}

// func (m *MethodCall) parseParams(params stream.List) {
// 	for i := 0; i < len(params); i++ {
// 		param := params[i]
// 		tok, ok1 := param.(stream.TokenType)
// 		if ok1 {
// 			if tok == stream.StartName {
// 				_, ok1 := params[i+1].([]byte) // name
// 				_, ok2 := params[i+2].(uint)   // val2
// 				endToc, ok3 := params[i+3].(stream.TokenType)
// 				if ok1 && ok2 && ok3 && endToc == stream.EndName {
// 					//log.Printf("got tokenType %q on params[%d]. name: %s, index: %d", tok.String(), i, string(name), uint(val2))
// 					i += 3
// 				} else {
// 					_, ok1 := params[i+1].(uint)   //val1
// 					_, ok2 := params[i+2].([]byte) //name
// 					endToc, ok3 := params[i+3].(stream.TokenType)
// 					if ok1 && ok2 && ok3 && endToc == stream.EndName {
// 						//log.Printf("got tokenType %q on params[%d]. index: %d, name: %s", tok.String(), i, uint(val1), string(name))
// 						i += 3
// 					} else {
// 						//log.Printf("got tokenType %q on params[%d]", tok.String(), i)
// 					}
// 				}
// 			}
// 			continue
// 		}
// 		lst, ok2 := param.(stream.List)
// 		if ok2 {
// 			//log.Printf("got list on params[%d]", i)
// 			m.parseParams(lst)
// 			continue
// 		}
// 		_, ok3 := param.(uint)
// 		if ok3 {
// 			//log.Printf("got uint on params[%d]", i)
// 			continue
// 		}
// 		_, ok4 := param.([]byte)
// 		if ok4 {
// 			//log.Printf("got []byte on params[%d]", i)
// 			continue
// 		}
// 		log.Printf("unknown value at params[%d], %v", i, param)
// 	}
// }

// func (m *MethodCall) String() string {
// 	var buff bytes.Buffer
// 	b, err := m.MarshalBinary()
// 	if err != nil {
// 		log.Printf("failed to marshal binary")
// 		return ""
// 	}

// 	req, err := stream.Decode(b)
// 	if err != nil {
// 		log.Printf("failed to decode m")
// 		return ""
// 	}
// 	if len(req) >= 6 {
// 		if tok, ok := req[0].(stream.TokenType); !ok || tok != stream.Call {
// 			log.Printf("expected first token to be Call")
// 			return ""
// 		}
// 		var uidStr uid.InvokingID
// 		var midstr uid.MethodID
// 		if iid, ok := req[1].([]byte); !ok {
// 			log.Printf("expected uid")
// 			return ""
// 		} else {
// 			copy(uidStr[:], iid[:8])
// 		}
// 		if mid, ok := req[2].([]byte); !ok {
// 			return ""
// 		} else {
// 			copy(midstr[:], mid[:8])
// 		}
// 		if params, ok := req[3].(stream.List); ok {
// 			m.parseParams(params)
// 		}

// 		// if tok, ok := req[4].(stream.TokenType); ok {
// 		// 	if lst, ok := req[5].(stream.List); ok {
// 		// 		log.Printf("got tokenType %q: %v", tok.String(), lst)
// 		// 	}
// 		// }
// 		if midstr == uid.MethodIDSMStartSession {
// 			reqBytes := req[3].(stream.List)
// 			hsn := reqBytes[0]
// 			spid := reqBytes[1].([]byte)
// 			readOnly := reqBytes[2].(uint)
// 			// for some reason readOnly is opposite - we set mc.Bool(!s.ReadOnly) so 0 means read-only -- go figure
// 			buff.WriteString(fmt.Sprintf("[MethodID: StartSession] hsn: %d, spid: %v, readOnly: %t", hsn, spid, readOnly == 0))
// 		} else if midstr == uid.OpalGet {
// 			buff.WriteString("[MethodID: Get]")
// 		} else if midstr == uid.OpalNext {
// 			buff.WriteString("[MethodID: Next]")
// 		} else if midstr == uid.OpalSet {
// 			reqBytes := req[3].(stream.List)
// 			_ = reqBytes[0].(stream.TokenType) //startList
// 			_ = reqBytes[1].(uint)
// 			lst := reqBytes[2].(stream.List)
// 			_ = lst[0].(stream.TokenType) //startName
// 			length := lst[1].(uint)       //len
// 			name, ok := lst[2].([]byte)   //data
// 			if !ok {
// 				buff.WriteString(fmt.Sprintf("[MethodID: Set] %v,", req[3]))
// 			}
// 			_ = lst[3].(stream.TokenType) //endName
// 			buff.WriteString(fmt.Sprintf("[MethodID: Set] len: %v, name: %v", length, string(name)))
// 		} else if midstr == uid.OpalAuthenticate {
// 			reqBytes := req[3].(stream.List)
// 			lockingAuthority := reqBytes[0].([]byte)
// 			_ = reqBytes[1].(stream.TokenType) //startName
// 			length := reqBytes[2].(uint)
// 			name := reqBytes[3].([]byte)
// 			_ = reqBytes[4].(stream.TokenType) //endName
// 			buff.WriteString(fmt.Sprintf("[MethodID: Authenticate] authority: %v, len: %d, name: %s", lockingAuthority, length, string(name)))
// 		} else if midstr == uid.OpalRevert {
// 			buff.WriteString("[MethodID: Revert]")
// 		} else if midstr == uid.MethodIDSMProperties {
// 			// var tp TPerProperties
// 			// parseTPerProperties(req[3].(stream.List), &tp)
// 			// buff.WriteString(fmt.Sprintf("[MethodID: Properties] tp: %v", tp))
// 			// reqBytes := req[3].(stream.List)
// 			// for i := 0; i < len(reqBytes); i++ {
// 			// 	tokenType := reqBytes[i].(stream.TokenType)
// 			// 	if tokenType == stream.StartName {
// 			// 		length := reqBytes[i+1].(uint)
// 			// 		name := reqBytes[i+2].(stream.List)
// 			// 		_ = reqBytes[i+3].(stream.TokenType)
// 			// 		buff.WriteString(fmt.Sprintf("[MethodID: Properties] tokenType: %v, length: %d, name: %s", stream.StartName, length, name))
// 			// 		i += 3
// 			// 	}
// 		} else {
// 			buff.WriteString(fmt.Sprintf("default to invoking-id: %s, method-id: %s. req: %+v", uidStr.String(), midstr.String(), req[3]))
// 		}
// 		return buff.String()
// 	}

// 	return fmt.Sprintf("unhandled method call - invoking-id: %s, method-id: %s. req: %+v", m.iid.String(), m.mid.String(), req)
// }

// Copy the current state of a method call into a new independent copy
func (m *MethodCall) Clone() *MethodCall {
	mn := &MethodCall{
		iid:   m.iid,
		mid:   m.mid,
		buf:   bytes.Buffer{},
		depth: m.depth,
		flags: m.flags,
	}
	mn.buf.Write(m.buf.Bytes())
	return mn
}

func (m *MethodCall) IsEOS() bool {
	return false
}

func (m *MethodCall) StartList() {
	m.depth++
	m.buf.Write(stream.Token(stream.StartList))
}

func (m *MethodCall) EndList() {
	m.depth--
	m.buf.Write(stream.Token(stream.EndList))
}

// Start an optional parameters group
//
// From "3.2.1.2 Method Signature Pseudo-code"
// > Optional parameters are submitted to the method invocation as Named value pairs.
// > The Name portion of the Named value pair SHALL be a uinteger. Starting at zero,
// > these uinteger values are assigned based on the ordering of the optional parameters
// > as defined in this document.
// The above is true for Core 2.0 things like OpalV2 but not for e.g. Enterprise.
// Thus, we provide a way for the code to switch between using uint or string.
func (m *MethodCall) StartOptionalParameter(id uint, name string) {
	m.depth++
	m.buf.Write(stream.Token(stream.StartName))
	if m.flags&MethodFlagOptionalAsName > 0 {
		m.buf.Write(stream.Bytes([]byte(name)))
	} else {
		m.buf.Write(stream.UInt(id))
	}
}

// Add a named value (uint) pair
func (m *MethodCall) NamedUInt(name string, val uint) {
	m.buf.Write(stream.Token(stream.StartName))
	m.buf.Write(stream.Bytes([]byte(name)))
	m.buf.Write(stream.UInt(val))
	m.buf.Write(stream.Token(stream.EndName))
}

// Add a named value (bool) pair
func (m *MethodCall) NamedBool(name string, val bool) {
	if val {
		m.NamedUInt(name, 1)
	} else {
		m.NamedUInt(name, 0)
	}
}

// Token adds a specific token to the MethodCall buffer.
func (m *MethodCall) Token(t stream.TokenType) {
	m.buf.Write(stream.Token(t))
}

// EndOptionalParameter ends the current optional parameter group
func (m *MethodCall) EndOptionalParameter() {
	m.depth--
	m.buf.Write(stream.Token(stream.EndName))
}

// Bytes adds a bytes atom
func (m *MethodCall) Bytes(b []byte) {
	m.buf.Write(stream.Bytes(b))
}

// UInt adds an uint atom
func (m *MethodCall) UInt(v uint) {
	m.buf.Write(stream.UInt(v))
}

// Bool adds a bool atom (as uint)
func (m *MethodCall) Bool(v bool) {
	if v {
		m.UInt(1)
	} else {
		m.UInt(0)
	}
}

// Marshal the complete method call to the data stream representation
func (m *MethodCall) MarshalBinary() ([]byte, error) {
	mn := *m
	mn.EndList() // End argument list
	// Finish method call
	mn.buf.Write(stream.Token(stream.EndOfData))
	mn.StartList() // Status code list
	mn.buf.Write(stream.UInt(MethodStatusSuccess))
	mn.buf.Write(stream.UInt(0)) // Reserved
	mn.buf.Write(stream.UInt(0)) // Reserved
	mn.EndList()
	if mn.depth != 0 {
		return nil, ErrMethodListUnbalanced
	}
	return mn.buf.Bytes(), nil
}

type EOSMethodCall struct {
}

func (m *EOSMethodCall) MarshalBinary() ([]byte, error) {
	return stream.Token(stream.EndOfSession), nil
}

func (m *EOSMethodCall) IsEOS() bool {
	return true
}
