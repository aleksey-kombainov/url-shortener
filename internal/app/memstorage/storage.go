package memstorage

import "errors"

type Storager interface {
	Put(key string, val string)
	GetValueByKey(key string) (string, error)
	GetKeyByValue(val string) (string, error)
}

type Storage struct {
	storage map[string]string
}

func (s *Storage) Put(key string, val string) {
	s.storage[key] = val
}

func (s Storage) GetValueByKey(key string) (string, error) {
	if val, ok := s.storage[key]; ok {
		return val, nil
	}
	return "", errors.New("key not found")
}

func (s Storage) GetKeyByValue(val string) (string, error) {
	for key, storageVal := range s.storage {
		if val == storageVal {
			return key, nil
		}
	}
	return "", errors.New("value not found")
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]string),
	}
}
