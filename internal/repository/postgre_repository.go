package repository

func (db *Storage) Close() error {

	return db.postgre.Close()
}

func (db *Storage) CheckConnection() error {
	return db.postgre.CheckConnection()
}
