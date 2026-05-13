// Package app содержит основную структуру приложения и его инициализацию.
//
// Здесь происходит сборка всех компонентов: генераторов, репозитория,
// HTTP-сервера и эндпоинтов.
package app

import (
	"fmt"

	"go.uber.org/zap"

	"encomers/license/internal/api/http/endpoints"
	"encomers/license/internal/api/http/server"
	lettergenerator "encomers/license/internal/domain/services/generator/realizations/lettergenerator"
	numbergenerator "encomers/license/internal/domain/services/generator/realizations/numbergenerator"
	plategenerator "encomers/license/internal/domain/services/generator/realizations/plategenerator"
	regiongenerator "encomers/license/internal/domain/services/generator/realizations/regiongenerator"
	repository "encomers/license/internal/domain/services/repository/realizations"
	models "encomers/license/internal/domain/valueObjects"
)

// vocab — алфавит, используемый для генерации буквенной части номера.
// Содержит только разрешённые символы кириллицы.
var vocab = []rune("АЕТОРНУКХСВМ")

// region — фиксированный регион по умолчанию (используется в статическом генераторе).
var region, _ = models.NewRegionPart("116", "RUS")

// App — главная структура приложения.
type App struct {
	server *server.Server
	logger *zap.Logger
}

// New создаёт и настраивает новое приложение.
//
// Выполняет сборку всех зависимостей:
//   - In-memory репозиторий
//   - Генераторы букв, цифр и региона
//   - Основной генератор номеров (PlateGenerator)
//   - HTTP-эндпоинты и сервер
//
// При критических ошибках инициализации генераторов происходит panic.
func New(logger *zap.Logger) *App {
	repo := repository.New()

	// Создаём генератор буквенной части
	letterGen, err := lettergenerator.New(vocab)
	if err != nil {
		panic(fmt.Sprintf("failed to create letter generator: %v", err))
	}

	// Создаём генератор числовой части
	numberGen := numbergenerator.New(nil)

	// Создаём статический генератор региона
	regionGen := regiongenerator.New(region)

	// Создаём основной генератор номеров
	plateGen, err := plategenerator.New(numberGen, letterGen, regionGen, repo)
	if err != nil {
		panic(fmt.Sprintf("failed to create plate generator: %v", err))
	}

	// Инициализируем HTTP-эндпоинты
	endpoints := endpoints.New(plateGen, logger)

	// Создаём HTTP-сервер
	srv := server.New(logger)
	endpoints.RegisterEndpoints(srv.Router())

	return &App{
		server: srv,
		logger: logger,
	}
}

// Run запускает HTTP-сервер приложения на указанном адресе.
func (a *App) Run(addr string) error {
	a.logger.Info("starting application", zap.String("address", addr))
	return a.server.Run(addr)
}
