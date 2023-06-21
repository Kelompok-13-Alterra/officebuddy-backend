package office

import (
	"context"
	"fmt"
	officeDom "go-clean/src/business/domain/office"
	"go-clean/src/business/entity"
	"go-clean/src/lib/cloud_storage"
	"log"
	"mime/multipart"
	"time"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateOfficeParam) (entity.Office, error)
	UploadImage(ctx context.Context, blobFile multipart.File, param entity.OfficeParam) error
	GetList(ctx context.Context, param entity.OfficeParam) ([]entity.Office, error)
	Get(ctx context.Context, param entity.OfficeParam) (entity.Office, error)
	Update(param entity.OfficeParam, inputParam entity.UpdateOfficeParam) error
	Delete(param entity.OfficeParam) error
}

type office struct {
	office       officeDom.Interface
	cloudStorage cloud_storage.Interface
}

func Init(od officeDom.Interface, cs cloud_storage.Interface) Interface {
	o := &office{
		office:       od,
		cloudStorage: cs,
	}

	return o
}

func (o *office) Create(ctx context.Context, param entity.CreateOfficeParam) (entity.Office, error) {
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

func (o *office) UploadImage(ctx context.Context, blobFile multipart.File, param entity.OfficeParam) error {
	fileName := fmt.Sprintf("%d-image-%d", param.ID, time.Now().Unix())
	if err := o.cloudStorage.UploadFile(ctx, blobFile, fileName, "office-image/"); err != nil {
		return err
	}

	if err := o.office.Update(entity.OfficeParam{
		ID: param.ID,
	}, entity.UpdateOfficeParam{
		ImageUrl: fileName,
	}); err != nil {
		return err
	}

	return nil
}

func (o *office) convertStringToOfficeHours(s string) (entity.OfficeHours, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return entity.OfficeHours{}, err
	}
	return entity.OfficeHours{Time: t}, nil
}

func (o *office) GetList(ctx context.Context, param entity.OfficeParam) ([]entity.Office, error) {
	var (
		offices []entity.Office
		err     error
	)

	if param.Name != "" || param.Location != "" {
		offices, err = o.office.GetListByLike(param)
	} else {
		offices, err = o.office.GetList(param)

	}
	if err != nil {
		return offices, err
	}

	for i, of := range offices {
		if of.ImageUrl == "" {
			continue
		}
		url, err := o.office.GetPresignedURL(ctx, of.ImageUrl)
		if err != nil {
			log.Println(err)
		}
		offices[i].ImageUrl = url
	}

	return offices, nil
}

func (o *office) Get(ctx context.Context, param entity.OfficeParam) (entity.Office, error) {
	office, err := o.office.Get(param)
	if err != nil {
		return office, err
	}

	if office.ImageUrl != "" {
		url, err := o.office.GetPresignedURL(ctx, office.ImageUrl)
		if err != nil {
			log.Println(err)
		}
		office.ImageUrl = url
	}

	return office, nil
}

func (o *office) Update(param entity.OfficeParam, inputParam entity.UpdateOfficeParam) error {
	if inputParam.OpenHours != "" {
		openTime, err := o.convertStringToOfficeHours(inputParam.OpenHours)
		if err != nil {
			return err
		}
		inputParam.Open = openTime
	}

	if inputParam.CloseHours != "" {
		closeTime, err := o.convertStringToOfficeHours(inputParam.CloseHours)
		if err != nil {
			return err
		}
		inputParam.Close = closeTime
	}

	if err := o.office.Update(param, inputParam); err != nil {
		return err
	}

	return nil
}

func (o *office) Delete(param entity.OfficeParam) error {
	if err := o.office.Delete(entity.OfficeParam{
		ID: param.ID,
	}); err != nil {
		return err
	}

	return nil
}
