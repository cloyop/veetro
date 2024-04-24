package storage

import (
	"fmt"
	"time"
)

type Application struct {
	Id               string            `json:"id"`
	OfferId          string            `json:"offerId" bson:"offer_id"`
	EmployeeLocation string            `json:"employeeLocation" bson:"employee_location"`
	EmployeeId       string            `json:"employeeId" bson:"employee_id"`
	AppliedTime      int64             `json:"appliedTime" bson:"applied_time"`
	Resume           string            `json:"resume"`
	Anwsers          map[string]string `json:"answers"`
	CoverLetter      string            `json:"coverLetter" bson:"cover_letter"`
}

func NewApplication(offerId, employeeId, resume, coverLetter, employeeLocation string, answers *map[string]string) *Application {
	now := time.Now().Unix()
	return &Application{
		Id:               fmt.Sprintf("%s%d%s", employeeId[:8], now, employeeId[len(employeeId)-4:]),
		OfferId:          offerId,
		EmployeeId:       employeeId,
		Resume:           resume,
		CoverLetter:      coverLetter,
		AppliedTime:      now,
		EmployeeLocation: employeeLocation,
		Anwsers:          *answers,
	}
}
func (a *Application) ParseUpdateFields() (*map[string]string, *[]string) {
	updts := map[string]string{}
	errs := []string{}
	if a.CoverLetter != "" {
		updts["cover_letter"] = a.CoverLetter
	}
	if a.Resume != "" {
		if a.ValidResume() {
			updts["resume"] = a.Resume
		} else {
			errs = append(errs, "Invalid Resume")
		}
	}
	return &updts, &errs
}
func (a *Application) ValidCoverLetter() bool {
	return len(a.CoverLetter) > 250
}
func (a *Application) ValidResume() bool {
	return len(a.Resume) > 300
}
