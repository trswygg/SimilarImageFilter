package Everything

import (
	"fmt"
	"runtime"
	"syscall"

	"github.com/juju/errors"
)

// Description about this package
var Description = `
    ______                      __  __    _
   / ____/   _____  _______  __/ /_/ /_  (_)___  ______
  / __/ | | / / _ \/ ___/ / / / __/ __ \/ / __ \/ __  /
 / /___ | |/ /  __/ /  / /_/ / /_/ / / / / / / / /_/ /
/_____/ |___/\___/_/   \__, /\__/_/ /_/_/_/ /_/\__, /
                     /____/                  /____/
					             powered by trswygg
Everything Lib for golang
See https://www.voidtools.com/zh-cn/support/everything/sdk/ for more details
Thanks https://github.com/YCF/Everything
require "Everything64.dll" or "Everything32.dll" (Depends on the number of bits of your operating system)
`
var _ = `
This file is based on "Everything.h" (Copyright (C) 2016 David Carpenter)
`

const API_VERSION = "0.1.2 alpha"

// error
const (
	EVERYTHING_OK                     = 0 // no error detected
	EVERYTHING_ERROR_MEMORY           = 1 // out of memory.
	EVERYTHING_ERROR_IPC              = 2 // Everything search client is not running
	EVERYTHING_ERROR_REGISTERCLASSEX  = 3 // unable to register window class.
	EVERYTHING_ERROR_CREATEWINDOW     = 4 // unable to create listening window
	EVERYTHING_ERROR_CREATETHREAD     = 5 // unable to create listening thread
	EVERYTHING_ERROR_INVALIDINDEX     = 6 // invalid index
	EVERYTHING_ERROR_INVALIDCALL      = 7 // invalid call
	EVERYTHING_ERROR_INVALIDREQUEST   = 8 // invalid request data, request data first.
	EVERYTHING_ERROR_INVALIDPARAMETER = 9 // bad parameter.
)

// sort
const (
	EVERYTHING_SORT_NAME_ASCENDING                   = 1  // sort name ascending
	EVERYTHING_SORT_NAME_DESCENDING                  = 2  // sort name descending
	EVERYTHING_SORT_PATH_ASCENDING                   = 3  //sort path ascending
	EVERYTHING_SORT_PATH_DESCENDING                  = 4  //sort path descending
	EVERYTHING_SORT_SIZE_ASCENDING                   = 5  //sort size ascending
	EVERYTHING_SORT_SIZE_DESCENDING                  = 6  //sort size descending
	EVERYTHING_SORT_EXTENSION_ASCENDING              = 7  //sort extension ascending
	EVERYTHING_SORT_EXTENSION_DESCENDING             = 8  //sort extension descending
	EVERYTHING_SORT_TYPE_NAME_ASCENDING              = 9  //sort type name ascending
	EVERYTHING_SORT_TYPE_NAME_DESCENDING             = 10 //sort type name descending
	EVERYTHING_SORT_DATE_CREATED_ASCENDING           = 11 //sort date created ascending
	EVERYTHING_SORT_DATE_CREATED_DESCENDING          = 12 //sort date created descending
	EVERYTHING_SORT_DATE_MODIFIED_ASCENDING          = 13 // sort date modified ascending
	EVERYTHING_SORT_DATE_MODIFIED_DESCENDING         = 14 // sort date modified descending
	EVERYTHING_SORT_ATTRIBUTES_ASCENDING             = 15 // sort attributes ascending
	EVERYTHING_SORT_ATTRIBUTES_DESCENDING            = 16 // sort attributes descending
	EVERYTHING_SORT_FILE_LIST_FILENAME_ASCENDING     = 17 //sort file list filename ascending
	EVERYTHING_SORT_FILE_LIST_FILENAME_DESCENDING    = 18 //sort  file list filename descending
	EVERYTHING_SORT_RUN_COUNT_ASCENDING              = 19 //sort run count ascending
	EVERYTHING_SORT_RUN_COUNT_DESCENDING             = 20 // sort run count descending
	EVERYTHING_SORT_DATE_RECENTLY_CHANGED_ASCENDING  = 21 //sort date recently changed ascending
	EVERYTHING_SORT_DATE_RECENTLY_CHANGED_DESCENDING = 22 // sort date recently changed descending
	EVERYTHING_SORT_DATE_ACCESSED_ASCENDING          = 23 // sort date ascending
	EVERYTHING_SORT_DATE_ACCESSED_DESCENDING         = 24 //sort  sort date descending
	EVERYTHING_SORT_DATE_RUN_ASCENDING               = 25 // sort date run ascending
	EVERYTHING_SORT_DATE_RUN_DESCENDING              = 26 //sort date run descending
)

// request
const (
	EVERYTHING_REQUEST_FILE_NAME                           = 0x00000001 //request file name
	EVERYTHING_REQUEST_PATH                                = 0x00000002 // request path
	EVERYTHING_REQUEST_FULL_PATH_AND_FILE_NAME             = 0x00000004 // request full path and file name
	EVERYTHING_REQUEST_EXTENSION                           = 0x00000008 // request extension
	EVERYTHING_REQUEST_SIZE                                = 0x00000010 // request size
	EVERYTHING_REQUEST_DATE_CREATED                        = 0x00000020 // request date created
	EVERYTHING_REQUEST_DATE_MODIFIED                       = 0x00000040 // request date modified
	EVERYTHING_REQUEST_DATE_ACCESSED                       = 0x00000080 // request date accessed
	EVERYTHING_REQUEST_ATTRIBUTES                          = 0x00000100 // request attributes
	EVERYTHING_REQUEST_FILE_LIST_FILE_NAME                 = 0x00000200 // request file list file name
	EVERYTHING_REQUEST_RUN_COUNT                           = 0x00000400 // request run count
	EVERYTHING_REQUEST_DATE_RUN                            = 0x00000800 // request date run
	EVERYTHING_REQUEST_DATE_RECENTLY_CHANGED               = 0x00001000 // request date recently changed
	EVERYTHING_REQUEST_HIGHLIGHTED_FILE_NAME               = 0x00002000 // request highlighted file name
	EVERYTHING_REQUEST_HIGHLIGHTED_PATH                    = 0x00004000 // request highlighted path
	EVERYTHING_REQUEST_HIGHLIGHTED_FULL_PATH_AND_FILE_NAME = 0x00008000 // request highlighted full path and file name
)

var arch = runtime.GOARCH
var dll *syscall.DLL
var version string
var debug = true

// init Load dll
func init() {
	if arch == "amd64" {
		dll = syscall.MustLoadDLL("Everything64.dll")
	} else if arch == "386" {
		dll = syscall.MustLoadDLL("Everything32.dll")
	}
	version = GetVersion()
	fmt.Println(dll.Name + " version: " + version)
	if debug {
		fmt.Println(GetVersion())
		fmt.Printf("API_VERSION %s", API_VERSION)
	}

}

func CallError(caller string) error {
	return errors.Errorf("Everything error: failed call %s", caller)
}
