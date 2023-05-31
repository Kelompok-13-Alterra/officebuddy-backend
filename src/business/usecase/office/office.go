package office

import (
	officeDom "go-clean/src/business/domain/office"
	"go-clean/src/business/entity"
	"time"
)

type Interface interface {
	Create(param entity.CreateOfficeParam) (entity.Office, error)
	GetList(param entity.OfficeParam) ([]entity.Office, error)
	Get(param entity.OfficeParam) (entity.Office, error)
}

type office struct {
	office officeDom.Interface
}

func Init(od officeDom.Interface) Interface {
	o := &office{
		office: od,
	}

	return o
}

func (o *office) Create(param entity.CreateOfficeParam) (entity.Office, error) {
	office := entity.Office{}

	openTime, err := o.convertStringToOfficeHours(param.Open)
	if err != nil {
		return office, err
	}

	closeTime, err := o.convertStringToOfficeHours(param.Close)
	if err != nil {
		return office, err
	}

	office, err = o.office.Create(entity.Office{
		Name:        param.Name,
		Description: param.Description,
		Capacity:    param.Capacity,
		Type:        param.Type,
		Open:        openTime,
		Close:       closeTime,
		Location:    param.Location,
		Price:       param.Price,
		Facilities:  param.Facilities,
	})
	if err != nil {
		return office, err
	}

	return office, nil
}

func (o *office) convertStringToOfficeHours(s string) (entity.OfficeHours, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return entity.OfficeHours{}, err
	}
	return entity.OfficeHours{Time: t}, nil
}

func (o *office) GetList(param entity.OfficeParam) ([]entity.Office, error) {
	var (
		offices []entity.Office
		err     error
	)

	if param.Name != "" {
		offices, err = o.office.GetListByLikeName(param.Name)
	} else {
		offices, err = o.office.GetList(param)

	}
	if err != nil {
		return offices, err
	}

	return offices, nil
}

func (o *office) Get(param entity.OfficeParam) (entity.Office, error) {
	office, err := o.office.Get(param)
	if err != nil {
		return office, err
	}

	return office, nil
}
