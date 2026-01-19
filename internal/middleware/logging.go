package middleware

import (
	"net/http"
	"time"

	"hosting-backend/internal/logger"

	"github.com/sirupsen/logrus"
)

// LoggingMiddleware registra informações sobre cada requisição recebida.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Usamos um ResponseWriter customizado para capturar o status da resposta.
		lrw := newLoggingResponseWriter(w)

		// Chama o próximo handler na cadeia.
		next.ServeHTTP(lrw, r)

		latency := time.Since(startTime)
		statusCode := lrw.statusCode

		// Cria uma entrada de log estruturada com todos os detalhes da requisição.
		logger.Log.WithFields(logrus.Fields{
			"module":      "http_server",
			"method":      r.Method,
			"uri":         r.RequestURI,
			"remote_addr": r.RemoteAddr,
			"status":      statusCode,
			"latency_ms":  latency.Milliseconds(),
		}).Info("HTTP request completed") // Mensagem principal do log
	})
}

// loggingResponseWriter é um wrapper em torno de http.ResponseWriter que captura o status code.

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// O status padrão, caso WriteHeader não seja chamado, é http.StatusOK (200).
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
