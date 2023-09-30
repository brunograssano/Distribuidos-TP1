package main

import (
	dataStructures "github.com/brunograssano/Distribuidos-TP1/common/data_structures"
	"github.com/brunograssano/Distribuidos-TP1/common/middleware"
	"testing"
	"time"
)

type (
	mockConsumer struct {
		inputChannel chan []byte
		ok           bool
	}
)

func (m *mockConsumer) Pop() ([]byte, bool) {
	if !m.ok {
		return []byte{}, m.ok
	}
	msg, ok := <-m.inputChannel
	return msg, ok
}

type (
	mockProducer struct {
		outputChannel chan []byte
	}
)

func (m *mockProducer) Send(data []byte) {
	m.outputChannel <- data
}

func TestShouldGetAMessageProcessItAndSendItToAllChannels(t *testing.T) {
	pConfig := &ProcessorConfig{}
	input := make(chan []byte, 10)
	outputEx13 := make(chan []byte, 10)
	outputEx2 := make(chan []byte, 10)
	outputEx4 := make(chan []byte, 10)
	serializer := dataStructures.NewDynamicMapSerializer()

	mConsumer := &mockConsumer{
		inputChannel: input,
		ok:           true,
	}
	mProducer13 := &mockProducer{
		outputChannel: outputEx13,
	}
	mProducer2 := &mockProducer{
		outputChannel: outputEx2,
	}
	mProducer4 := &mockProducer{
		outputChannel: outputEx4,
	}

	processor := &DataProcessor{
		processorId:    0,
		c:              pConfig,
		consumer:       mConsumer,
		producersEx123: []middleware.ProducerInterface{mProducer2, mProducer13},
		producersEx4:   mProducer4,
		serializer:     serializer,
		ex123Columns:   []string{"startingAirport", "segmentsArrivalAirportCode", "totalStopovers", "route"},
		ex4Columns:     []string{"route"},
	}

	dynMap := make(map[string][]byte)
	dynMap["startingAirport"] = []byte("FRA")
	dynMap["segmentsArrivalAirportCode"] = []byte("EZE")
	dynMap["col"] = []byte("Even more data")

	row := dataStructures.NewDynamicMap(dynMap)

	go processor.ProcessData()

	input <- serializer.Serialize(row)
	close(input)

	sentResponseToAll := [3]bool{false, false, false}
	for i := 0; i < 3; i++ {
		select {
		case result := <-outputEx13:
			newRow := serializer.Deserialize(result)
			if newRow.GetColumnCount() != 4 {
				t.Errorf("RowCount expected was 2")
			}
			sentResponseToAll[0] = true
		case result := <-outputEx2:
			newRow := serializer.Deserialize(result)
			if newRow.GetColumnCount() != 4 {
				t.Errorf("RowCount expected was 2")
			}
			sentResponseToAll[1] = true
		case result := <-outputEx4:
			newRow := serializer.Deserialize(result)
			if newRow.GetColumnCount() != 1 {
				t.Errorf("RowCount expected was 2")
			}
			sentResponseToAll[2] = true
		case <-time.After(1 * time.Second):
			t.Errorf("Timeout! Should have finished by now...")

		}
	}
	for i := 0; i < 3; i++ {
		if !sentResponseToAll[i] {
			t.Errorf("Missing response from a channel")
		}
	}
}

func TestShouldProcessTheDataOfEx123(t *testing.T) {
	serializer := dataStructures.NewDynamicMapSerializer()
	processor := &DataProcessor{
		processorId:  0,
		serializer:   serializer,
		ex123Columns: []string{"totalStopovers", "route"},
	}
	dynMap := make(map[string][]byte)
	dynMap["startingAirport"] = []byte("FRA")
	dynMap["segmentsArrivalAirportCode"] = []byte("CDG||EZE")
	dynMap["col"] = []byte("Even more data")

	row := dataStructures.NewDynamicMap(dynMap)
	row, err := processor.processEx123Row(row)
	if err != nil {
		t.Errorf("Got error when processing ex123 row: %v", err)
	}

	route, err := row.GetAsString("route")
	if err != nil {
		t.Errorf("Got error when getting route: %v", err)
	}

	if route != "FRA||CDG||EZE" {
		t.Errorf("Expecting FRA||CDG||EZE route but got %v", route)
	}

	stopovers, err := row.GetAsInt("totalStopovers")
	if err != nil {
		t.Errorf("Got error when getting stopovers: %v", err)
	}

	if stopovers != 1 {
		t.Errorf("Expecting 1 stopover but got %v", stopovers)
	}
}

func TestShouldReturnAnErrorIfTheSegmentsColDoesNotExist(t *testing.T) {
	serializer := dataStructures.NewDynamicMapSerializer()
	processor := &DataProcessor{
		processorId:  0,
		serializer:   serializer,
		ex123Columns: []string{"totalStopovers", "route"},
	}
	dynMap := make(map[string][]byte)
	dynMap["startingAirport"] = []byte("FRA")

	row := dataStructures.NewDynamicMap(dynMap)
	row, err := processor.processEx123Row(row)
	if err == nil {
		t.Errorf("Didn't got an error when processing segments column")
	}

}

func TestShouldReturnAnErrorIfTheStartingAirportColDoesNotExist(t *testing.T) {
	serializer := dataStructures.NewDynamicMapSerializer()
	processor := &DataProcessor{
		processorId:  0,
		serializer:   serializer,
		ex123Columns: []string{"totalStopovers", "route"},
	}
	dynMap := make(map[string][]byte)
	dynMap["segmentsArrivalAirportCode"] = []byte("CDG||EZE")

	row := dataStructures.NewDynamicMap(dynMap)
	row, err := processor.processEx123Row(row)
	if err == nil {
		t.Errorf("Didn't got an error when processing starting column")
	}

}