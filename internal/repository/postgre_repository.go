package repository

func (s *Storage) Close() error {

	return s.postgre.Close()
}

func (s *Storage) CheckConnection() error {
	return s.postgre.CheckConnection()
}

func (s *Storage) SetGaugeStrorage(key string, value float64) {
	s.postgre.SetGauge(key, value)
}
func (s *Storage) SetCounterStorage(key string, value int64) {
	s.postgre.SetCounter(key, value)
}
func (s *Storage) AddCounterStorage(key string, value int64) bool {
	return s.postgre.AddCounter(key, value)
}
