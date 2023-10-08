package main

import (
	dataStructure "github.com/brunograssano/Distribuidos-TP1/common/data_structures"
	"github.com/brunograssano/Distribuidos-TP1/common/protocol"
	log "github.com/sirupsen/logrus"
)

type AvgCalculator struct {
	toInternalSaversChannels []protocol.ProducerProtocolInterface
	pricesConsumer           protocol.ConsumerProtocolInterface
}

func NewAvgCalculator(toInternalSaversChannels []protocol.ProducerProtocolInterface, pricesConsumer protocol.ConsumerProtocolInterface) *AvgCalculator {
	return &AvgCalculator{toInternalSaversChannels: toInternalSaversChannels, pricesConsumer: pricesConsumer}
}

// CalculateAvgLoop Waits for the final results from the journey savers,
// performs the calculation, and sends the results back
func (a *AvgCalculator) CalculateAvgLoop() {
	log.Infof("AvgCalculator | Started Avg Calculator loop")
	sumOfPrices := float32(0)
	sumOfRows := 0
	for sentResults := 0; sentResults < len(a.toInternalSaversChannels); sentResults++ {
		msg, ok := a.pricesConsumer.Pop()
		if !ok {
			log.Errorf("AvgCalculator | Consumer closed when not expected, exiting average calculator")
			return
		}
		log.Debugf("AvgCalculator | Received message %v", msg)
		if msg.TypeMessage != dataStructure.EOFFlightRows {
			log.Errorf("AvgCalculator | Received a message of type '%v' that was not expected, skipping...", msg.TypeMessage)
			sentResults--
			continue
		}

		prices, err := msg.DynMaps[0].GetAsFloat("localPrice")
		if err != nil {
			log.Errorf("AvgCalculator | Error getting localPrice | %v", err)
			continue
		}
		sumOfPrices += prices

		rows, err := msg.DynMaps[0].GetAsInt("localQuantity")
		if err != nil {
			log.Errorf("AvgCalculator | Error getting localQuantity | %v", err)
			continue
		}
		sumOfRows += rows
	}
	log.Infof("AvgCalculator | Received all local JourneySaver values, calculating average")
	avg := a.calculateAvg(sumOfRows, sumOfPrices)
	a.sendToJourneySavers(avg)

}

// sendToJourneySavers Sends the average to the journey savers
func (a *AvgCalculator) sendToJourneySavers(avg float32) {
	avgBytes := dataStructure.NewSerializer().SerializeFloat(avg)
	dynMap := make(map[string][]byte)
	dynMap["finalAvg"] = avgBytes
	data := []*dataStructure.DynamicMap{dataStructure.NewDynamicMap(dynMap)}
	msg := &dataStructure.Message{TypeMessage: dataStructure.FinalAvg, DynMaps: data}
	for i, channel := range a.toInternalSaversChannels {
		log.Infof("AvgCalculator | Sending average to saver %v", i)
		err := channel.Send(msg)
		if err != nil {
			log.Errorf("AvgCalculator | Error sending avg: %v", err)
		}
	}
}

// calculateAvg Performs the calculation of the average
// If the total rows is zero, returns zero
func (a *AvgCalculator) calculateAvg(sumOfRows int, sumOfPrices float32) float32 {
	if sumOfRows == 0 {
		log.Warnf("AvgCalculator | Total rows is zero")
		return float32(0)
	}
	avg := sumOfPrices / float32(sumOfRows)
	log.Infof("AvgCalculator | Sum of prices: %v | Total rows: %v | Avg: %v", sumOfPrices, sumOfRows, avg)
	return avg
}
