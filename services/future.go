package services

func LoopingExchangeInfo(path string) {
	store, err := OpenDataSQLite3Store(path)
	if err != nil {
		panic(err)
	}
}
