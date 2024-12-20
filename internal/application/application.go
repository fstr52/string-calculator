package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fstr52/calculator/internal/config"
	"github.com/fstr52/calculator/pkg/calculation"
)

type Config struct {
	Addr string
}

// Полчение конфига из переменной окружения (в конфиге есть возможность изменять порт запуска сервера)
// файл конфигурации (конфиг) создается в "../../config.txt", так что запуск приложения из редактора/IDE нужно выполнять СТРОГО ПО ПУТИ C://User/.../calc/cmd/app (прописано в README)
func ConfigFromEnv() *Config {
	rootDir := GetProjectRoot()
	configPath := filepath.Join(rootDir, "config.txt")
	config.GetConfig(configPath)
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	return config
}

// Получение корневой папки проекта
func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

type Logger struct {
	logger  *slog.Logger
	logFile *os.File
}

// Открытие/создание файла для хранения логов
func OpenLogger() *Logger {
	rootDir := GetProjectRoot()
	loggerPath := filepath.Join(rootDir, "log.txt")
	file, err := os.OpenFile(loggerPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(file, nil)
	logger := slog.New(handler)
	return &Logger{logger: logger, logFile: file}
}

func (a *Application) CloseLogger() {
	a.logger.logFile.Close()
}

type Application struct {
	config *Config
	logger *Logger
}

// Создание нового application-а с использованием config и logger
func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
		logger: OpenLogger(),
	}
}

// Запуск приложения в консоли/редакторе (НЕ СЕРВЕР)
func (a *Application) Run() error {
	logger := a.logger.logger
	for {
		logger.Info("expression inputing")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("failed to read expression",
				slog.Any("error", err),
			)
			fmt.Println("failed to read expression")
			return err
		}

		text = strings.TrimSpace(text)
		if text == "exit" {
			logger.Info("application closed succesfully")
			return nil
		}

		logger.Info("succesfull reading of expression",
			slog.String("expression", text),
		)

		result, err := calculation.Calculate(text)
		if err != nil {
			logger.Error("failed to calculate expression",
				slog.Any("error", err),
			)
			fmt.Println("failed to calculate expression")
		}

		fmt.Println("The result is:", result)
	}
}

type Request struct {
	Expression string `json:"expression"`
}

// Хэндлер для логирования
func (a *Application) LoggingHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := a.logger.logger
		timeStart := time.Now()
		logger.Info("new request",
			slog.Group("request info",
				slog.String("url", r.URL.String()),
				slog.Any("header", r.Header),
				slog.String("method", r.Method),
			),
		)

		next.ServeHTTP(w, r)

		duration := time.Since(timeStart)

		logger.Info("response",
			slog.Duration("duration", duration),
		)
	})
}

// Хэндлер для обработки метода запроса
func (a *Application) RequestHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		logger := a.logger.logger
		if method != http.MethodPost {
			response := response{Err: "Method not allowed"}
			responseJson, err := json.Marshal(&response)
			logger.Info("new request with bad method",
				slog.Group("request info",
					slog.String("url", r.URL.String()),
					slog.Any("header", r.Header),
					slog.String("method", r.Method),
				),
			)
			if err != nil {
				logger.Error("error encoding response with err",
					slog.Any("error", err),
				)
			}
			http.Error(w, string(responseJson), http.StatusMethodNotAllowed)
		} else if r.Header.Get("Content-Type") != "application/json" {
			response := response{Err: "Expression is not valid. Expected JSON format input"}
			responseJson, err := json.Marshal(&response)
			logger.Info("new request with bad Content-Type",
				slog.Group("request info",
					slog.String("url", r.URL.String()),
					slog.Any("header", r.Header),
					slog.String("method", r.Method),
				),
			)
			if err != nil {
				logger.Error("error encoding response with err",
					slog.Any("error", err),
				)
			}
			http.Error(w, string(responseJson), http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type response struct {
	Result string `json:"result,omitempty"`
	Err    string `json:"error,omitempty"`
}

// Хэндлер для обработки выражения и вывода ответа (СЕРВЕР)
func (a *Application) CalcHandler(w http.ResponseWriter, r *http.Request) {
	logger := a.logger.logger
	response := new(response)
	request := new(Request)

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Error("request body decode error",
			slog.Any("error", err),
		)
		return
	}

	result, err := calculation.Calculate(request.Expression)
	if err != nil {
		if err == calculation.ErrInvalidExpression || err == calculation.ErrDivisionByZero {
			response.Err = "Expression is not valid"
			responseJson, err := json.Marshal(&response)
			if err != nil {
				logger.Error("error encoding response with err",
					slog.Any("error", err),
				)
			}
			http.Error(w, string(responseJson), http.StatusUnprocessableEntity)
		} else {
			response.Err = "Internal server error"
			responseJson, err := json.Marshal(&response)
			if err != nil {
				logger.Error("error encoding response with err",
					slog.Any("error", err),
				)
			}
			http.Error(w, string(responseJson), http.StatusInternalServerError)
		}
	} else {
		response.Result = strconv.FormatFloat(result, 'f', -1, 64)
		responseJson, err := json.Marshal(&response)
		if err != nil {
			logger.Error("error encoding response with err",
				slog.Any("error", err),
			)
		}
		fmt.Fprint(w, string(responseJson))
	}
}

// Запуск сервера
func (a *Application) RunServer() error {
	app := New()
	logger := app.logger.logger
	logger.Info("Trying to start server",
		slog.String("port", a.config.Addr),
	)

	http.HandleFunc("/api/v1/calculate", app.RequestHandler(app.LoggingHandler(app.CalcHandler)))
	err := http.ListenAndServe(":"+a.config.Addr, nil)
	if err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		return err
	}

	return nil
}
