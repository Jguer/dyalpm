package lib

// FunctionCaller is implemented by MockFunctionRegistry for tests.
type FunctionCaller interface {
	GetFunc(name string) (uintptr, error)
	CallFunc(name string, args ...uintptr) (uintptr, error)
}
