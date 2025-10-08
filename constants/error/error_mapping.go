package error

func ErrMapping(err error) bool {
	allErrors := append([]error{}, GeneralError...)
	allErrors = append(allErrors, UserError...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}
	return false
}