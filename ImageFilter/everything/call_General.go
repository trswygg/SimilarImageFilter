package Everything

// CleanUp function resets the result list and search state, freeing any allocated memory by the library.
// void Everything_CleanUp( void);
func CleanUp() error {
	_func := dll.MustFindProc("Everything_CleanUp")
	_, _, err := _func.Call()
	return err
}

// DeleteRunHistory function deletes all run history.
// BOOL Everything_DeleteRunHistory( void);
func DeleteRunHistory() (error, bool) {
	_func := dll.MustFindProc("Everything_DeleteRunHistory")
	r, _, err := _func.Call()
	res := r != 0
	return err, res
}

