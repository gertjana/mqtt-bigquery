package main

import (
	BQ "cloud.google.com/go/bigquery"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/org.eclipse.paho.mqtt.golang"
	"golang.org/x/net/context"
	"os"
)

var proj string = "iot-lora-165218"

//define a function for the default message handler
var publishHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var subscribeHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	item := decode(msg.Payload())
	fmt.Printf("MSG: %s\n", item)

	insertRow("IoT", "pressure", item)
}

func main() {
	done := make(chan bool)

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://eu.thethings.network:1883")
	opts.SetUsername("test_gertjan")
	opts.SetPassword("ttn-account-v2.ZhjCVxlzIP1I_7izCc7qn8ZWVWsBYx1DzYEN3RA3_Xw")
	opts.SetClientID("mqtt-bigquery-bridge")
	opts.SetDefaultPublishHandler(publishHandler)

	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("MQTT Client connected")
	}

	if token := mqttClient.Subscribe("+/devices/+/up", 1, subscribeHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	} else {
		fmt.Println("Subscribed")
	}

	<-done

	//unsubscribe from /go-mqtt/sample
	if token := mqttClient.Unsubscribe("+/devices/+/up"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	} else {
		fmt.Println("Unsubcribed")
	}

	mqttClient.Disconnect(250)
}

func decode(bytes []byte) Item {

	var f interface{}
	err := json.Unmarshal(bytes, &f)

	if err != nil {
		fmt.Println(err)
	}

	m := f.(map[string]interface{})

	fields := m["payload_fields"].(map[string]interface{})
	metadata := m["metadata"].(map[string]interface{})

	celcius := fields["celcius"].(float64)
	mbar := fields["mbar"].(float64)
	timestamp := metadata["time"].(string)

	return Item{celcius, mbar, timestamp}
}

type Item struct {
	Celcius   float64
	Mbar      float64
	Timestamp string
}

func insertRow(datasetID, tableID string, item Item) error {
	ctx := context.Background()
	bqClient, err := BQ.NewClient(ctx, proj)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("BigQuery Client initialised")
	}

	u := bqClient.Dataset(datasetID).Table(tableID).Uploader()

	if err := u.Put(ctx, item); err != nil {
		return err
	}

	return nil
}
