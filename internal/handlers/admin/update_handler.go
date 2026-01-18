package admin

import (
	"encoding/json"
	"net/http"

	"hosting-backend/internal/deployment"
	"hosting-backend/internal/logger"
)

// ArtifactPayload define a estrutura para receber a localização do artefato de deploy.
type ArtifactPayload struct {
	// A URL ou caminho para o artefato a ser analisado (ex: .tar.gz de um build de CI).
	ArtifactURL string `json:"artifact_url"`
	// O hash SHA256 do artefato para verificação de integridade.
	ArtifactHash string `json:"artifact_hash"`
	// A versão que está sendo proposta para o deploy.
	Version string `json:"version"`
}

// AnalyzeUpdateHandler é o ponto de entrada para o sistema de deploy controlado.
// Ele recebe um artefato, aciona o Backend Analyzer e retorna um relatório de compatibilidade.
// @Summary Analisa uma nova versão para deploy
// @Description Recebe um artefato, executa uma série de validações (variáveis de ambiente, migrations, dependências) e retorna um relatório indicando se o deploy é seguro.
// @Tags Admin
// @Accept json
// @Produce json
// @Param payload body ArtifactPayload true "Informações do Artefato"
// @Success 200 {object} deployment.AnalysisReport "Relatório de análise indicando se a atualização é segura ou bloqueada."
// @Failure 400 {object} map[string]string "Payload inválido."
// @Failure 500 {object} map[string]string "Erro interno durante a análise."
// @Router /admin/updates/analyze [post]
// @Security ApiKeyAuth
func AnalyzeUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ArtifactPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error": "Payload da requisição inválido"}`, http.StatusBadRequest)
			return
		}

		// Aqui chamaremos o módulo central do Analyzer.
		// Por enquanto, vamos retornar um relatório mockado.
		analyzer := deployment.NewAnalyzer(payload.ArtifactURL, payload.ArtifactHash, payload.Version)
		report, err := analyzer.RunChecks()

		if err != nil {
			logger.Log.Errorf("Erro ao executar a análise de deploy: %v", err)
			http.Error(w, `{"error": "Ocorreu um erro inesperado durante a análise"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	}
}
