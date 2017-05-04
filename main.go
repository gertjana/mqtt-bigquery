package main

import (
	BQ "cloud.google.com/go/bigquery"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/org.eclipse.paho.mqtt.golang"
	"golang.org/x/net/context"
	"os"
)

var ttn_app string = os.Getenv("TTN_APP")
var ttn_key string = os.Getenv("TTN_KEY")
var project string = os.Getenv("GCP_PROJECT")

func main() {
	done := make(chan bool)

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://eu.thethings.network:1883")
	opts.SetUsername(ttn_app)
	opts.SetPassword(ttn_key)
	opts.SetClientID("mqtt-bigquery-bridge")

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

	<-done //block forever

	//unsubscribe from /go-mqtt/sample
	if token := mqttClient.Unsubscribe("+/devices/+/up"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	} else {
		fmt.Println("Unsubcribed")
	}

	mqttClient.Disconnect(250)
}

var subscribeHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	appid, devid, item := decode(msg.Payload())
	fmt.Printf("MSG: %s\n", item)

	err := insertRow(appid, devid, item)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

func decode(bytes []byte) (string, string, interface{}) {
	fmt.Println("1")
	var f interface{}
	err := json.Unmarshal(bytes, &f)

	if err != nil {
		fmt.Println(err)
	}

	m := f.(map[string]interface{})

	app_id := m["app_id"].(string)
	metadata := m["metadata"].(map[string]interface{})
	fields := m["payload_fields"].(map[string]interface{})

	dev_id := m["dev_id"].(string)
	timestamp := metadata["time"].(string)

	item := deviceFromFields(dev_id, timestamp, fields)
	return app_id, dev_id, item
}

func insertRow(datasetID, tableID string, item interface{}) error {
	ctx := context.Background()
	bqClient, err := BQ.NewClient(ctx, project)
	if err != nil {
		panic(err)
	}

	u := bqClient.Dataset(datasetID).Table(tableID).Uploader()

	if err := u.Put(ctx, item); err != nil {
		return err
	}

	return nil
}
