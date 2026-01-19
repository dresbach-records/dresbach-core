package deployment

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/docker/go-connections/nat"
)

// ... (structs e constantes permanecem as mesmas)

// AnalysisStatus define se a atualização é segura ou bloqueada.
type AnalysisStatus string

const (
	StatusSafe    AnalysisStatus = "✔️ Atualização Segura"
	StatusBlocked AnalysisStatus = "❌ Atualização Bloqueada"
)

// ProblemSeverity define o nível de criticidade de um problema.
type ProblemSeverity string

const (
	SeverityRisk    ProblemSeverity = "RISCO"
	SeverityBlocker ProblemSeverity = "BLOQUEADOR"
)

// Problem detalha uma inconsistência encontrada durante a análise.
type Problem struct {
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Severity    ProblemSeverity `json:"severity"`
}

// AnalysisReport é o resultado final da análise do Backend Analyzer.
type AnalysisReport struct {
	Version           string         `json:"version"`
	Status            AnalysisStatus `json:"status"`
	Problems          []Problem      `json:"problems"`
	RecommendedAction string         `json:"recommended_action"`
}

// Analyzer é o orquestrador do processo de análise de deploy.
type Analyzer struct {
	ArtifactURL   string
	ArtifactHash  string
	Version       string
	extractedPath string
}

// NewAnalyzer cria uma nova instância do analisador.
func NewAnalyzer(artifactURL, artifactHash, version string) *Analyzer {
	return &Analyzer{
		ArtifactURL:   artifactURL,
		ArtifactHash:  artifactHash,
		Version:       version,
		extractedPath: ".", // Simulação
	}
}

func (a *Analyzer) RunChecks() (*AnalysisReport, error) {
	report := &AnalysisReport{
		Version: a.Version,
		Status:  StatusSafe,
	}

	checks := []func() ([]Problem, error){
		a.checkEnvVariables,
		a.checkDatabaseMigrations,
		a.checkDependencies,
		a.checkSandboxHealth, // Verificação no Sandbox
	}

	for _, check := range checks {
		problems, err := check()
		if err != nil {
			return nil, fmt.Errorf("falha ao executar a verificação: %w", err)
		}
		if problems != nil {
			report.Problems = append(report.Problems, problems...)
		}
	}

	if len(report.Problems) > 0 {
		report.Status = StatusBlocked
		report.RecommendedAction = "Corrija os problemas bloqueadores listados antes de tentar o deploy novamente."
	}

	return report, nil
}

// --- Seção de Funções de Verificação ---

// ... (checkEnvVariables, checkDatabaseMigrations, checkDependencies permanecem os mesmos)
func (a *Analyzer) checkEnvVariables() ([]Problem, error) {
	var problems []Problem
	exampleFilePath := filepath.Join(a.extractedPath, ".env.example")

	requiredVars, err := parseEnvFile(exampleFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Não é um bloqueio se o .env.example não existir
		}
		return nil, err
	}

	for _, key := range requiredVars {
		if os.Getenv(key) == "" {
			problems = append(problems, Problem{
				Code:        "ENV_VAR_MISSING",
				Description: "Variável de ambiente obrigatória '" + key + "' não está definida.",
				Severity:    SeverityBlocker,
			})
		}
	}
	return problems, nil
}

func (a *Analyzer) checkDatabaseMigrations() ([]Problem, error) {
	var problems []Problem
	migrationsPath := filepath.Join(a.extractedPath, "migrations")
	destructiveKeywords := []string{"DROP TABLE", "DROP COLUMN", "TRUNCATE TABLE"}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return nil, nil // Nenhum diretório de migrations, nada a fazer.
	}

	err := filepath.Walk(migrationsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
			content, _ := os.ReadFile(path)
			textContent := strings.ToUpper(string(content))
			for _, keyword := range destructiveKeywords {
				if strings.Contains(textContent, keyword) {
					problems = append(problems, Problem{
						Code:        "DB_DESTRUCTIVE_MIGRATION",
						Description: "Comando destrutivo '" + keyword + "' encontrado em: " + info.Name(),
						Severity:    SeverityBlocker,
					})
				}
			}
		}
		return nil
	})

	return problems, err
}

type GovulncheckFinding struct {
	OSV   string `json:"osv"`
	Fixed string `json:"fixed"`
	Trace []struct {
		Module string `json:"module"`
	} `json:"trace"`
}
type GovulncheckOutput struct {
	Finding GovulncheckFinding `json:"finding"`
}

func (a *Analyzer) checkDependencies() ([]Problem, error) {
	var problems []Problem
	cmd := exec.Command("govulncheck", "-json", "./...")
	cmd.Dir = a.extractedPath

	var out bytes.Buffer
	cmd.Stdout = &out

	_ = cmd.Run()

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		var line GovulncheckOutput
		if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
			continue
		}

		if line.Finding.OSV != "" {
			finding := line.Finding
			problem := Problem{
				Code: "DEPENDENCY_VULNERABILITY",
				Description: fmt.Sprintf(
					"Vulnerabilidade [%s] encontrada no módulo '%s'. Versão corrigida: %s.",
					finding.OSV,
					finding.Trace[0].Module,
					finding.Fixed,
				),
				Severity: SeverityBlocker,
			}
			problems = append(problems, problem)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler a saída do govulncheck: %w", err)
	}

	return problems, nil
}

func (a *Analyzer) checkSandboxHealth() ([]Problem, error) {
	ctx := context.Background()
	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []Problem{{Code: "DOCKER_UNAVAILABLE", Description: "Não foi possível conectar ao Docker daemon."}}, nil
	}

	imageName := fmt.Sprintf("sandbox-test-%s:latest", strings.ToLower(a.Version))

	// Garante que a imagem e o contêiner serão removidos no final.
	defer a.cleanupSandbox(ctx, dockerCli, imageName, "")

	// Constrói a imagem Docker
	buildResponse, err := a.buildDockerImage(ctx, dockerCli, imageName)
	if err != nil {
		return []Problem{{Code: "DOCKER_BUILD_FAILED", Description: fmt.Sprintf("Falha ao construir a imagem Docker: %v", err)}}, nil
	}
	defer buildResponse.Body.Close()
	// É importante consumir o output do build para garantir que ele finalize
	io.Copy(io.Discard, buildResponse.Body)

	// Cria e inicia o contêiner
	containerID, exposedPort, err := a.runContainer(ctx, dockerCli, imageName)
	if err != nil {
		return []Problem{{Code: "DOCKER_RUN_FAILED", Description: fmt.Sprintf("Falha ao iniciar o contêiner de sandbox: %v", err)}}, nil
	}
	// Atualiza o defer para garantir a remoção do contêiner criado
	defer a.cleanupSandbox(ctx, dockerCli, "", containerID)

	// Executa o Health Check
	healthy, err := a.performHealthCheck(exposedPort)
	if err != nil || !healthy {
		return []Problem{{Code: "HEALTH_CHECK_FAILED", Description: "A aplicação no sandbox falhou no health check."}}, nil
	}

	return nil, nil // Tudo certo!
}

func (a *Analyzer) buildDockerImage(ctx context.Context, cli *client.Client, imageName string) (types.ImageBuildResponse, error) {
	return cli.ImageBuild(ctx, nil, types.ImageBuildOptions{
		Context:    a.extractedPath, // Usa o diretório do artefato como contexto
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true, // Remove contêineres intermediários
	})
}

func (a *Analyzer) runContainer(ctx context.Context, cli *client.Client, imageName string) (string, string, error) {
	config := &container.Config{
		Image:        imageName,
		ExposedPorts: nat.PortSet{"8080/tcp": {}},
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "0"}}, // "0" para porta aleatória
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", err
	}

	// Inspeciona o contêiner para descobrir a porta mapeada
	inspect, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", err
	}
	exposedPort := inspect.NetworkSettings.Ports["8080/tcp"][0].HostPort

	return resp.ID, exposedPort, nil
}

func (a *Analyzer) performHealthCheck(port string) (bool, error) {
	// Tenta o health check por até 15 segundos
	for i := 0; i < 15; i++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
		if err == nil && resp.StatusCode == http.StatusOK {
			return true, nil
		}
		time.Sleep(1 * time.Second)
	}
	return false, fmt.Errorf("health check timed out")
}

func (a *Analyzer) cleanupSandbox(ctx context.Context, cli *client.Client, imageName, containerID string) {
	if containerID != "" {
		// Para o contêiner, não se importa com erros
		_ = cli.ContainerStop(ctx, containerID, container.StopOptions{})
		// Remove o contêiner
		_ = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	}
	if imageName != "" {
		// Remove a imagem
		_, _ = cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{Force: true})
	}
}

// --- Funções Auxiliares ---
func parseEnvFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if parts := strings.SplitN(line, "=", 2); len(parts) > 0 {
			keys = append(keys, parts[0])
		}
	}
	return keys, scanner.Err()
}
