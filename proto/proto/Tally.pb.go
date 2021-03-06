// Code generated by protoc-gen-go.
// source: proto/Tally.proto
// DO NOT EDIT!

/*
Package Tally is a generated protocol buffer package.

It is generated from these files:
	proto/Tally.proto

It has these top-level messages:
	Content
	Location
	Attribute
	Event
*/
package Tally

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Content struct {
	Mime             *string `protobuf:"bytes,1,req,name=mime" json:"mime,omitempty"`
	Data             []byte  `protobuf:"bytes,2,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Content) Reset()         { *m = Content{} }
func (m *Content) String() string { return proto.CompactTextString(m) }
func (*Content) ProtoMessage()    {}

func (m *Content) GetMime() string {
	if m != nil && m.Mime != nil {
		return *m.Mime
	}
	return ""
}

func (m *Content) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// To be transformed to GeoJson - ex)  {"location" : [-71.34, 41.12]}
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/mapping-geo-point-type.html
type Location struct {
	Lon              *float64 `protobuf:"fixed64,1,req,name=lon" json:"lon,omitempty"`
	Lat              *float64 `protobuf:"fixed64,2,req,name=lat" json:"lat,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Location) Reset()         { *m = Location{} }
func (m *Location) String() string { return proto.CompactTextString(m) }
func (*Location) ProtoMessage()    {}

func (m *Location) GetLon() float64 {
	if m != nil && m.Lon != nil {
		return *m.Lon
	}
	return 0
}

func (m *Location) GetLat() float64 {
	if m != nil && m.Lat != nil {
		return *m.Lat
	}
	return 0
}

type Attribute struct {
	Key              *string  `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	StringValue      *string  `protobuf:"bytes,2,opt,name=string_value" json:"string_value,omitempty"`
	IntValue         *int64   `protobuf:"varint,3,opt,name=int_value" json:"int_value,omitempty"`
	DoubleValue      *float64 `protobuf:"fixed64,4,opt,name=double_value" json:"double_value,omitempty"`
	BoolValue        *bool    `protobuf:"varint,5,opt,name=bool_value" json:"bool_value,omitempty"`
	ContentValue     *Content `protobuf:"bytes,6,opt,name=content_value" json:"content_value,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Attribute) Reset()         { *m = Attribute{} }
func (m *Attribute) String() string { return proto.CompactTextString(m) }
func (*Attribute) ProtoMessage()    {}

func (m *Attribute) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Attribute) GetStringValue() string {
	if m != nil && m.StringValue != nil {
		return *m.StringValue
	}
	return ""
}

func (m *Attribute) GetIntValue() int64 {
	if m != nil && m.IntValue != nil {
		return *m.IntValue
	}
	return 0
}

func (m *Attribute) GetDoubleValue() float64 {
	if m != nil && m.DoubleValue != nil {
		return *m.DoubleValue
	}
	return 0
}

func (m *Attribute) GetBoolValue() bool {
	if m != nil && m.BoolValue != nil {
		return *m.BoolValue
	}
	return false
}

func (m *Attribute) GetContentValue() *Content {
	if m != nil {
		return m.ContentValue
	}
	return nil
}

type Event struct {
	Timestamp        *float64     `protobuf:"fixed64,1,req,name=timestamp" json:"timestamp,omitempty"`
	Type             *string      `protobuf:"bytes,2,req,name=type" json:"type,omitempty"`
	Source           *string      `protobuf:"bytes,3,req,name=source" json:"source,omitempty"`
	Context          *string      `protobuf:"bytes,4,opt,name=context" json:"context,omitempty"`
	Location         *Location    `protobuf:"bytes,5,opt,name=location" json:"location,omitempty"`
	Attributes       []*Attribute `protobuf:"bytes,6,rep,name=attributes" json:"attributes,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}

func (m *Event) GetTimestamp() float64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *Event) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func (m *Event) GetSource() string {
	if m != nil && m.Source != nil {
		return *m.Source
	}
	return ""
}

func (m *Event) GetContext() string {
	if m != nil && m.Context != nil {
		return *m.Context
	}
	return ""
}

func (m *Event) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *Event) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func init() {
}
