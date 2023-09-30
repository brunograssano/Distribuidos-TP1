package main

import (
	"encoding/binary"
	"filters_config"
	"github.com/brunograssano/Distribuidos-TP1/common/data_structures"
	"github.com/brunograssano/Distribuidos-TP1/common/filters"
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

func TestGettingARowWithTotalStopoversLessThanThreeShouldNotSendIt(t *testing.T) {
	input := make(chan []byte)
	output := make(chan []byte)
	serializer := data_structures.NewDynamicMapSerializer()

	mockCons := &mockConsumer{
		inputChannel: input,
		ok:           true,
	}
	arrayProducers := make([]middleware.ProducerInterface, 1)
	arrayProducers[0] = &mockProducer{
		outputChannel: output,
	}
	filterEscalas := &FilterEscalas{
		filterId:   0,
		config:     &filters_config.FilterConfig{},
		consumer:   mockCons,
		producers:  arrayProducers,
		serializer: data_structures.NewDynamicMapSerializer(),
		filter:     filters.NewFilter(),
	}
	go filterEscalas.FilterEscalas()

	dynMap := make(map[string][]byte)
	dynMap["totalStopovers"] = make([]byte, 4)
	binary.BigEndian.PutUint32(dynMap["totalStopovers"], uint32(2))
	row := data_structures.NewDynamicMap(dynMap)
	input <- serializer.Serialize(row)
	close(input)
	select {
	case <-output:
		t.Errorf("RowCount expected was 0")

	case <-time.After(2 * time.Second):
	}
}

func TestGettingARowWithTotalStopoversEqualToThreeShouldSendIt(t *testing.T) {
	input := make(chan []byte)
	output := make(chan []byte)
	serializer := data_structures.NewDynamicMapSerializer()

	mockCons := &mockConsumer{
		inputChannel: input,
		ok:           true,
	}
	arrayProducers := make([]middleware.ProducerInterface, 1)
	arrayProducers[0] = &mockProducer{
		outputChannel: output,
	}
	filterEscalas := &FilterEscalas{
		filterId:   0,
		config:     &filters_config.FilterConfig{},
		consumer:   mockCons,
		producers:  arrayProducers,
		serializer: data_structures.NewDynamicMapSerializer(),
		filter:     filters.NewFilter(),
	}
	go filterEscalas.FilterEscalas()

	dynMap := make(map[string][]byte)
	dynMap["totalStopovers"] = make([]byte, 4)
	binary.BigEndian.PutUint32(dynMap["totalStopovers"], uint32(3))
	row := data_structures.NewDynamicMap(dynMap)
	input <- serializer.Serialize(row)
	close(input)
	select {
	case result := <-output:
		newRow := serializer.Deserialize(result)
		ts, err := newRow.GetAsInt("totalStopovers")
		if err != nil {
			t.Errorf("Error getting totalStopovers...")
		}
		if ts < 3 {
			t.Errorf("Received a row that was not expected, has less than 3 stopovers...")
		}

	case <-time.After(1 * time.Second):
		t.Errorf("Timeout! Should have finished by now...")
	}
}

func TestGettingARowWithTotalStopoversGreaterThanThreeShouldSendIt(t *testing.T) {
	input := make(chan []byte)
	output := make(chan []byte)
	serializer := data_structures.NewDynamicMapSerializer()

	mockCons := &mockConsumer{
		inputChannel: input,
		ok:           true,
	}
	arrayProducers := make([]middleware.ProducerInterface, 1)
	arrayProducers[0] = &mockProducer{
		outputChannel: output,
	}
	filterEscalas := &FilterEscalas{
		filterId:   0,
		config:     &filters_config.FilterConfig{},
		consumer:   mockCons,
		producers:  arrayProducers,
		serializer: data_structures.NewDynamicMapSerializer(),
		filter:     filters.NewFilter(),
	}
	go filterEscalas.FilterEscalas()

	dynMap := make(map[string][]byte)
	dynMap["totalStopovers"] = make([]byte, 4)
	binary.BigEndian.PutUint32(dynMap["totalStopovers"], uint32(4))
	row := data_structures.NewDynamicMap(dynMap)
	input <- serializer.Serialize(row)
	close(input)
	select {
	case result := <-output:
		newRow := serializer.Deserialize(result)
		ts, err := newRow.GetAsInt("totalStopovers")
		if err != nil {
			t.Errorf("Error getting totalStopovers...")
		}
		if ts < 3 {
			t.Errorf("Received a row that was not expected, has less than 3 stopovers...")
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout! Should have finished by now...")
	}
}

func TestWithLessEqualAndGreaterCasesTogetherShouldSendTwoOutOfThree(t *testing.T) {
	input := make(chan []byte, 3)
	output := make(chan []byte, 3)
	serializer := data_structures.NewDynamicMapSerializer()

	mockCons := &mockConsumer{
		inputChannel: input,
		ok:           true,
	}
	arrayProducers := make([]middleware.ProducerInterface, 1)
	arrayProducers[0] = &mockProducer{
		outputChannel: output,
	}
	filterEscalas := &FilterEscalas{
		filterId:   0,
		config:     &filters_config.FilterConfig{},
		consumer:   mockCons,
		producers:  arrayProducers,
		serializer: data_structures.NewDynamicMapSerializer(),
		filter:     filters.NewFilter(),
	}
	go filterEscalas.FilterEscalas()

	for i := 0; i < 3; i++ {
		dynMap := make(map[string][]byte)
		dynMap["totalStopovers"] = make([]byte, 4)
		binary.BigEndian.PutUint32(dynMap["totalStopovers"], uint32(2+i))
		row := data_structures.NewDynamicMap(dynMap)
		input <- serializer.Serialize(row)
	}
	close(input)
	rowCountRecvd := 0
	for i := 0; i < 3; i++ {
		select {
		case result := <-output:
			newRow := serializer.Deserialize(result)
			rowCountRecvd++
			ts, err := newRow.GetAsInt("totalStopovers")
			if err != nil {
				t.Errorf("Error getting totalStopovers...")
			}
			if ts < 3 {
				t.Errorf("Received a row that was not expected, has less than 3 stopovers...")
			}

		case <-time.After(1 * time.Second):
			// Should only read two messages from channel (3 and 4 stopovers)
			if rowCountRecvd != 2 {
				t.Errorf("Timeout! Should have finished by now...")
			}
		}
	}
	if rowCountRecvd != 2 {
		t.Errorf("Expected to receive only 2 rows, but %v were received", rowCountRecvd)
	}
}