package Everything

/*
typedef const wchar_t* LPCWSTR;
typedef unsigned long DWORD;

const char* API_VERSION = "0.1.2 alpha\0";
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// SetSearch function sets the search string for the IPC Query.
// void Everything_SetSearch( LPCWSTR lpString);
// @Parameter <lpString> : Pointer to a null-terminated string to be used as the new search text.
func SetSearch(lpString string) error {
	return callSetSearchW(lpString)
}

// callSetSearchW implementation Everything_SetSearch
func callSetSearchW(lpString string) error {
	_func := dll.MustFindProc("Everything_SetSearchW")
	_lpString := String2LPCWSTR(lpString)
	_, _, err := _func.Call(uintptr(unsafe.Pointer(_lpString)))
	return err
}

// callSetMatchPath function enables or disables full path matching for the next call to Everything_Query.
// void Everything_SetMatchPath( BOOL bEnable);
// @Parameter <bEnable> : Specifies whether to enable or disable full path matching.
func callSetMatchPath(bEnable bool) error {
	_func := dll.MustFindProc("Everything_SetMatchPath")
	_, _, err := _func.Call(uintptr(Bool2BOOL(bEnable)))
	return err
}

// callSetMatchCase function enables or disables full path matching for the next call to Everything_Query.
// void Everything_SetMatchCase( BOOL bEnable);
// @Parameter <bEnable> :Specifies whether the search is case sensitive or insensitive.
func callSetMatchCase(bEnable bool) error {
	_func := dll.MustFindProc("Everything_SetMatchCase")
	_, _, err := _func.Call(uintptr(Bool2BOOL(bEnable)))
	return err
}

// callSetMatchWholeWord function enables or disables matching whole words for the next call to Everything_Query.
// void Everything_SetMatchWholeWord( BOOL bEnable);
// @Parameter <bEnable> :Specifies whether the search matches whole words, or can match anywhere.
func callSetMatchWholeWord(bEnable bool) error {
	_func := dll.MustFindProc("Everything_SetMatchWholeWord")
	_, _, err := _func.Call(uintptr(Bool2BOOL(bEnable)))
	return err
}

// callSetRegex function enables or disables Regular Expression searching.
// void Everything_SetRegex( BOOL bEnabled);
// @Parameter <bEnable> :Set to non-zero to enable regex, set to zero to disable regex.
func callSetRegex(bEnable bool) error {
	_func := dll.MustFindProc("Everything_SetRegex")
	_, _, err := _func.Call(uintptr(Bool2BOOL(bEnable)))
	return err
}

// callSetMax function set the maximum number of results to return from Everything_Query.
// void Everything_SetMax( DWORD dwMaxResults);
// @Parameter <bEnable> :Specifies the maximum number of results to return.
func callSetMax(dwMax uint32) error {
	_func := dll.MustFindProc("Everything_SetMax")
	_, _, err := _func.Call(uintptr(C.ulong(dwMax)))
	return err
}

// callSetOffset function set the first result offset to return from a call to Everything_Query.
// void Everything_SetOffset( DWORD dwOffset);
// @Parameter <dwOffset> :Specifies the first result from the available results.
func callSetOffset(dwOffset string) error {
	_func := dll.MustFindProc("Everything_SetOffset")
	_dwOffset := String2LPCWSTR(dwOffset)
	_, _, err := _func.Call(uintptr(unsafe.Pointer(_dwOffset)))
	return err
}

// callSetReplyWindow function sets the window that will receive the the IPC Query results.
//void Everything_SetReplyWindow( HWND hWnd);
func callSetReplyWindow(hWnd uintptr) error {
	fmt.Print(fmt.Errorf("[Everything] error:func SetReplyWindow() unrealized"))
	_func := dll.MustFindProc("Everything_SetReplyWindow")
	_, _, err := _func.Call(uintptr(hWnd))
	return err
}

// callSetReplyID function sets the unique number to identify the next query.
// void Everything_SetReplyID(DWORD dwId);
func callSetReplyID(nId uint32) error {
	_func := dll.MustFindProc("Everything_SetReplyID")
	_, _, err := _func.Call(uintptr(C.ulong(nId)))
	return err
}

// callSetSort function sets how the results should be ordered.
// void Everything_SetSort(DWORD dwSort);
// @Parameter <dwSortType> :dwSortType The sort type, can be one of the following values:
func callSetSort(dwSortType uint32) error {
	_func := dll.MustFindProc("Everything_SetSort")
	_, _, err := _func.Call(uintptr(C.ulong(dwSortType)))
	return err
}