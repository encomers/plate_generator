package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"encomers/license/internal/domain/services/generator"
)

type Endpoints struct {
	plateGen generator.IPlateGenerator
	logger   *zap.Logger
}

func New(plateGen generator.IPlateGenerator, logger *zap.Logger) *Endpoints {
	return &Endpoints{
		plateGen: plateGen,
		logger:   logger,
	}
}

// NextPlate godoc
// @Summary      Получить следующий номер
// @Description  Возвращает следующий сгенерированный номерной знак
// @Tags         plates
// @Produce      plain
// @Success      200  {string}  string  "У 123 ТЕ 116 RUS"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /next [get]
func (e *Endpoints) NextPlate(c *gin.Context) {
	plate, err := e.plateGen.GetNext()
	if err != nil {
		e.logger.Error("failed to get next plate", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, plate.String())
}

// RandomPlate godoc
// @Summary      Получить случайный номер
// @Description  Возвращает случайный номерной знак
// @Tags         plates
// @Produce      plain
// @Success      200  {string}  string  "X 789 УА 116 RUS"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /random [get]
func (e *Endpoints) RandomPlate(c *gin.Context) {
	plate, err := e.plateGen.GetRandom()
	if err != nil {
		e.logger.Error("failed to get random plate", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, plate.String())
}

func (e *Endpoints) RegisterEndpoints(router *gin.Engine) {
	e.logger.Info("registering endpoints")
	router.GET("/next", e.NextPlate)
	router.GET("/random", e.RandomPlate)
}
