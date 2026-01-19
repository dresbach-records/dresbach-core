package main

import (
	"database/sql"
	"net/http"
	"os"

	"hosting-backend/internal/config"
	_ "hosting-backend/docs" // Importa os docs gerados pelo swag
	"hosting-backend/internal/handlers/admin"
	"hosting-backend/internal/handlers/client"
	"hosting-backend/internal/handlers/domain"
	"hosting-backend/internal/handlers/products"
	"hosting-backend/internal/handlers/webhooks"
	"hosting-backend/internal/logger" // Importa nosso novo logger
	"hosting-backend/internal/middleware"
	"hosting-backend/internal/services"
	"hosting-backend/internal/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API de Backend de Hospedagem
// @version 1.0
// @description Este é o backend para um serviço de hospedagem, com gerenciamento de clientes, serviços e faturamento.
// @termsOfService http://swagger.io/terms/
// @contact.name Suporte da API
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()

	logger.Log.Info("Iniciando a API de Hospedagem...")

	if err := utils.InitCrypto(); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"module": "crypto",
			"error":  err.Error(),
		}).Fatal("Erro ao inicializar o módulo de criptografia")
	}

	config.InitAsaas() // Adiciona a inicialização do Asaas

	// A conexão com o banco de dados foi comentada para o teste.
	// db, err := database.ConnectMySQL()
	// if err != nil {
	// 	logger.Log.WithFields(logrus.Fields{
	// 		"module": "database",
	// 		"error":  err.Error(),
	// 	}).Fatal("Erro ao conectar no MySQL")
	// }
	// defer db.Close()
	var db *sql.DB // Declara db como nil para a compilação funcionar.
	logger.Log.Info("Conexão com o banco de dados MySQL pulada para teste.")

	// Instancia os serviços com um db nil.
	clientService := services.NewClientService(db)
	adminService := services.NewAdminService(db)

	// Os workers foram comentados pois causariam pânico com um db nil.
	// go workers.GenerateInvoicesWorker(db)
	// go workers.SuspensionWorker(db)

	r := mux.NewRouter()

	// Aplica o middleware de logging a todas as rotas
	r.Use(middleware.LoggingMiddleware)

	// Rota para a documentação do Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// --- Rotas Públicas ---
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.HandleFunc("/domains/check", domain.CheckDomainHandler()).Methods("POST")
	r.HandleFunc("/products/vps", products.GetVpsProductsHandler()).Methods("GET")

	// --- Rotas de Autenticação ---
	authRouter := r.PathPrefix("/api").Subrouter()
	authRouter.HandleFunc("/register", client.RegisterHandler(db)).Methods("POST")
	authRouter.HandleFunc("/login", client.LoginHandler(db, os.Getenv("JWT_SECRET"))).Methods("POST")

	// --- Rotas da Área do Cliente (Protegidas por JWT) ---
	clientRouter := r.PathPrefix("/").Subrouter()
	clientRouter.Use(middleware.AuthMiddleware)
	clientRouter.HandleFunc("/api/me", client.MeHandler(db)).Methods("GET")
	clientRouter.HandleFunc("/api/my-services", client.GetServicesHandler(clientService)).Methods("GET")
	clientRouter.HandleFunc("/api/my-services/{id:[0-9]+}", client.GetServiceDetailsHandler(db)).Methods("GET")
	clientRouter.HandleFunc("/api/my-invoices", client.GetInvoicesHandler(clientService)).Methods("GET")
	clientRouter.HandleFunc("//api/checkout", client.CheckoutHandler(db)).Methods("POST")

	// --- Rotas de Administração ---
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware, middleware.RBACMiddleware(db), middleware.RequirePermission("admin"))
	adminRouter.HandleFunc("/clients", admin.GetClientsHandler(adminService)).Methods("GET")
	adminRouter.HandleFunc("/clients", admin.CreateClientHandler(adminService)).Methods("POST")
	adminRouter.HandleFunc("/clients/{id:[0-a-z0-9]+}", admin.GetClientHandler(adminService)).Methods("GET")
	adminRouter.HandleFunc("/clients/{id:[0-9]+}", admin.UpdateClientHandler(adminService)).Methods("PUT")
	adminRouter.HandleFunc("/clients/{id:[0-9]+}", admin.DeleteClientHandler(adminService)).Methods("DELETE")
	adminRouter.HandleFunc("/domain-orders", admin.GetDomainOrdersHandler(db)).Methods("GET")
	adminRouter.HandleFunc("/domain-orders/{id:[0-9]+}", admin.UpdateDomainOrderHandler(db)).Methods("PUT")
	adminRouter.HandleFunc("/financials/balance", admin.GetBalanceHandler(db)).Methods("GET")
	adminRouter.HandleFunc("/financials/transactions", admin.GetTransactionsHandler(db)).Methods("GET")

	// Rotas de Configurações Fiscais (Admin)
	adminRouter.HandleFunc("/fiscal/settings", admin.GetFiscalSettingsHandler(db)).Methods("GET")
	adminRouter.HandleFunc("//fiscal/settings", admin.UpdateFiscalSettingsHandler(db)).Methods("PUT")

	// Rotas de gerenciamento de serviços (Admin)
	adminRouter.HandleFunc("/services", admin.CreateServiceHandler(db)).Methods("POST")
	adminRouter.HandleFunc("/services/{id:[0-9]+}/suspend", admin.SuspendServiceHandler(db)).Methods("PUT")
	adminRouter.HandleFunc("/services/{id:[0-9]+}/reactivate", admin.ReactivateServiceHandler(db)).Methods("PUT")

	// Rota de Monitoramento (Admin)
	adminRouter.HandleFunc("/monitoring/logs", admin.GetSystemLogsHandler()).Methods("GET")

	// Rota para Análise de Deploy (Admin)
	adminRouter.HandleFunc("/updates/analyze", admin.AnalyzeUpdateHandler()).Methods("POST")

	// --- Rotas de Webhook (Externas) ---
	r.HandleFunc("/webhooks/asaas", webhooks.AsaasWebhookHandler(db)).Methods("POST")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.WithField("port", port).Info("API rodando na porta")

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Log.WithField("error", err.Error()).Fatal("Erro ao iniciar o servidor HTTP")
	}
}
