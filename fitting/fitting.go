package fitting
/*
#cgo LDFLAGS: -L. -ldl -lopencv_core
#include "cfitting.h"
*/
import "C"
import "unsafe"

func CurveFittingArr(X []float64,Y []float64,MaxLen int) (uint64) {
	_x := (*C.double)(unsafe.Pointer(&X[0]))
	_y := (*C.double)(unsafe.Pointer(&Y[0]))
	Len := C.int(len(X))
	Max:=C.int(MaxLen)
	key := C.GetCurveArr(_x,_y,Len,Max)
	return uint64(key)
}
func CurveFitting(data []byte,MaxLen int) (uint64,error) {
	db :=(*C.char)(unsafe.Pointer(&data[0]))
	Len:=C.int(len(data))
	Max:=C.int(MaxLen)
	key :=C.GetCurveData(db,Len,Max)
	return uint64(key),nil
}
