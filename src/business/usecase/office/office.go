package office

import (
	officeDom "go-clean/src/business/domain/office"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetList(param entity.OfficeParam) ([]entity.Office, error)
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