package repository

import (
	"context"
	"next-terminal/server/model"
)

var LicenseRepository = new(licenseRepository)

type licenseRepository struct {
	baseRepository
}

func (r licenseRepository) FindLicense(c context.Context) (model.License, error) {
	var license model.License
	err := r.GetDB(c).First(&license).Error
	return license, err
}
func (r licenseRepository) Create(c context.Context, o *model.License) (err error) {
	return r.GetDB(c).Create(o).Error
}
func (r licenseRepository) UpdateById(c context.Context, o *model.License, id string) error {
	return r.GetDB(c).Updates(o).Error
}
func (r licenseRepository) DeleteById(c context.Context, id string) (err error) {
	return r.GetDB(c).Where("id = ?", id).Delete(&model.License{}).Error
}
