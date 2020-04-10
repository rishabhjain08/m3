// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/cmd/services/m3coordinator/downsample (interfaces: Downsampler,MetricsAppender,SamplesAppender)

// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package downsample is a generated GoMock package.
package downsample

import (
	"reflect"
	"time"

	"github.com/golang/mock/gomock"
)

// MockDownsampler is a mock of Downsampler interface
type MockDownsampler struct {
	ctrl     *gomock.Controller
	recorder *MockDownsamplerMockRecorder
}

// MockDownsamplerMockRecorder is the mock recorder for MockDownsampler
type MockDownsamplerMockRecorder struct {
	mock *MockDownsampler
}

// NewMockDownsampler creates a new mock instance
func NewMockDownsampler(ctrl *gomock.Controller) *MockDownsampler {
	mock := &MockDownsampler{ctrl: ctrl}
	mock.recorder = &MockDownsamplerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDownsampler) EXPECT() *MockDownsamplerMockRecorder {
	return m.recorder
}

// NewMetricsAppender mocks base method
func (m *MockDownsampler) NewMetricsAppender() (MetricsAppender, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMetricsAppender")
	ret0, _ := ret[0].(MetricsAppender)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMetricsAppender indicates an expected call of NewMetricsAppender
func (mr *MockDownsamplerMockRecorder) NewMetricsAppender() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMetricsAppender", reflect.TypeOf((*MockDownsampler)(nil).NewMetricsAppender))
}

// MockMetricsAppender is a mock of MetricsAppender interface
type MockMetricsAppender struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsAppenderMockRecorder
}

// MockMetricsAppenderMockRecorder is the mock recorder for MockMetricsAppender
type MockMetricsAppenderMockRecorder struct {
	mock *MockMetricsAppender
}

// NewMockMetricsAppender creates a new mock instance
func NewMockMetricsAppender(ctrl *gomock.Controller) *MockMetricsAppender {
	mock := &MockMetricsAppender{ctrl: ctrl}
	mock.recorder = &MockMetricsAppenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetricsAppender) EXPECT() *MockMetricsAppenderMockRecorder {
	return m.recorder
}

// AddTag mocks base method
func (m *MockMetricsAppender) AddTag(arg0, arg1 []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddTag", arg0, arg1)
}

// AddTag indicates an expected call of AddTag
func (mr *MockMetricsAppenderMockRecorder) AddTag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTag", reflect.TypeOf((*MockMetricsAppender)(nil).AddTag), arg0, arg1)
}

// Finalize mocks base method
func (m *MockMetricsAppender) Finalize() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Finalize")
}

// Finalize indicates an expected call of Finalize
func (mr *MockMetricsAppenderMockRecorder) Finalize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finalize", reflect.TypeOf((*MockMetricsAppender)(nil).Finalize))
}

// IsDropPolicyApplied mocks base method
func (m *MockMetricsAppender) IsDropPolicyApplied() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDropPolicyApplied")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsDropPolicyApplied indicates an expected call of IsDropPolicyApplied
func (mr *MockMetricsAppenderMockRecorder) IsDropPolicyApplied() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDropPolicyApplied", reflect.TypeOf((*MockMetricsAppender)(nil).IsDropPolicyApplied))
}

// Reset mocks base method
func (m *MockMetricsAppender) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset
func (mr *MockMetricsAppenderMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockMetricsAppender)(nil).Reset))
}

// SamplesAppender mocks base method
func (m *MockMetricsAppender) SamplesAppender(arg0 SampleAppenderOptions) (SamplesAppender, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SamplesAppender", arg0)
	ret0, _ := ret[0].(SamplesAppender)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SamplesAppender indicates an expected call of SamplesAppender
func (mr *MockMetricsAppenderMockRecorder) SamplesAppender(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SamplesAppender", reflect.TypeOf((*MockMetricsAppender)(nil).SamplesAppender), arg0)
}

// MockSamplesAppender is a mock of SamplesAppender interface
type MockSamplesAppender struct {
	ctrl     *gomock.Controller
	recorder *MockSamplesAppenderMockRecorder
}

// MockSamplesAppenderMockRecorder is the mock recorder for MockSamplesAppender
type MockSamplesAppenderMockRecorder struct {
	mock *MockSamplesAppender
}

// NewMockSamplesAppender creates a new mock instance
func NewMockSamplesAppender(ctrl *gomock.Controller) *MockSamplesAppender {
	mock := &MockSamplesAppender{ctrl: ctrl}
	mock.recorder = &MockSamplesAppenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSamplesAppender) EXPECT() *MockSamplesAppenderMockRecorder {
	return m.recorder
}

// AppendCounterSample mocks base method
func (m *MockSamplesAppender) AppendCounterSample(arg0 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendCounterSample", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendCounterSample indicates an expected call of AppendCounterSample
func (mr *MockSamplesAppenderMockRecorder) AppendCounterSample(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendCounterSample", reflect.TypeOf((*MockSamplesAppender)(nil).AppendCounterSample), arg0)
}

// AppendCounterTimedSample mocks base method
func (m *MockSamplesAppender) AppendCounterTimedSample(arg0 time.Time, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendCounterTimedSample", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendCounterTimedSample indicates an expected call of AppendCounterTimedSample
func (mr *MockSamplesAppenderMockRecorder) AppendCounterTimedSample(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendCounterTimedSample", reflect.TypeOf((*MockSamplesAppender)(nil).AppendCounterTimedSample), arg0, arg1)
}

// AppendGaugeSample mocks base method
func (m *MockSamplesAppender) AppendGaugeSample(arg0 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendGaugeSample", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendGaugeSample indicates an expected call of AppendGaugeSample
func (mr *MockSamplesAppenderMockRecorder) AppendGaugeSample(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendGaugeSample", reflect.TypeOf((*MockSamplesAppender)(nil).AppendGaugeSample), arg0)
}

// AppendGaugeTimedSample mocks base method
func (m *MockSamplesAppender) AppendGaugeTimedSample(arg0 time.Time, arg1 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendGaugeTimedSample", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendGaugeTimedSample indicates an expected call of AppendGaugeTimedSample
func (mr *MockSamplesAppenderMockRecorder) AppendGaugeTimedSample(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendGaugeTimedSample", reflect.TypeOf((*MockSamplesAppender)(nil).AppendGaugeTimedSample), arg0, arg1)
}
