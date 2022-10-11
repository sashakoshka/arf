package analyzer

func typeMismatchErrorMessage (source Type, destination Type) (message string) {
	message += source.Describe()
	message += " cannot be used as "
	message += destination.Describe()
	return
}
