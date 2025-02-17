// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	types "github.com/smartcontractkit/ocr2keepers/pkg/types"
	mock "github.com/stretchr/testify/mock"
)

// ReportEncoder is an autogenerated mock type for the ReportEncoder type
type ReportEncoder struct {
	mock.Mock
}

// DecodeReport provides a mock function with given fields: _a0
func (_m *ReportEncoder) DecodeReport(_a0 []byte) ([]types.UpkeepResult, error) {
	ret := _m.Called(_a0)

	var r0 []types.UpkeepResult
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) ([]types.UpkeepResult, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]byte) []types.UpkeepResult); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.UpkeepResult)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncodeReport provides a mock function with given fields: _a0
func (_m *ReportEncoder) EncodeReport(_a0 []types.UpkeepResult) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func([]types.UpkeepResult) ([]byte, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func([]types.UpkeepResult) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func([]types.UpkeepResult) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewReportEncoder interface {
	mock.TestingT
	Cleanup(func())
}

// NewReportEncoder creates a new instance of ReportEncoder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReportEncoder(t mockConstructorTestingTNewReportEncoder) *ReportEncoder {
	mock := &ReportEncoder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
