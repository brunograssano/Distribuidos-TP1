package saver

import (
	dataStructures "github.com/brunograssano/Distribuidos-TP1/common/data_structures"
	"github.com/brunograssano/Distribuidos-TP1/common/filemanager"
	"github.com/brunograssano/Distribuidos-TP1/common/middleware"
	queueProtocol "github.com/brunograssano/Distribuidos-TP1/common/protocol/queues"
	"github.com/brunograssano/Distribuidos-TP1/common/serializer"
	"github.com/brunograssano/Distribuidos-TP1/common/utils"
	log "github.com/sirupsen/logrus"
)

// SimpleSaver Structure that handles the final results
type SimpleSaver struct {
	c        *Config
	consumer queueProtocol.ConsumerProtocolInterface
	canSend  chan string
}

// NewSimpleSaver Creates a new saver for the results
func NewSimpleSaver(qMiddleware *middleware.QueueMiddleware, c *Config, canSend chan string) *SimpleSaver {
	consumer := queueProtocol.NewConsumerQueueProtocolHandler(qMiddleware.CreateConsumer(c.InputQueueName, true))
	return &SimpleSaver{c: c, consumer: consumer, canSend: canSend}
}

// SaveData Saves the results from the queue in a file
func (s *SimpleSaver) SaveData() {
	log.Infof("SimpleSaver | Goroutine started")
	for {
		msgStruct, ok := s.consumer.Pop()
		if !ok {
			log.Infof("SimpleSaver | Exiting saver")
			return
		}
		if msgStruct.TypeMessage == dataStructures.EOFFlightRows {
			log.Infof("SimpleSaver | Received all results. Closing saver...")
			folder, err := filemanager.MoveFiles([]string{s.c.OutputFileName})
			if err != nil {
				log.Errorf("SimpleSaver | Error moving to file to folder | %v", err)
				return
			}
			s.canSend <- folder
		} else if msgStruct.TypeMessage == dataStructures.FlightRows {
			err := s.handleFlightRows(msgStruct)
			if err != nil {
				log.Errorf("SimpleSaver | Error handling flight rows. Closing saver...")
				return
			}
		}
	}
}

func (s *SimpleSaver) handleFlightRows(msgStruct *dataStructures.Message) error {
	writer, err := filemanager.NewFileWriter(s.c.OutputFileName)
	if err != nil {
		log.Errorf("SimpleSaver | Error opening file writer of output")
		return err
	}
	defer utils.CloseFileAndNotifyError(writer.FileManager)
	return s.writeRowsToFile(msgStruct.DynMaps, writer)
}

func (s *SimpleSaver) writeRowsToFile(rows []*dataStructures.DynamicMap, writer filemanager.OutputManagerInterface) error {
	for _, row := range rows {
		line := serializer.SerializeToString(row)
		log.Debugf("SimpleSaver | Saving line: %v", line)
		err := writer.WriteLine(line)
		if err != nil {
			log.Errorf("SimpleSaver | Error writing to file | %v", err)
			return err
		}
	}
	return nil
}