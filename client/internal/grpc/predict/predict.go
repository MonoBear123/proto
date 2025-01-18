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

func New(address string) *PredictClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("Failed to connect to Predictor" + err.Error())
	}
	return &PredictClient{predict.NewStonksPredictorClient(client)}
}

func (p *PredictClient) GetPrediction(secid string) ([]float32, error) {
	res, err := p.client.Predictor(context.Background(), &predict.PredictorRequest{Query: secid})
	if err != nil {
		return nil, err
	}
	return res.Numbers, nil
}
