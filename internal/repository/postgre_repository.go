package repository

type Storage struct {
	postgre Postgre
}

func NewStorage(postgre Postgre) *Storage {
	return &Storage{postgre: postgre}
}
func (s *Storage) Close() error {

	return s.postgre.Close()
}

func (s *Storage) CheckConnection() error {

	return s.postgre.CheckConnection()
}

func (s *Storage) SetGauge(key string, value float64) {
	s.postgre.SetGauge(key, value)
}
func (s *Storage) SetCounter(key string, value int64) {
	s.postgre.SetCounter(key, value)
}

func (s *Storage) GetAllGauges() (map[string]float64, error) {
	return s.postgre.GetAllGauges()
}

func (s *Storage) GetAllCounters() (map[string]int64, error) {
	return s.postgre.GetAllCounters()
}
