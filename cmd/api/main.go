package main

import (
	"net/http"
	"os"

	"hosting-backend/internal/config"
	"hosting-backend/internal/database"
	_ "hosting-backend/docs" // Importa os docs gerados pelo swag
	"hosting-backend/internal/handlers/admin"
	"hosting-backend/internal/handlers/client"
	"hosting-backend/internal/handlers/domain"
	"hosting-backend/internal/handlers/products"
	"hosting-backend/internal/handlers/webhooks"
	"hosting-backend/internal/logger" // Importa nosso novo logger
	"hosting-backend/internal/middleware"
	"hosting-backend/internal/workers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/sirupsen/logrus"
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

	config.InitStripe()

	db, err := database.ConnectMySQL()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"module": "database",
			"error":  err.Error(),
		}).Fatal("Erro ao conectar no MySQL")
	}
	defer db.Close()
	logger.Log.Info("Conexão com o banco de dados MySQL estabelecida com sucesso.")

	go workers.GenerateInvoicesWorker(db)
	go workers.SuspensionWorker(db)

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
	clientRouter.Use(middleware.JWTAuthMiddleware)
	clientRouter.HandleFunc("/api/me", client.MeHandler(db)).Methods("GET")
	clientRouter.HandleFunc("/api/my-services", client.ListMyServicesHandler(db)).Methods("GET")
	clientRouter.HandleFunc("/api/my-services/{id:[0-9]+}", client.GetServiceDetailsHandler(db)).Methods("GET")
	clientRouter.HandleFunc("/api/my-invoices", client.ListMyInvoicesHandler(db)).Methods("GET")

	// --- Rotas de Administração ---
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.JWTAuthMiddleware, middleware.RBACMiddleware(db), middleware.RequirePermission("admin"))
	adminRouter.HandleFunc("/clients", admin.GetClientsHandler(db)).Methods("GET")
	adminRouter.HandleFunc("/clients", admin.CreateClientHandler(db)).Methods("POST")
	adminRouter.HandleFunc("/clients/{id:[0-a-z0-9]+}", admin.GetClientHandler(db)).Methods("GET")
	adminRouter.HandleFunc("/clients/{id:[0-9]+}", admin.UpdateClientHandler(db)).Methods("PUT")
	adminRouter.HandleFunc("/clients/{id:[0-9]+}", admin.DeleteClientHandler(db)).Methods("DELETE")
	adminRouter.HandleFunc("/domain-orders", admin.GetDomainOrdersHandler(db)).Methods("GET")
	adminRouter.HandleFunc("/domain-orders/{id:[0-9]+}", admin.UpdateDomainOrderHandler(db)).Methods("PUT")
	adminRouter.HandleFunc("/financials/balance", admin.GetBalanceHandler(db)).Methods("GET")
	adminRouter.HandleFunc("/financials/transactions", admin.GetTransactionsHandler(db)).Methods("GET")

	// Rotas de gerenciamento de serviços (Admin)
	adminRouter.HandleFunc("/services", admin.CreateServiceHandler(db)).Methods("POST")
	adminRouter.HandleFunc("/services/{id:[0-9]+}/suspend", admin.SuspendServiceHandler(db)).Methods("PUT")
	adminRouter.HandleFunc("/services/{id:[0-9]+}/reactivate", admin.ReactivateServiceHandler(db)).Methods("PUT")

	// Rota de Monitoramento (Admin)
	adminRouter.HandleFunc("/monitoring/logs", admin.GetSystemLogsHandler()).Methods("GET")

	// Rota para Análise de Deploy (Admin)
	adminRouter.HandleFunc("/updates/analyze", admin.AnalyzeUpdateHandler()).Methods("POST")

	// --- Rotas de Webhook (Externas) ---
	r.HandleFunc("/webhooks/stripe", webhooks.StripeWebhookHandler(db)).Methods("POST")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.WithField("port", port).Info("API rodando na porta")

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Log.WithField("error", err.Error()).Fatal("Erro ao iniciar o servidor HTTP")
	}
}
