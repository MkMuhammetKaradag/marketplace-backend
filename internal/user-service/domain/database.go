package domain

type PostgresRepository interface {
	Close() error
}
