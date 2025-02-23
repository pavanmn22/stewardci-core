// /*
// #########################
// #  SAP Steward-CI       #
// #########################
//
// THIS CODE IS GENERATED! DO NOT TOUCH!
//
// Copyright SAP SE.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */
//

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SAP/stewardci-core/pkg/runctl/metrics (interfaces: CounterMetric,PipelineRunsMetric,StateItemsMetric,ResultsMetric)

// Package testing is a generated GoMock package.
package testing

import (
	v1alpha1 "github.com/SAP/stewardci-core/pkg/apis/steward/v1alpha1"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCounterMetric is a mock of CounterMetric interface
type MockCounterMetric struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMetricMockRecorder
}

// MockCounterMetricMockRecorder is the mock recorder for MockCounterMetric
type MockCounterMetricMockRecorder struct {
	mock *MockCounterMetric
}

// NewMockCounterMetric creates a new mock instance
func NewMockCounterMetric(ctrl *gomock.Controller) *MockCounterMetric {
	mock := &MockCounterMetric{ctrl: ctrl}
	mock.recorder = &MockCounterMetricMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCounterMetric) EXPECT() *MockCounterMetricMockRecorder {
	return m.recorder
}

// Inc mocks base method
func (m *MockCounterMetric) Inc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc
func (mr *MockCounterMetricMockRecorder) Inc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounterMetric)(nil).Inc))
}

// MockPipelineRunsMetric is a mock of PipelineRunsMetric interface
type MockPipelineRunsMetric struct {
	ctrl     *gomock.Controller
	recorder *MockPipelineRunsMetricMockRecorder
}

// MockPipelineRunsMetricMockRecorder is the mock recorder for MockPipelineRunsMetric
type MockPipelineRunsMetricMockRecorder struct {
	mock *MockPipelineRunsMetric
}

// NewMockPipelineRunsMetric creates a new mock instance
func NewMockPipelineRunsMetric(ctrl *gomock.Controller) *MockPipelineRunsMetric {
	mock := &MockPipelineRunsMetric{ctrl: ctrl}
	mock.recorder = &MockPipelineRunsMetricMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPipelineRunsMetric) EXPECT() *MockPipelineRunsMetricMockRecorder {
	return m.recorder
}

// Observe mocks base method
func (m *MockPipelineRunsMetric) Observe(arg0 *v1alpha1.PipelineRun) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Observe", arg0)
}

// Observe indicates an expected call of Observe
func (mr *MockPipelineRunsMetricMockRecorder) Observe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Observe", reflect.TypeOf((*MockPipelineRunsMetric)(nil).Observe), arg0)
}

// MockStateItemsMetric is a mock of StateItemsMetric interface
type MockStateItemsMetric struct {
	ctrl     *gomock.Controller
	recorder *MockStateItemsMetricMockRecorder
}

// MockStateItemsMetricMockRecorder is the mock recorder for MockStateItemsMetric
type MockStateItemsMetricMockRecorder struct {
	mock *MockStateItemsMetric
}

// NewMockStateItemsMetric creates a new mock instance
func NewMockStateItemsMetric(ctrl *gomock.Controller) *MockStateItemsMetric {
	mock := &MockStateItemsMetric{ctrl: ctrl}
	mock.recorder = &MockStateItemsMetricMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStateItemsMetric) EXPECT() *MockStateItemsMetricMockRecorder {
	return m.recorder
}

// Observe mocks base method
func (m *MockStateItemsMetric) Observe(arg0 *v1alpha1.StateItem) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Observe", arg0)
}

// Observe indicates an expected call of Observe
func (mr *MockStateItemsMetricMockRecorder) Observe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Observe", reflect.TypeOf((*MockStateItemsMetric)(nil).Observe), arg0)
}

// MockResultsMetric is a mock of ResultsMetric interface
type MockResultsMetric struct {
	ctrl     *gomock.Controller
	recorder *MockResultsMetricMockRecorder
}

// MockResultsMetricMockRecorder is the mock recorder for MockResultsMetric
type MockResultsMetricMockRecorder struct {
	mock *MockResultsMetric
}

// NewMockResultsMetric creates a new mock instance
func NewMockResultsMetric(ctrl *gomock.Controller) *MockResultsMetric {
	mock := &MockResultsMetric{ctrl: ctrl}
	mock.recorder = &MockResultsMetricMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResultsMetric) EXPECT() *MockResultsMetricMockRecorder {
	return m.recorder
}

// Observe mocks base method
func (m *MockResultsMetric) Observe(arg0 v1alpha1.Result) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Observe", arg0)
}

// Observe indicates an expected call of Observe
func (mr *MockResultsMetricMockRecorder) Observe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Observe", reflect.TypeOf((*MockResultsMetric)(nil).Observe), arg0)
}
