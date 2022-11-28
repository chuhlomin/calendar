// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: request.proto

package main

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Size     string `protobuf:"bytes,1,opt,name=size,proto3" json:"size,omitempty"`
	FirstDay int32  `protobuf:"varint,2,opt,name=firstDay,proto3" json:"firstDay,omitempty"`
	Year     int32  `protobuf:"varint,3,opt,name=year,proto3" json:"year,omitempty"`
	Month    int32  `protobuf:"varint,4,opt,name=month,proto3" json:"month,omitempty"`
	Language string `protobuf:"bytes,5,opt,name=language,proto3" json:"language,omitempty"`
	// Days
	DaysFontSize     int32  `protobuf:"varint,10,opt,name=daysFontSize,proto3" json:"daysFontSize,omitempty"`
	DaysFontFamily   string `protobuf:"bytes,11,opt,name=daysFontFamily,proto3" json:"daysFontFamily,omitempty"`
	TextColor        string `protobuf:"bytes,12,opt,name=textColor,proto3" json:"textColor,omitempty"`
	WeekendColor     string `protobuf:"bytes,13,opt,name=weekendColor,proto3" json:"weekendColor,omitempty"`
	DaysX            int32  `protobuf:"varint,14,opt,name=daysX,proto3" json:"daysX,omitempty"`
	DaysY            int32  `protobuf:"varint,15,opt,name=daysY,proto3" json:"daysY,omitempty"`
	DaysXStep        int32  `protobuf:"varint,16,opt,name=daysXStep,proto3" json:"daysXStep,omitempty"`
	DaysYStep        int32  `protobuf:"varint,17,opt,name=daysYStep,proto3" json:"daysYStep,omitempty"`
	ShowInactiveDays bool   `protobuf:"varint,18,opt,name=showInactiveDays,proto3" json:"showInactiveDays,omitempty"`
	InactiveColor    string `protobuf:"bytes,19,opt,name=inactiveColor,proto3" json:"inactiveColor,omitempty"`
	// Month
	ShowMonth       bool   `protobuf:"varint,100,opt,name=showMonth,proto3" json:"showMonth,omitempty"`
	MonthFontFamily string `protobuf:"bytes,101,opt,name=monthFontFamily,proto3" json:"monthFontFamily,omitempty"`
	MonthFontSize   int32  `protobuf:"varint,102,opt,name=monthFontSize,proto3" json:"monthFontSize,omitempty"`
	MonthColor      string `protobuf:"bytes,103,opt,name=monthColor,proto3" json:"monthColor,omitempty"`
	MonthY          int32  `protobuf:"varint,104,opt,name=monthY,proto3" json:"monthY,omitempty"`
	MonthFormat     string `protobuf:"bytes,105,opt,name=monthFormat,proto3" json:"monthFormat,omitempty"`
	// Weekdays
	ShowWeekdays       bool   `protobuf:"varint,200,opt,name=showWeekdays,proto3" json:"showWeekdays,omitempty"`
	WeekdaysFontFamily string `protobuf:"bytes,201,opt,name=weekdaysFontFamily,proto3" json:"weekdaysFontFamily,omitempty"`
	WeekdaysFontSize   int32  `protobuf:"varint,202,opt,name=weekdaysFontSize,proto3" json:"weekdaysFontSize,omitempty"`
	WeekdaysColor      string `protobuf:"bytes,203,opt,name=weekdaysColor,proto3" json:"weekdaysColor,omitempty"`
	WeekdaysX          int32  `protobuf:"varint,204,opt,name=weekdaysX,proto3" json:"weekdaysX,omitempty"`
	WeekdaysY          int32  `protobuf:"varint,205,opt,name=weekdaysY,proto3" json:"weekdaysY,omitempty"`
	// WeekNumbers
	ShowWeekNumbers       bool   `protobuf:"varint,300,opt,name=showWeekNumbers,proto3" json:"showWeekNumbers,omitempty"`
	WeeknumbersFontFamily string `protobuf:"bytes,301,opt,name=weeknumbersFontFamily,proto3" json:"weeknumbersFontFamily,omitempty"`
	WeeknumbersFontSize   int32  `protobuf:"varint,302,opt,name=weeknumbersFontSize,proto3" json:"weeknumbersFontSize,omitempty"`
	WeeknumbersColor      string `protobuf:"bytes,303,opt,name=weeknumbersColor,proto3" json:"weeknumbersColor,omitempty"`
	WeeknumbersX          int32  `protobuf:"varint,304,opt,name=weeknumbersX,proto3" json:"weeknumbersX,omitempty"`
	WeeknumbersY          int32  `protobuf:"varint,305,opt,name=weeknumbersY,proto3" json:"weeknumbersY,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetSize() string {
	if x != nil {
		return x.Size
	}
	return ""
}

func (x *Request) GetFirstDay() int32 {
	if x != nil {
		return x.FirstDay
	}
	return 0
}

func (x *Request) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *Request) GetMonth() int32 {
	if x != nil {
		return x.Month
	}
	return 0
}

func (x *Request) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *Request) GetDaysFontSize() int32 {
	if x != nil {
		return x.DaysFontSize
	}
	return 0
}

func (x *Request) GetDaysFontFamily() string {
	if x != nil {
		return x.DaysFontFamily
	}
	return ""
}

func (x *Request) GetTextColor() string {
	if x != nil {
		return x.TextColor
	}
	return ""
}

func (x *Request) GetWeekendColor() string {
	if x != nil {
		return x.WeekendColor
	}
	return ""
}

func (x *Request) GetDaysX() int32 {
	if x != nil {
		return x.DaysX
	}
	return 0
}

func (x *Request) GetDaysY() int32 {
	if x != nil {
		return x.DaysY
	}
	return 0
}

func (x *Request) GetDaysXStep() int32 {
	if x != nil {
		return x.DaysXStep
	}
	return 0
}

func (x *Request) GetDaysYStep() int32 {
	if x != nil {
		return x.DaysYStep
	}
	return 0
}

func (x *Request) GetShowInactiveDays() bool {
	if x != nil {
		return x.ShowInactiveDays
	}
	return false
}

func (x *Request) GetInactiveColor() string {
	if x != nil {
		return x.InactiveColor
	}
	return ""
}

func (x *Request) GetShowMonth() bool {
	if x != nil {
		return x.ShowMonth
	}
	return false
}

func (x *Request) GetMonthFontFamily() string {
	if x != nil {
		return x.MonthFontFamily
	}
	return ""
}

func (x *Request) GetMonthFontSize() int32 {
	if x != nil {
		return x.MonthFontSize
	}
	return 0
}

func (x *Request) GetMonthColor() string {
	if x != nil {
		return x.MonthColor
	}
	return ""
}

func (x *Request) GetMonthY() int32 {
	if x != nil {
		return x.MonthY
	}
	return 0
}

func (x *Request) GetMonthFormat() string {
	if x != nil {
		return x.MonthFormat
	}
	return ""
}

func (x *Request) GetShowWeekdays() bool {
	if x != nil {
		return x.ShowWeekdays
	}
	return false
}

func (x *Request) GetWeekdaysFontFamily() string {
	if x != nil {
		return x.WeekdaysFontFamily
	}
	return ""
}

func (x *Request) GetWeekdaysFontSize() int32 {
	if x != nil {
		return x.WeekdaysFontSize
	}
	return 0
}

func (x *Request) GetWeekdaysColor() string {
	if x != nil {
		return x.WeekdaysColor
	}
	return ""
}

func (x *Request) GetWeekdaysX() int32 {
	if x != nil {
		return x.WeekdaysX
	}
	return 0
}

func (x *Request) GetWeekdaysY() int32 {
	if x != nil {
		return x.WeekdaysY
	}
	return 0
}

func (x *Request) GetShowWeekNumbers() bool {
	if x != nil {
		return x.ShowWeekNumbers
	}
	return false
}

func (x *Request) GetWeeknumbersFontFamily() string {
	if x != nil {
		return x.WeeknumbersFontFamily
	}
	return ""
}

func (x *Request) GetWeeknumbersFontSize() int32 {
	if x != nil {
		return x.WeeknumbersFontSize
	}
	return 0
}

func (x *Request) GetWeeknumbersColor() string {
	if x != nil {
		return x.WeeknumbersColor
	}
	return ""
}

func (x *Request) GetWeeknumbersX() int32 {
	if x != nil {
		return x.WeeknumbersX
	}
	return 0
}

func (x *Request) GetWeeknumbersY() int32 {
	if x != nil {
		return x.WeeknumbersY
	}
	return 0
}

var File_request_proto protoreflect.FileDescriptor

var file_request_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x83, 0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x72, 0x73, 0x74, 0x44, 0x61,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x66, 0x69, 0x72, 0x73, 0x74, 0x44, 0x61,
	0x79, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x79, 0x65, 0x61, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x61, 0x79, 0x73, 0x46,
	0x6f, 0x6e, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x64,
	0x61, 0x79, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x64,
	0x61, 0x79, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x64, 0x61, 0x79, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x46, 0x61, 0x6d,
	0x69, 0x6c, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x65, 0x78, 0x74, 0x43, 0x6f, 0x6c, 0x6f, 0x72,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x65, 0x78, 0x74, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x12, 0x22, 0x0a, 0x0c, 0x77, 0x65, 0x65, 0x6b, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x77, 0x65, 0x65, 0x6b, 0x65, 0x6e, 0x64,
	0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x61, 0x79, 0x73, 0x58, 0x18, 0x0e,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x64, 0x61, 0x79, 0x73, 0x58, 0x12, 0x14, 0x0a, 0x05, 0x64,
	0x61, 0x79, 0x73, 0x59, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x64, 0x61, 0x79, 0x73,
	0x59, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x61, 0x79, 0x73, 0x58, 0x53, 0x74, 0x65, 0x70, 0x18, 0x10,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x64, 0x61, 0x79, 0x73, 0x58, 0x53, 0x74, 0x65, 0x70, 0x12,
	0x1c, 0x0a, 0x09, 0x64, 0x61, 0x79, 0x73, 0x59, 0x53, 0x74, 0x65, 0x70, 0x18, 0x11, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x64, 0x61, 0x79, 0x73, 0x59, 0x53, 0x74, 0x65, 0x70, 0x12, 0x2a, 0x0a,
	0x10, 0x73, 0x68, 0x6f, 0x77, 0x49, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x44, 0x61, 0x79,
	0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x73, 0x68, 0x6f, 0x77, 0x49, 0x6e, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x44, 0x61, 0x79, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x69, 0x6e, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x69, 0x6e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x77, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x18, 0x64, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x09, 0x73, 0x68, 0x6f, 0x77, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x12, 0x28, 0x0a,
	0x0f, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x46, 0x6f, 0x6e, 0x74, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79,
	0x18, 0x65, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x46, 0x6f, 0x6e,
	0x74, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x6d, 0x6f, 0x6e, 0x74, 0x68,
	0x46, 0x6f, 0x6e, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x66, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d,
	0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x46, 0x6f, 0x6e, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x0a,
	0x0a, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x67, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x16, 0x0a,
	0x06, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x59, 0x18, 0x68, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6d,
	0x6f, 0x6e, 0x74, 0x68, 0x59, 0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x46, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x18, 0x69, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x6f, 0x6e, 0x74,
	0x68, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x23, 0x0a, 0x0c, 0x73, 0x68, 0x6f, 0x77, 0x57,
	0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x18, 0xc8, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x73, 0x68, 0x6f, 0x77, 0x57, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x12, 0x2f, 0x0a, 0x12,
	0x77, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x46, 0x61, 0x6d, 0x69,
	0x6c, 0x79, 0x18, 0xc9, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x77, 0x65, 0x65, 0x6b, 0x64,
	0x61, 0x79, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x2b, 0x0a,
	0x10, 0x77, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x53, 0x69, 0x7a,
	0x65, 0x18, 0xca, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x77, 0x65, 0x65, 0x6b, 0x64, 0x61,
	0x79, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x25, 0x0a, 0x0d, 0x77, 0x65,
	0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0xcb, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x77, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x12, 0x1d, 0x0a, 0x09, 0x77, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x58, 0x18, 0xcc,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x77, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x58,
	0x12, 0x1d, 0x0a, 0x09, 0x77, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x59, 0x18, 0xcd, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x77, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x73, 0x59, 0x12,
	0x29, 0x0a, 0x0f, 0x73, 0x68, 0x6f, 0x77, 0x57, 0x65, 0x65, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x18, 0xac, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x73, 0x68, 0x6f, 0x77, 0x57,
	0x65, 0x65, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x35, 0x0a, 0x15, 0x77, 0x65,
	0x65, 0x6b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x46, 0x61, 0x6d,
	0x69, 0x6c, 0x79, 0x18, 0xad, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x77, 0x65, 0x65, 0x6b,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x46, 0x6f, 0x6e, 0x74, 0x46, 0x61, 0x6d, 0x69, 0x6c,
	0x79, 0x12, 0x31, 0x0a, 0x13, 0x77, 0x65, 0x65, 0x6b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x46, 0x6f, 0x6e, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x18, 0xae, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x13, 0x77, 0x65, 0x65, 0x6b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x46, 0x6f, 0x6e, 0x74,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x2b, 0x0a, 0x10, 0x77, 0x65, 0x65, 0x6b, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0xaf, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x10, 0x77, 0x65, 0x65, 0x6b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x12, 0x23, 0x0a, 0x0c, 0x77, 0x65, 0x65, 0x6b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x58, 0x18, 0xb0, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x77, 0x65, 0x65, 0x6b, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x73, 0x58, 0x12, 0x23, 0x0a, 0x0c, 0x77, 0x65, 0x65, 0x6b, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x73, 0x59, 0x18, 0xb1, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x77,
	0x65, 0x65, 0x6b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x59, 0x42, 0x24, 0x5a, 0x22, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x68, 0x75, 0x68, 0x6c, 0x6f,
	0x6d, 0x69, 0x6e, 0x2f, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x2f, 0x6d, 0x61, 0x69,
	0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_request_proto_rawDescOnce sync.Once
	file_request_proto_rawDescData = file_request_proto_rawDesc
)

func file_request_proto_rawDescGZIP() []byte {
	file_request_proto_rawDescOnce.Do(func() {
		file_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_request_proto_rawDescData)
	})
	return file_request_proto_rawDescData
}

var file_request_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_request_proto_goTypes = []interface{}{
	(*Request)(nil), // 0: main.Request
}
var file_request_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_request_proto_init() }
func file_request_proto_init() {
	if File_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_request_proto_goTypes,
		DependencyIndexes: file_request_proto_depIdxs,
		MessageInfos:      file_request_proto_msgTypes,
	}.Build()
	File_request_proto = out.File
	file_request_proto_rawDesc = nil
	file_request_proto_goTypes = nil
	file_request_proto_depIdxs = nil
}
