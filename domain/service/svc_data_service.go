package service

import (
	"context"
	"errors"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/ljjgs/svc/domain/model"
	"github.com/ljjgs/svc/domain/repository"
	"github.com/ljjgs/svc/proto/svc"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type ISvcDataService interface {
	AddSvc(*model.Svc) (int64, error)
	DeleteSvc(int64) error
	UpdateSvc(*model.Svc) error
	FindSvcByID(int64) (*model.Svc, error)
	FindAllSvc() ([]model.Svc, error)
	CreateSvcToK8s(*svc.SvcInfo) error
	UpdateSvcToK8s(*svc.SvcInfo) error
	DeleteFromK8s(*model.Svc) error
}

type SvcDataService struct {
	SvcRepository repository.ISvcRepository
	K8sClientSet  *kubernetes.Clientset
}

func (s SvcDataService) AddSvc(m *model.Svc) (int64, error) {
	return s.SvcRepository.CreateSvc(m)
}

func (s SvcDataService) DeleteSvc(i int64) error {
	return s.SvcRepository.DeleteSvcByID(i)
}

func (s SvcDataService) UpdateSvc(m *model.Svc) error {

	return s.SvcRepository.UpdateSvc(m)
}

func (s SvcDataService) FindSvcByID(i int64) (*model.Svc, error) {
	return s.SvcRepository.FindSvcByID(i)
}

func (s SvcDataService) FindAllSvc() ([]model.Svc, error) {
	return s.SvcRepository.FindAll()
}

func (s SvcDataService) CreateSvcToK8s(info *svc.SvcInfo) (err error) {
	service := s.setService(info)
	//查找是否纯在指定的服务
	if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Get(context.TODO(), info.SvcName, v12.GetOptions{}); err != nil {
		//查找不到,就创建
		if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Create(context.TODO(), service, v12.CreateOptions{}); err != nil {
			log.Error(err)
			return err
		}
		return nil
	} else {
		log.Error("Service " + info.SvcName + "已经存在")
		return errors.New("Service " + info.SvcName + "已经存在")
	}

}

func (s SvcDataService) UpdateSvcToK8s(info *svc.SvcInfo) (err error) {
	service := s.setService(info)
	//查找是否纯在指定的服务
	if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Get(context.TODO(), info.SvcName, v12.GetOptions{}); err != nil {
		//查找不到
		log.Error(err)
		return errors.New("Service" + info.SvcName + "不存在请先创建")
	} else {
		if _, err = s.K8sClientSet.CoreV1().Services(info.SvcNamespace).Update(context.TODO(), service, v12.UpdateOptions{}); err != nil {
			log.Error(err)
			return err
		}
		log.Info("Service " + info.SvcName + "更新成功")
		return nil
	}
}

func (s SvcDataService) DeleteFromK8s(m *model.Svc) (err error) {
	if err = s.K8sClientSet.CoreV1().Services(m.SvcNamespace).Delete(context.TODO(), m.SvcName, v12.DeleteOptions{}); err != nil {
		log.Error(err)
		return err
	} else {
		if err := s.DeleteSvc(m.ID); err != nil {
			log.Error(err)
			return err
		}
		log.Info("删除Service ID：" + strconv.FormatInt(m.ID, 10) + "成功！")
	}
	return
}

func NewSvcDataService(svcRepository repository.ISvcRepository, clientset *kubernetes.Clientset) ISvcDataService {
	return &SvcDataService{
		SvcRepository: svcRepository,
		K8sClientSet:  clientset,
	}
}

func (s SvcDataService) setService(svcInfo *svc.SvcInfo) *v1.Service {
	service := &v1.Service{}
	//设置服务类型
	service.TypeMeta = v12.TypeMeta{
		Kind:       "v1",
		APIVersion: "Service",
	}
	//设置服务基础信息
	service.ObjectMeta = v12.ObjectMeta{
		Name:      svcInfo.SvcName,
		Namespace: svcInfo.SvcNamespace,
		Labels: map[string]string{
			"app-name": svcInfo.SvcPodName,
			"author":   "Caplost",
		},
		Annotations: map[string]string{
			"k8s/generated-by-cap": "由Cap老师代码创建",
		},
	}
	//设置服务的spec信息，课程中采用ClusterIP模式
	service.Spec = v1.ServiceSpec{
		Ports: s.getSvcPort(svcInfo),
		Selector: map[string]string{
			"app-name": svcInfo.SvcPodName,
		},
		Type: "ClusterIP",
	}
	return service
}
func (p SvcDataService) getSvcPort(svcInfo *svc.SvcInfo) (servicePort []v1.ServicePort) {
	for _, v := range svcInfo.SvcPort {
		servicePort = append(servicePort, v1.ServicePort{
			Name:       "port-" + strconv.FormatInt(int64(v.SvcPort), 10),
			Protocol:   v1.Protocol(v.SvcPortProtocol),
			Port:       v.SvcPort,
			TargetPort: intstr.FromInt(int(v.SvcTargetPort)),
		})
	}
	return
}
