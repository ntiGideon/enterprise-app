package service

//
//import (
//	"Enterprise/prisma/db"
//	"context"
//	"github.com/stretchr/testify/mock"
//	"testing"
//)
//
//type MockUserDB struct {
//	mock.Mock
//}
//
//func (m *MockUserDB) FindMany() *db.UserQuery {
//	args := m.Called()
//	return args.Get(0).(*db.UserManyTxResult)
//}
//
//type MockUserQuery struct {
//	mock.Mock
//}
//
//func (q *MockUserQuery) Select(field ...interface{}) *MockUserQuery {
//	q.Called(field)
//	return q
//}
//
//func (q *MockUserQuery) With(roleFetch interface{}) *MockUserQuery {
//	q.Called(roleFetch)
//	return q
//}
//
//func (q *MockUserQuery) Exec(ctx context.Context) ([]db.UserManyTxResult, error) {
//	args := q.Called(ctx)
//	return args.Get(0).([]db.UserManyTxResult), args.Error(1)
//}
//
//
//func TestUserService_GetAllUsers(t *testing.T) {
//	mockDB := new(MockUserDB)
//	mockQuery := new(MockUserQuery)
//
//	mockUsers := []db.UserQuery
//}
