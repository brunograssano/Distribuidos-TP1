package main

import (
	"data_processor/processor"
	"github.com/brunograssano/Distribuidos-TP1/common/middleware"
	"github.com/brunograssano/Distribuidos-TP1/common/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	sigs := utils.CreateSignalListener()

	env, err := processor.InitEnv()
	if err != nil {
		log.Fatalf("Main - DataProcessor | Error initializing env | %s", err)
	}

	processorConfig, err := processor.GetConfig(env)
	if err != nil {
		log.Fatalf("Main - DataProcessor | Error initializing config | %s", err)
	}

	qMiddleware := middleware.NewQueueMiddleware(processorConfig.RabbitAddress)
	for i := 0; i < processorConfig.GoroutinesCount; i++ {
		r := processor.NewDataProcessor(i, qMiddleware, processorConfig)
		go r.ProcessData()
	}
	<-sigs
	qMiddleware.Close()
}
