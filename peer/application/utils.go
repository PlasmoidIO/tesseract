package application

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
