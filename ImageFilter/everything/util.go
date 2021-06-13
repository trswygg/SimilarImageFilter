package Everything

// #include <windows.h>
// #include <wchar.h>
// #include <WinNT.h>
import "C"

import (
	"unicode/utf16"
	"unsafe"
)

const _ = `
This file 
THANKS https://github.com/icholy
`

// maxRunes = 1073741823
const maxRunes = 1<<30 - 1

// C.LPCWSTR
func LPCWSTR2String(lpwstr C.LPCWSTR) string {
	ptr := unsafe.Pointer(&lpwstr)
	sz := C.wcslen((*C.wchar_t)(ptr))
	wstr := (*[maxRunes]uint16)(ptr)[:sz:sz]
	return string(utf16.Decode(wstr))
}
func LPCWSTRPtr2String(ptr uintptr) string {
	p := unsafe.Pointer(ptr)
	sz := C.wcslen((*C.wchar_t)(p))
	wstr := (*[maxRunes]uint16)(p)[:sz:sz]
	return string(utf16.Decode(wstr))
}

func String2LPCWSTR(s string) C.LPCWSTR {
	awstr := utf16.Encode([]rune(s))
	awstr = append(awstr, 0x00)
	return (C.LPCWSTR)(unsafe.Pointer(&awstr[0]))
}

func Bool2BOOL(bEnable bool) int {
	var pEnable = 0
	if bEnable {
		pEnable = 1
	}
	return pEnable
}