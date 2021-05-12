package main
//go install github.com/golang/mock/mockgen@v1.5.0

//go:generate mockgen -destination mock_db/mock_db.go -source=db.go -package=mock_db

type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}
	return -1
}
