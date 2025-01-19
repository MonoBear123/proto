package predictHandler

import (
	"client/internal/grpc/predict"
	"client/internal/lib/parseJWT"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PredictHandler struct {
	client *grpcPredict.PredictClient
}

func New(client *grpcPredict.PredictClient) *PredictHandler {
	return &PredictHandler{
		client: client,
	}
}

func (p *PredictHandler) Predict(ctx echo.Context) error {
	cookie, err := ctx.Cookie("token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if err := parseJWT.ParseToken(cookie.Value); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	secid := ctx.FormValue("secid")
	res, err := p.client.GetPrediction(secid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(res) == 1 {
		return ctx.JSON(http.StatusOK, map[string]string{"message:": "Модель уже  находится на обучении, пожалуйста подождите "})
	}
	if len(res) == 2 {
		return ctx.JSON(http.StatusOK, map[string]string{"message:": "Достигнуто максимальное количество одновременных обучений,пожалуйста подождите"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"result": res})
}
