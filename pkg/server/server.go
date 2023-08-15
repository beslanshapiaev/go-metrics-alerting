// pkg/server/server.go

package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/beslanshapiaev/go-metrics-alerting/common"
	"github.com/beslanshapiaev/go-metrics-alerting/internal/middleware"
	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
	"github.com/jackc/pgx/v5"

	"github.com/gorilla/mux"
)

// test
type MetricServer struct {
	storage       storage.MetricStorage
	router        *mux.Router
	storeInterval time.Duration
}

var connectionString string

func NewMetricServer(storage storage.MetricStorage, connString string) *MetricServer {
	connectionString = connString
	return &MetricServer{
		storage: storage,
		router:  mux.NewRouter(),
	}
}

func (s *MetricServer) SetStoreInterval(interval time.Duration) {
	s.storeInterval = interval
}

func (s *MetricServer) handleMetricUpdate(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		s.handleMetricUpdateJSON(w, r)
	} else {
		s.handleMetricUpdateForm(w, r)
	}
}

func (s *MetricServer) handleMetricUpdates(w http.ResponseWriter, r *http.Request) {
	var metrics []common.Metric

	var reader io.Reader = r.Body

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()

		var buf bytes.Buffer
		teeReader := io.TeeReader(gzipReader, &buf)

		data, err := io.ReadAll(teeReader)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		reader = &buf
	}
	err := json.NewDecoder(reader).Decode(&metrics)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = s.storage.AddMetricsBatch(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *MetricServer) handleMetricUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var metric common.Metric

	var reader io.Reader = r.Body

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()

		var buf bytes.Buffer
		teeReader := io.TeeReader(gzipReader, &buf)

		data, err := io.ReadAll(teeReader)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		reader = &buf
	}

	err := json.NewDecoder(reader).Decode(&metric)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":
		s.storage.AddGaugeMetric(metric.ID, *metric.Value)
	case "counter":
		s.storage.AddCounterMetric(metric.ID, *metric.Delta)
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	updatedMetric, err := s.getMetricValue(metric.ID, metric.MType)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMetric)
}

func (s *MetricServer) handleMetricUpdateForm(w http.ResponseWriter, r *http.Request) {
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

func (s *MetricServer) getMetricValue(metricID, metricType string) (*common.Metric, error) {
	var metricValue interface{}
	var ok bool

	switch metricType {
	case "gauge":
		metricValue, ok = s.storage.GetGaugeMetric(metricID)
		if !ok {
			return nil, fmt.Errorf("metric not found: %s", metricID)
		}
	default:
		metricValue, ok = s.storage.GetCounterMetric(metricID)
		if !ok {
			return nil, fmt.Errorf("metric not found: %s", metricID)
		}
		value := metricValue.(int64)
		metricValue = value
	}

	metric := &common.Metric{
		ID:    metricID,
		MType: metricType,
	}

	switch v := metricValue.(type) {
	case float64:
		metric.Value = &v
	case int64:
		metric.Delta = &v
	default:
		return nil, fmt.Errorf("unsupported metric value type")
	}

	return metric, nil
}

func (s *MetricServer) handleMetricValue(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		s.handleMetricValueJSON(w, r)
	} else {
		s.handleMetricValueForm(w, r)
	}
}

func (s *MetricServer) handleMetricValueJSON(w http.ResponseWriter, r *http.Request) {
	var metric common.Metric
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if metric.ID == "" || metric.MType == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	metricValue, err := s.getMetricValue(metric.ID, metric.MType)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricValue)
}

func (s *MetricServer) handleMetricValueForm(w http.ResponseWriter, r *http.Request) {
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

func (s *MetricServer) handlePing(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect database %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())
	if err = conn.Ping(context.Background()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *MetricServer) Start(addr string) error {
	s.router.Use(middleware.LoggingMiddleware)
	s.router.Use(middleware.GzipMiddleware)
	s.router.HandleFunc("/update/{type}/{name}/{value}", s.handleMetricUpdate).Methods("POST")
	s.router.HandleFunc("/update/", s.handleMetricUpdate).Methods("POST")
	s.router.HandleFunc("/updates/", s.handleMetricUpdates).Methods("POST")
	s.router.HandleFunc("/value/{type}/{name}", s.handleMetricValue).Methods("GET")
	s.router.HandleFunc("/value/", s.handleMetricValue).Methods("POST")
	s.router.HandleFunc("/", s.handleMetricsList).Methods("GET")
	s.router.HandleFunc("/ping", s.handlePing).Methods("GET")

	fmt.Printf("Server is listening on %s\n", addr)

	if s.storeInterval > 0 {
		ticker := time.NewTicker(s.storeInterval)
		go func() {
			for {
				<-ticker.C
				if err := s.storage.SaveToFile(); err != nil {
					fmt.Printf("Error saving metrics to file: %v\n", err)
				}
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		if err := s.storage.SaveToFile(); err != nil {
			fmt.Printf("Error saving metrics to file during graceful shutdown: %v\n", err)
		}
		os.Exit(0)
	}()

	return http.ListenAndServe(addr, s.router)
}
