package repository

type PostgreMethods interface {
	Close() error
	CheckConnection() error
}

type PostgreDB struct {
	Methods PostgreMethods
}

func NewPostgreDB(methods PostgreMethods) *PostgreDB {
	return &PostgreDB{Methods: methods}
}

func (db *PostgreDB) Close() error {

	return db.Methods.Close()
}

func (db *PostgreDB) CheckConnection() error {
	return db.Methods.CheckConnection()
}
