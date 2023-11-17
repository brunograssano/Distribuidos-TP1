package main

import (
	"github.com/brunograssano/Distribuidos-TP1/common/middleware"
	queueProtocol "github.com/brunograssano/Distribuidos-TP1/common/protocol/queues"
	"github.com/brunograssano/Distribuidos-TP1/common/queuefactory"
	"github.com/brunograssano/Distribuidos-TP1/common/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	sigs := utils.CreateSignalListener()
	env, err := InitEnv()
	if err != nil {
		log.Fatalf("Main - Ex4 Avg Calculator | Error initializing env | %s", err)
	}
	config, err := GetConfig(env)
	if err != nil {
		log.Fatalf("Main - Ex4 Avg Calculator | Error initializing Config | %s", err)
	}

	var toJourneySavers []queueProtocol.ProducerProtocolInterface
	qMiddleware := middleware.NewQueueMiddleware(config.RabbitAddress)
	qFactory := queuefactory.NewDirectExchangeProducerSimpleConsQueueFactory(qMiddleware)
	inputQueue := qFactory.CreateConsumer(config.InputQueueName)

	for i := uint(0); i < config.SaversCount; i++ {
		producer := qFactory.CreateProducer(config.OutputQueueName)
		toJourneySavers = append(toJourneySavers, producer)
	}

	avgCalculator := NewAvgCalculator(toJourneySavers, inputQueue, config)
	go avgCalculator.CalculateAvgLoop()
	<-sigs
	qMiddleware.Close()
}
