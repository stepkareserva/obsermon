// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=../../mocks/mock_service.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/stepkareserva/obsermon/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockGaugesService is a mock of GaugesService interface.
type MockGaugesService struct {
	ctrl     *gomock.Controller
	recorder *MockGaugesServiceMockRecorder
	isgomock struct{}
}

// MockGaugesServiceMockRecorder is the mock recorder for MockGaugesService.
type MockGaugesServiceMockRecorder struct {
	mock *MockGaugesService
}

// NewMockGaugesService creates a new mock instance.
func NewMockGaugesService(ctrl *gomock.Controller) *MockGaugesService {
	mock := &MockGaugesService{ctrl: ctrl}
	mock.recorder = &MockGaugesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGaugesService) EXPECT() *MockGaugesServiceMockRecorder {
	return m.recorder
}

// GetGauge mocks base method.
func (m *MockGaugesService) GetGauge(name string) (*models.Gauge, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", name)
	ret0, _ := ret[0].(*models.Gauge)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockGaugesServiceMockRecorder) GetGauge(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockGaugesService)(nil).GetGauge), name)
}

// ListGauges mocks base method.
func (m *MockGaugesService) ListGauges() ([]models.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGauges")
	ret0, _ := ret[0].([]models.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGauges indicates an expected call of ListGauges.
func (mr *MockGaugesServiceMockRecorder) ListGauges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGauges", reflect.TypeOf((*MockGaugesService)(nil).ListGauges))
}

// UpdateGauge mocks base method.
func (m *MockGaugesService) UpdateGauge(val models.Gauge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGauge", val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGauge indicates an expected call of UpdateGauge.
func (mr *MockGaugesServiceMockRecorder) UpdateGauge(val any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGauge", reflect.TypeOf((*MockGaugesService)(nil).UpdateGauge), val)
}

// MockCountersService is a mock of CountersService interface.
type MockCountersService struct {
	ctrl     *gomock.Controller
	recorder *MockCountersServiceMockRecorder
	isgomock struct{}
}

// MockCountersServiceMockRecorder is the mock recorder for MockCountersService.
type MockCountersServiceMockRecorder struct {
	mock *MockCountersService
}

// NewMockCountersService creates a new mock instance.
func NewMockCountersService(ctrl *gomock.Controller) *MockCountersService {
	mock := &MockCountersService{ctrl: ctrl}
	mock.recorder = &MockCountersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCountersService) EXPECT() *MockCountersServiceMockRecorder {
	return m.recorder
}

// GetCounter mocks base method.
func (m *MockCountersService) GetCounter(name string) (*models.Counter, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", name)
	ret0, _ := ret[0].(*models.Counter)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockCountersServiceMockRecorder) GetCounter(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockCountersService)(nil).GetCounter), name)
}

// ListCounters mocks base method.
func (m *MockCountersService) ListCounters() ([]models.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCounters")
	ret0, _ := ret[0].([]models.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCounters indicates an expected call of ListCounters.
func (mr *MockCountersServiceMockRecorder) ListCounters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCounters", reflect.TypeOf((*MockCountersService)(nil).ListCounters))
}

// UpdateCounter mocks base method.
func (m *MockCountersService) UpdateCounter(val models.Counter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounter", val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounter indicates an expected call of UpdateCounter.
func (mr *MockCountersServiceMockRecorder) UpdateCounter(val any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockCountersService)(nil).UpdateCounter), val)
}

// MockMetricsService is a mock of MetricsService interface.
type MockMetricsService struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsServiceMockRecorder
	isgomock struct{}
}

// MockMetricsServiceMockRecorder is the mock recorder for MockMetricsService.
type MockMetricsServiceMockRecorder struct {
	mock *MockMetricsService
}

// NewMockMetricsService creates a new mock instance.
func NewMockMetricsService(ctrl *gomock.Controller) *MockMetricsService {
	mock := &MockMetricsService{ctrl: ctrl}
	mock.recorder = &MockMetricsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsService) EXPECT() *MockMetricsServiceMockRecorder {
	return m.recorder
}

// GetMetric mocks base method.
func (m *MockMetricsService) GetMetric(t models.MetricType, name string) (*models.Metrics, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", t, name)
	ret0, _ := ret[0].(*models.Metrics)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockMetricsServiceMockRecorder) GetMetric(t, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockMetricsService)(nil).GetMetric), t, name)
}

// UpdateMetric mocks base method.
func (m *MockMetricsService) UpdateMetric(val models.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockMetricsServiceMockRecorder) UpdateMetric(val any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricsService)(nil).UpdateMetric), val)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetCounter mocks base method.
func (m *MockService) GetCounter(name string) (*models.Counter, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", name)
	ret0, _ := ret[0].(*models.Counter)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockServiceMockRecorder) GetCounter(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockService)(nil).GetCounter), name)
}

// GetGauge mocks base method.
func (m *MockService) GetGauge(name string) (*models.Gauge, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", name)
	ret0, _ := ret[0].(*models.Gauge)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockServiceMockRecorder) GetGauge(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockService)(nil).GetGauge), name)
}

// GetMetric mocks base method.
func (m *MockService) GetMetric(t models.MetricType, name string) (*models.Metrics, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", t, name)
	ret0, _ := ret[0].(*models.Metrics)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockServiceMockRecorder) GetMetric(t, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockService)(nil).GetMetric), t, name)
}

// ListCounters mocks base method.
func (m *MockService) ListCounters() ([]models.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCounters")
	ret0, _ := ret[0].([]models.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCounters indicates an expected call of ListCounters.
func (mr *MockServiceMockRecorder) ListCounters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCounters", reflect.TypeOf((*MockService)(nil).ListCounters))
}

// ListGauges mocks base method.
func (m *MockService) ListGauges() ([]models.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGauges")
	ret0, _ := ret[0].([]models.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGauges indicates an expected call of ListGauges.
func (mr *MockServiceMockRecorder) ListGauges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGauges", reflect.TypeOf((*MockService)(nil).ListGauges))
}

// UpdateCounter mocks base method.
func (m *MockService) UpdateCounter(val models.Counter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounter", val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounter indicates an expected call of UpdateCounter.
func (mr *MockServiceMockRecorder) UpdateCounter(val any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockService)(nil).UpdateCounter), val)
}

// UpdateGauge mocks base method.
func (m *MockService) UpdateGauge(val models.Gauge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGauge", val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGauge indicates an expected call of UpdateGauge.
func (mr *MockServiceMockRecorder) UpdateGauge(val any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGauge", reflect.TypeOf((*MockService)(nil).UpdateGauge), val)
}

// UpdateMetric mocks base method.
func (m *MockService) UpdateMetric(val models.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockServiceMockRecorder) UpdateMetric(val any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockService)(nil).UpdateMetric), val)
}
