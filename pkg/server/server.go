package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
	"github.com/gorilla/mux"
)

type MetricServer struct {
	storage storage.MetricStorage
	router  *mux.Router
}

func NewMetricServer(storage storage.MetricStorage) *MetricServer {
	return &MetricServer{
		storage: storage,
		router:  mux.NewRouter(),
	}
}

func (s *MetricServer) handleMetricUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	metricType := vars["type"]
	metricName := vars["name"]
	metricValue := vars["value"]

	if metricName == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	switch metricType {
	case "gauge":
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		s.storage.AddGaugeMetric(metricName, floatValue)
		// result, _ := s.storage.GetGaugeMetric(metricName)
		// fmt.Println(result)
	case "counter":
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		s.storage.AddCounterMetric(metricName, intValue)
		// result, _ := s.storage.GetCounterMetric(metricName)
		// fmt.Println(result)
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *MetricServer) handleMetricValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	metricType := vars["type"]
	metricName := vars["name"]

	if metricName == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var metricValue string

	switch metricType {
	case "gauge":
		val, ok := s.storage.GetGaugeMetric(metricName)
		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		metricValue = strconv.FormatFloat(val, 'f', -1, 64)
	case "counter":
		val, ok := s.storage.GetCounterMetric(metricName)
		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		metricValue = strconv.FormatInt(val, 10)
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricValue))
}

func (s *MetricServer) handleMetricsList(w http.ResponseWriter, r *http.Request) {
	metrics := s.storage.GetAllMetrics()

	var html strings.Builder
	html.WriteString("<html><body><h1>Metric List</h1><ul>")

	for name, value := range metrics {
		var stringValue string
		switch v := value.(type) {
		case float64:
			stringValue = strconv.FormatFloat(v, 'f', -1, 64)
		case int64:
			stringValue = strconv.FormatInt(v, 10)
		}
		html.WriteString("<li>" + name + ": " + stringValue + "</li>")
	}
	html.WriteString("</ul></body></html>")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html.String()))
}

func (s *MetricServer) Start(addr string) error {
	s.router.HandleFunc("/update/{type}/{name}/{value}", s.handleMetricUpdate).Methods("POST")
	s.router.HandleFunc("/value/{type}/{name}", s.handleMetricValue).Methods("GET")
	s.router.HandleFunc("/", s.handleMetricsList).Methods("GET")
	fmt.Printf("Server is listening on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
