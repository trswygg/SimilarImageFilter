package Everything

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetVersion returns DLL version
// return MajorVersion + '+ MinorVersion + '.' + Revision
// https://githuom/YCF/Everything/blob/79286e2a608cecc879d808093a714cfc3fd207e9/everything_windows_amd64.go#L164
func GetVersion() string {
	if dll != nil {
		fmajor := dll.MustFindProc("Everything_GetMajorVersion")
		fminor := dll.MustFindProc("Everything_GetMinorVersion")
		frevision := dll.MustFindProc("Everything_GetRevision")
		p1, _, _ := fmajor.Call()
		p2, _, _ := fminor.Call()
		p3, _, _ := frevision.Call()
		return fmt.Sprintf("%d.%d.%d", int(p1), int(p2), int(p3))
	}
	return ""
}

// GetCurrentDirectory returns the absolute path to which the program runs
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}