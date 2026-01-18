package admin

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"

	"hosting-backend/internal/logger"

	"github.com/sirupsen/logrus"
)

// LogEntry representa uma única entrada de log para ser exibida na UI.
type LogEntry struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp string    `json:"timestamp"`
	Module    string    `json:"module,omitempty"`
	Method    string    `json:"method,omitempty"`
	URI       string    `json:"uri,omitempty"`
	Status    int       `json:"status,omitempty"`
	Latency   int64     `json:"latency_ms,omitempty"`
	RemoteIP  string    `json:"remote_addr,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// GetSystemLogsHandler lê e retorna os logs do sistema.
// @Summary Obtém os logs do sistema
// @Description Retorna uma lista de todas as entradas de log do sistema para monitoramento.
// @Tags Admin
// @Produce json
// @Success 200 {array} LogEntry
// @Failure 500 {object} map[string]string
// @Router /admin/monitoring/logs [get]
// @Security ApiKeyAuth
func GetSystemLogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Em um cenário real, os logs viriam de um arquivo, um serviço centralizado
		// (como ELK, Loki, etc.) ou de uma captura em memória.
		// Para este exemplo, vamos simular a leitura de um arquivo de log.
		logs, err := readLogsFromFile("app.log") // Simulação

		if err != nil {
			logger.Log.WithField("error", err).Error("Falha ao ler o arquivo de log")
			http.Error(w, `{"error": "Falha ao ler os logs do sistema"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// readLogsFromFile é uma função auxiliar para simular a leitura e o parsing de um arquivo de log.
func readLogsFromFile(filePath string) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		// Se o arquivo não existir, retorna uma lista vazia, o que é normal se a app acabou de iniciar.
		if os.IsNotExist(err) {
			return []LogEntry{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var rawLog map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &rawLog); err != nil {
			// Ignora linhas que não são JSON válido.
			continue
		}

		// Mapeia o log JSON bruto para nossa estrutura LogEntry.
		entry := LogEntry{
			Level:     rawLog["level"].(string),
			Message:   rawLog["msg"].(string), // logrus usa 'msg' para a mensagem
			Timestamp: rawLog["time"].(string), // logrus usa 'time'
		}

		// Adiciona campos opcionais se eles existirem no log.
		if module, ok := rawLog["module"]; ok {
			entry.Module = module.(string)
		}
		if method, ok := rawLog["method"]; ok {
			entry.Method = method.(string)
		}
		if uri, ok := rawLog["uri"]; ok {
			entry.URI = uri.(string)
		}
		if status, ok := rawLog["status"].(float64); ok { // JSON decodifica números como float64
			entry.Status = int(status)
		}
		if latency, ok := rawLog["latency_ms"].(float64); ok {
			entry.Latency = int64(latency)
		}
		if remoteIP, ok := rawLog["remote_addr"]; ok {
			entry.RemoteIP = remoteIP.(string)
		}
		if err, ok := rawLog["error"]; ok {
			entry.Error = err.(string)
		}

		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}
