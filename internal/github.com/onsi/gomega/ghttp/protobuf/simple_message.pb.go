// Code generated by protoc-gen-go.
// source: simple_message.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	simple_message.proto

It has these top-level messages:
	SimpleMessage
*/
package protobuf

import proto "github.com/duskhacker/cqrsnu/internal/github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type SimpleMessage struct {
	Description      *string `protobuf:"bytes,1,req,name=description" json:"description,omitempty"`
	Id               *int32  `protobuf:"varint,2,req,name=id" json:"id,omitempty"`
	Metadata         *string `protobuf:"bytes,3,opt,name=metadata" json:"metadata,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SimpleMessage) Reset()         { *m = SimpleMessage{} }
func (m *SimpleMessage) String() string { return proto.CompactTextString(m) }
func (*SimpleMessage) ProtoMessage()    {}

func (m *SimpleMessage) GetDescription() string {
	if m != nil && m.Description != nil {
		return *m.Description
	}
	return ""
}

func (m *SimpleMessage) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *SimpleMessage) GetMetadata() string {
	if m != nil && m.Metadata != nil {
		return *m.Metadata
	}
	return ""
}
