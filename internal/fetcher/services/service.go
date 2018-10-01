package services

type Service interface {
	FetchAndSave(done chan<- bool)
	GetStoreName() string
}
