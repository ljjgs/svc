package main

import (
	"flag"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/util/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/ljjgs/svc/common"
	"github.com/ljjgs/svc/domain/repository"
	service2 "github.com/ljjgs/svc/domain/service"
	"github.com/ljjgs/svc/handler"
	"github.com/ljjgs/svc/proto/svc"
	"github.com/opentracing/opentracing-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strconv"
)

var (
	consulHost       = "162.14.77.55"
	consulPort int64 = 8500
	tracerHost       = "162.14.77.55"
	tracerPort       = 16686
	//hystrixPort          = 9002
	prometheusPort = "9090"
)

func main() {

	dir := homedir.HomeDir()
	print("dir" + dir)
	newRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulHost + ":" + strconv.FormatInt(consulPort, 10),
		}
	})

	config, err := common.GetConsulConfig(consulHost, consulPort, "micro/config")
	if err != nil {
		panic(err)
	}

	mysqlInfo := common.GetMysqlFromConsul(config, "mysql")
	dsn := mysqlInfo.User + ":" + mysqlInfo.Password + "@tcp(" + mysqlInfo.Host + ":" + mysqlInfo.Port + ")/" + mysqlInfo.DataSource + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		log.Errorf("数据库连接成功")
	}
	defer db.Close()
	db.SingularTable(true)

	tracer, closer, err := common.NewTracer("base", tracerHost+":"+strconv.Itoa(tracerPort))
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	//go func() {
	//	err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", strconv.Itoa(hystrixPort)), hystrixStreamHandler)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	//	err = common.PrometheusBoot("162.14.77.55", prometheusPort)
	//	if err != nil {
	//		panic(err)
	//	}
	/**

	 */
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kube_config 在当前系统中")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kube_config 在当前系统中")
	}
	flag.Parse()

	flags, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		panic(err)
	}
	//	//	在集群中使用
	//	forConfig, err := kubernetes.NewForConfig(flags)
	//    if err!=nil {
	//        panic(err)
	//    }
	//   //
	//   forConfig, err := rest.InClusterConfig()
	clientSet, err := kubernetes.NewForConfig(flags)
	if err != nil {
		panic(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.service.svc"),
		micro.Version("latest"),
		micro.Registry(newRegistry),
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(ratelimit.NewClientWrapper(1000)),
	)
	service.Init()

	//err = repository.NewSvcRepository(db).InitTable()

	if err != nil {
		panic(err)
	}

	// 注册句柄，可以快速操作已开发的服务
	svcDataService := service2.NewSvcDataService(repository.NewSvcRepository(db), clientSet)
	svc.RegisterSvcHandler(service.Server(), &handler.SvcHandler{SvcDataService: svcDataService})

	err = service.Run()
	if err != nil {
		panic(err)
	}
}
