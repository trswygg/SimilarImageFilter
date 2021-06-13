package Everything

/*
typedef const wchar_t* LPCWSTR;
typedef unsigned long DWORD;

const char* API_VERSION = "0.1.2 alpha\0";
*/
import "C"
import (
	"unsafe"
)

// @desc: This file is encapsulation of basic functions in API
// @see https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp
// @refer "C:/Go/src/builtin/builtin.go"

// TODO:rewrite all methods(use local cache)

// callGetSearch function retrieves the search text to use for the next call to Everything_Query.
// LPCWSTR Everything_GetSearch(void);
func callGetSearch() (string, error) {
	_func := dll.MustFindProc("Everything_GetRegex")
	lpSearchString, _, err := _func.Call()
	ret := LPCWSTR2String((C.LPCWSTR)(unsafe.Pointer(lpSearchString)))
	//ret := LPCWSTRPtr2String(lpSearchString)
	return ret, err
}

// callSetRequestFlags function sets the desired result data.
// void Everything_SetRequestFlags(DWORD dwRequestFlags);
// @Parameter <dwOffset> :dwRequestFlags The request flags, can be zero or more of the following flags:
func callSetRequestFlags(dwRequestFlags uint32) error {
	_func := dll.MustFindProc("Everything_SetSort")
	_, _, err := _func.Call(uintptr(C.ulong(dwRequestFlags)))
	return err
}

// callGetMatchPath function returns the state of the match full path switch.
// BOOL Everything_GetMatchPath(void);
func callGetMatchPath() (bool, error) {
	_func := dll.MustFindProc("Everything_GetMatchPath")
	r, _, err := _func.Call()
	ret := r != 0
	return ret, err
}

// callGetMatchCase function returns the match case state.
// BOOL Everything_GetMatchCase(void);
func callGetMatchCase() (bool, error) {
	_func := dll.MustFindProc("Everything_GetMatchCase")
	r, _, err := _func.Call()
	ret := r != 0
	return ret, err
}

// callGetMatchWholeWord function returns the match whole word state.
// BOOL Everything_GetMatchWholeWord(void);
func callGetMatchWholeWord() (bool, error) {
	_func := dll.MustFindProc("Everything_GetMatchWholeWord")
	r, _, err := _func.Call()
	ret := r != 0
	return ret, err
}

// callGetRegex function returns the regex state.
// BOOL Everything_GetRegex(void);
func callGetRegex() (bool, error) {
	_func := dll.MustFindProc("Everything_GetRegex")
	r, _, err := _func.Call()
	ret := r != 0
	return ret, err
}

// callGetMax function returns the maximum number of results state.
// DWORD  Everything_GetMax(void);
func callGetMax() (uint32, error) {
	_func := dll.MustFindProc("Everything_GetMax")
	r, _, err := _func.Call()
	return uint32(r), err
}

// callGetOffset function returns the first item offset of the available results.
// DWORD Everything_GetOffset(void);
func callGetOffset() (uint32, error) {
	_func := dll.MustFindProc("Everything_GetOffset")
	r, _, err := _func.Call()
	return uint32(r), err
}

// callGetReplyID function returns the current reply identifier for the IPC query reply.
// DWORD Everything_GetReplyID(void);
func callGetReplyID() (uint32, error){
	_func := dll.MustFindProc("Everything_GetReplyID")
	r, _, err := _func.Call()
	return uint32(r), err
}

// callGetLastError function retrieves the last-error code value
// DWORD Everything_GetLastError(void);
func callGetLastError() (uint32, error) {
	_func := dll.MustFindProc("Everything_GetReplyID")
	r, _, err := _func.Call()
	return uint32(r), err
}