package db

type DB struct {
	dbName string
	ownsCache bool
	ownsInfoLog bool
	logFileName uint64
}

func Write() Status {

}