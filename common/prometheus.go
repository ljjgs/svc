package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func PrometheusBoot(host, port string) error {
	promauto.NewCounter(prometheus.CounterOpts{
		Name: "my_counter",
		Help: "This is my counter",
	})
	// 导出指标
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(host+":"+port, nil); err != nil {
			panic(err)
		}
	}()

	return nil
}
