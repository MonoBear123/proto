// Пакет grpcPredict предоставляет клиент для взаимодействия с gRPC-сервисом предсказаний.
package grpcPredict

import (
	"context"
	"github.com/MonoBear123/proto/protos/gen/go/predict"
	"google.golang.org/grpc"
	"time"
)

type PredictClient struct {
	client predict.StonksPredictorClient
}

// New создает новый клиент PredictClient для подключения к gRPC-серверу предсказаний.
//
// Параметры:
//   - address: строка, содержащая адрес сервера (например, "localhost:50052").
//
// Возвращает:
//   - *PredictClient: клиент для взаимодействия с gRPC-сервисом.
//
// В случае невозможности установить соединение вызывается panic().
func New(address string) *PredictClient {
	// Создание контекста с таймаутом для попытки подключения.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("Failed to connect to Predictor" + err.Error())
	}
	return &PredictClient{predict.NewStonksPredictorClient(client)}
}

// GetPrediction отправляет запрос на предсказание и возвращает результат.
//
// Возвращает:
//   - []float32: список предсказанных чисел.
//   - error: ошибка в случае неудачного запроса.
func (p *PredictClient) GetPrediction() ([]float32, error) {
	res, err := p.client.Predictor(context.Background(), &predict.PredictorRequest{Query: ""})
	if err != nil {
		return nil, err
	}
	return res.Numbers, nil
}
