package storage

type StorageService interface {
	Init()
	Close()
	userStorageService
	offerStorageService
	appliationStorageService
}
type userStorageService interface {
	GetUserByEmail(string) (*User, bool, error)
	CreateUser(*User) error
	UserExist(string) (bool, error)
	UpdateUser(userId string, updts *map[string]interface{}) (bool, error)
	DeleteCustomer(userId string) (bool, error)
	DeleteEmployee(userId string) (bool, error)
}
type offerStorageService interface {
	GetOfferWithApplications(offerId, userId string) (*Offer, bool, error)
	GetAllOffersWithApplications(userId string) (*[]Offer, error)
	// With Filters
	GetAllOffers(keyword, role, location string, page int) (*[]Offer, int64, error)
	GetOffer(offerId string) (*Offer, bool, error)
	//
	OfferExist(string) (bool, error)
	UpdateOffer(offerId, userId string, updts *map[string]interface{}) (bool, bool, error)
	CreateOffer(*Offer) error
	DeleteOffer(string, string) (bool, error)
	DeleteOffersByIds(*[]string, string) (int64, error)
}
type appliationStorageService interface {
	GetAllApplications(offerId, employeeId string) (*[]Application, error)
	GetApplication(id, offerId, employeeId string) (*Application, bool, error)
	ApplicationExist(id, offerId, employeeId string) (bool, error)
	UpdateApplication(id, employeeId string, updts *map[string]string) (bool, bool, error)
	CreateApplication(*Application) error
	DeleteApplication(id, offerId, employeeId string) (bool, error)
	DeleteApplicationsByIds(*[]string, string) (int64, error)
}
