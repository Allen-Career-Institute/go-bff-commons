// Code generated by MockGen. DO NOT EDIT.
// Source: v1/intrnl/clients/interfaces.go

// Package clients is a generated GoMock package.
package clients

import (
	reflect "reflect"

	request "github.com/Allen-Career-Institute/common-protos/cal/v1/request"
	response "github.com/Allen-Career-Institute/common-protos/cal/v1/response"
	request0 "github.com/Allen-Career-Institute/common-protos/resource/v1/request"
	response0 "github.com/Allen-Career-Institute/common-protos/resource/v1/response"
	response1 "github.com/Allen-Career-Institute/common-protos/user_management/v1/response"
	config "github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	grpc "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/grpc"
	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// GetAncestorsOfAFacility mocks base method.
func (m *MockManager) GetAncestorsOfAFacility(c echo.Context, request *request0.GetAncestorsOfAFacilityRequest) (*response0.GetAncestorsOfAFacilityResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAncestorsOfAFacility", c, request)
	ret0, _ := ret[0].(*response0.GetAncestorsOfAFacilityResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAncestorsOfAFacility indicates an expected call of GetAncestorsOfAFacility.
func (mr *MockManagerMockRecorder) GetAncestorsOfAFacility(c, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAncestorsOfAFacility", reflect.TypeOf((*MockManager)(nil).GetAncestorsOfAFacility), c, request)
}

// GetLearningMaterial mocks base method.
func (m *MockManager) GetLearningMaterial(c echo.Context, cnf *config.Config, req *request.GetLearningMaterialRequest) (*response.GetLearningMaterialResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLearningMaterial", c, cnf, req)
	ret0, _ := ret[0].(*response.GetLearningMaterialResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLearningMaterial indicates an expected call of GetLearningMaterial.
func (mr *MockManagerMockRecorder) GetLearningMaterial(c, cnf, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLearningMaterial", reflect.TypeOf((*MockManager)(nil).GetLearningMaterial), c, cnf, req)
}

// GetNAC mocks base method.
func (m *MockManager) GetNAC(c echo.Context, cnf *config.Config, req *request.GetNACRequest) (*response.GetNACResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNAC", c, cnf, req)
	ret0, _ := ret[0].(*response.GetNACResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNAC indicates an expected call of GetNAC.
func (mr *MockManagerMockRecorder) GetNAC(c, cnf, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNAC", reflect.TypeOf((*MockManager)(nil).GetNAC), c, cnf, req)
}

// GetStudentBatchDetails mocks base method.
func (m *MockManager) GetStudentBatchDetails(c echo.Context, arg1 *config.Config, request *request0.GetStudentBatchDetailsRequest) (*response0.GetStudentBatchDetailsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStudentBatchDetails", c, arg1, request)
	ret0, _ := ret[0].(*response0.GetStudentBatchDetailsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStudentBatchDetails indicates an expected call of GetStudentBatchDetails.
func (mr *MockManagerMockRecorder) GetStudentBatchDetails(c, arg1, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStudentBatchDetails", reflect.TypeOf((*MockManager)(nil).GetStudentBatchDetails), c, arg1, request)
}

// GetUser mocks base method.
func (m *MockManager) GetUser(c echo.Context, cnf *config.Config, grpcHandler grpc.Manager) (*response1.GetUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", c, cnf, grpcHandler)
	ret0, _ := ret[0].(*response1.GetUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockManagerMockRecorder) GetUser(c, cnf, grpcHandler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockManager)(nil).GetUser), c, cnf, grpcHandler)
}
