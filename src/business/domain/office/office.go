package office

import (
	"fmt"
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(office entity.Office) (entity.Office, error)
	GetList(param entity.OfficeParam) ([]entity.Office, error)
	GetListByLike(param entity.OfficeParam) ([]entity.Office, error)
	Get(param entity.OfficeParam) (entity.Office, error)
	Update(selectParam entity.OfficeParam, updateParam entity.UpdateOfficeParam) error
	Delete(param entity.OfficeParam) error
}

type office struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	o := &office{
		db: db,
	}

	return o
}

func (o *office) Create(office entity.Office) (entity.Office, error) {
	if err := o.db.Create(&office).Error; err != nil {
		return office, err
	}

	return office, nil
}

func (o *office) GetList(param entity.OfficeParam) ([]entity.Office, error) {
	offices := []entity.Office{}

	if err := o.db.Where(param).Find(&offices).Error; err != nil {
		return offices, err
	}

	return offices, nil
}

func (o *office) GetListByLike(param entity.OfficeParam) ([]entity.Office, error) {
	offices := []entity.Office{}

	if err := o.db.Where("name LIKE ? AND location LIKE ?", fmt.Sprintf("%%%s%%", param.Name), fmt.Sprintf("%%%s%%", param.Location)).Find(&offices).Error; err != nil {
		return offices, err
	}

	return offices, nil
}

func (o *office) Get(param entity.OfficeParam) (entity.Office, error) {
	office := entity.Office{}

	if err := o.db.Where(param).First(&office).Error; err != nil {
		return office, err
	}

	return office, nil
}

func (o *office) Update(selectParam entity.OfficeParam, updateParam entity.UpdateOfficeParam) error {
	if err := o.db.Model(entity.Office{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}

func (o *office) Delete(param entity.OfficeParam) error {
	if err := o.db.Where(param).Delete(&entity.Office{}).Error; err != nil {
		return err
	}

	return nil
}
