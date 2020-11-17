package database

type DbService interface {
	CreateDb(params ...string) ([]string, error)
	DeleteDb() error
}