package office

import (
	"context"
	"fmt"
	"go-clean/src/business/entity"
	"go-clean/src/lib/cloud_storage"

	"gorm.io/gorm"
)

type Interface interface {
	Create(office entity.Office) (entity.Office, error)
	GetList(param entity.OfficeParam) ([]entity.Office, error)
	GetPresignedURL(ctx context.Context, filename string) (string, error)
	GetListByLike(param entity.OfficeParam) ([]entity.Office, error)
	GetListByID(officeIDs []uint) ([]entity.Office, error)
	Get(param entity.OfficeParam) (entity.Office, error)
	Update(selectParam entity.OfficeParam, updateParam entity.UpdateOfficeParam) error
	Delete(param entity.OfficeParam) error
	GetCount(param entity.OfficeParam) (int64, error)
}

type office struct {
	db           *gorm.DB
	cloudStorage cloud_storage.Interface
}

func Init(db *gorm.DB, cs cloud_storage.Interface) Interface {
	o := &office{
		db:           db,
		cloudStorage: cs,
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

func (o *office) GetPresignedURL(ctx context.Context, filename string) (string, error) {
	url, err := o.cloudStorage.GetSignedURL(ctx, filename, "office-image/")
	if err != nil {
		return url, err
	}

	return url, nil
}

func (o *office) GetListByLike(param entity.OfficeParam) ([]entity.Office, error) {
	offices := []entity.Office{}

	if err := o.db.Where("name LIKE ? AND location LIKE ?", fmt.Sprintf("%%%s%%", param.Name), fmt.Sprintf("%%%s%%", param.Location)).Find(&offices).Error; err != nil {
		return offices, err
	}

	return offices, nil
}

func (o *office) GetListByID(officeIDs []uint) ([]entity.Office, error) {
	offices := []entity.Office{}
	if err := o.db.Find(&offices, officeIDs).Error; err != nil {
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

func (o *office) GetCount(param entity.OfficeParam) (int64, error) {
	var result int64
	if err := o.db.Model(entity.Office{}).Where(param).Count(&result).Error; err != nil {
		return result, err
	}

	return result, nil
}
