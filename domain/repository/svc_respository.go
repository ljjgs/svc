package repository

import (
	"github.com/asim/go-micro/v3/util/log"
	"github.com/jinzhu/gorm"
	"github.com/ljjgs/svc/domain/model"
)

type ISvcRepository interface {
	//初始化表
	InitTable() error
	//根据ID查处找数据s
	FindSvcByID(int64) (*model.Svc, error)
	//创建一条 svc 数据
	CreateSvc(*model.Svc) (int64, error)
	//根据ID删除一条 svc 数据
	DeleteSvcByID(int64) error
	//修改更新数据
	UpdateSvc(*model.Svc) error
	//查找svc所有数据
	FindAll() ([]model.Svc, error)
}

func NewSvcRepository(db *gorm.DB) ISvcRepository {
	return &SvcRepository{
		mysqlDb: db,
	}
}

type SvcRepository struct {
	mysqlDb *gorm.DB
}

func (s *SvcRepository) InitTable() error {
	return s.mysqlDb.CreateTable(&model.Svc{}, &model.SvcPort{}).Error
}

//根据ID查找Svc信息
func (u *SvcRepository) FindSvcByID(svcID int64) (svc *model.Svc, err error) {
	svc = &model.Svc{}
	return svc, u.mysqlDb.First(svc, svcID).Error
}

func (s *SvcRepository) CreateSvc(svc *model.Svc) (int64, error) {
	return svc.ID, s.mysqlDb.Create(svc).Error
}

func (s *SvcRepository) DeleteSvcByID(i int64) error {
	svc := &model.Svc{}
	begin := s.mysqlDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			begin.Rollback()
		}
	}()
	if begin.Error != nil {
		log.Info(begin.Error)
		return nil
	}
	if erro := s.mysqlDb.Where("id = ?", i).Delete(svc).Error; erro != nil {
		begin.Rollback()
		log.Error(erro)
		return erro
	}
	if erro := s.mysqlDb.Where("svc_id = ?", i).Delete(&model.SvcPort{}).Error; erro != nil {
		log.Error(erro)
		return erro
	}

	return nil
}

func (s *SvcRepository) UpdateSvc(svc *model.Svc) error {
	return s.mysqlDb.Model(svc).Update(svc).Error
}

func (s *SvcRepository) FindAll() (svcAll []model.Svc, error error) {
	return svcAll, s.mysqlDb.Find(&svcAll).Error
}
