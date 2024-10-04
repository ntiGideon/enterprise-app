package helpers

func PanicAllErrors(err error) {
	if err != nil {
		panic(err)
	}
}
