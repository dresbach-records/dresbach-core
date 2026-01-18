package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Log é a instância global do nosso logger.
var Log = logrus.New()

func init() {
	// Abre ou cria o arquivo de log.
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Fatalf("Falha ao abrir o arquivo de log: %v", err)
	}

	// Configura a saída para ser tanto no arquivo quanto no stdout.
	// Isso é útil para desenvolvimento, permitindo ver os logs no console e tê-los no arquivo.
	Log.SetOutput(io.MultiWriter(os.Stdout, file))

	// Define o formato de saída para JSON. Isso é fundamental para o parsing e análise de logs.
	Log.SetFormatter(&logrus.JSONFormatter{})

	// Define o nível mínimo de log a ser registrado. Info é um bom padrão.
	// Níveis: Debug, Info, Warn, Error, Fatal, Panic.
	Log.SetLevel(logrus.InfoLevel)
}
