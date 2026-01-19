package admin

import (
	"encoding/json"
	"net/http"

	"hosting-backend/internal/logger"
)

// ArtifactPayload defines the structure for receiving the deployment artifact location.
type ArtifactPayload struct {
	// The URL or path to the artifact to be analyzed (e.g., .tar.gz from a CI build).
	ArtifactURL string `json:"artifact_url"`
	// The SHA256 hash of the artifact for integrity verification.
	ArtifactHash string `json:"artifact_hash"`
	// The version being proposed for deployment.
	Version string `json:"version"`
}

// AnalyzeUpdateHandler is the entry point for the controlled deployment system.
// It receives an artifact, triggers the Backend Analyzer, and returns a compatibility report.
// @Summary Analyzes a new version for deployment
// @Description Receives an artifact, runs a series of validations (environment variables, migrations, dependencies), and returns a report indicating if the deployment is safe.
// @Tags Admin
// @Accept json
// @Produce json
// @Param payload body ArtifactPayload true "Artifact Information"
// @Success 200 {object} map[string]string "Mocked analysis report for now."
// @Failure 400 {object} map[string]string "Invalid payload."
// @Failure 500 {object} map[string]string "Internal error during analysis."
// @Router /admin/updates/analyze [post]
// @Security ApiKeyAuth
func AnalyzeUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ArtifactPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		// The analyzer module is currently disabled. 
		// We will return a mocked report.
		logger.Log.Infof("Received analysis request for version %s. Returning mocked report.", payload.Version)

		mockedReport := map[string]string{
			"status": "Analysis Disabled",
			"message": "The deployment analysis module is temporarily disabled.",
			"version": payload.Version,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockedReport)
	}
}
