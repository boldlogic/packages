package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Registry описывает минимальный контракт Prometheus registry,
// который умеет регистрировать и собирать метрики.
type Registry interface {
	prometheus.Registerer
	prometheus.Gatherer
}

// New создаёт новый Prometheus registry и регистрирует в нём
// стандартные Go- и process-коллекторы.
func New() *prometheus.Registry {
	reg := prometheus.NewRegistry()

	_ = reg.Register(collectors.NewGoCollector())
	_ = reg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return reg
}
