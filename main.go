package main

import (
	BQ "cloud.google.com/go/bigquery"
  MQTT "github.com/eclipse/org.eclipse.paho.mqtt.golang"
	"golang.org/x/net/context"
  "google.golang.org/api/iterator"
  "encoding/json"
  "fmt"
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

//TODO Cache datasets
func createDataSetIfNotExists(bqClient *BQ.Client, datasetID string) error {
  ctx := context.Background()

  it := bqClient.Datasets(ctx)
  for {
    dataset, err := it.Next()
    if err == iterator.Done {
      fmt.Println("No dataset found")
      break
    }
    if err != nil {
      return err
    }
    if dataset.DatasetID == datasetID {
      return nil
    }
  }
  fmt.Println("Creating Dataset ", datasetID)
  return bqClient.Dataset(datasetID).Create(ctx)
}

//TODO cache tables
func createTableIfNotExists(bqClient *BQ.Client, datasetID string, tableID string, item interface{}) error {
  ctx := context.Background()

  ts := bqClient.Dataset(datasetID).Tables(ctx)
  for {
    t, err := ts.Next()
    if err == iterator.Done {
      fmt.Println("No Table found")
      break
    }
    if err != nil {
      return err
    }
    if t.TableID == tableID {
      return nil
    }
  }

  schema, err := BQ.InferSchema(item)
  if err != nil {
    return err
  }

  fmt.Println("Creating Table ", tableID)
  return bqClient.Dataset(datasetID).Table(tableID).Create(ctx, schema)
}

func insertRow(datasetID, tableID string, item interface{}) error {
	ctx := context.Background()
	bqClient, err := BQ.NewClient(ctx, project) 
  if err != nil {
		return err
	}
  if err := createDataSetIfNotExists(bqClient, datasetID);err != nil {
    return err
  }
  if err := createTableIfNotExists(bqClient, datasetID, tableID, item);err != nil {
    return err
  }

	u := bqClient.Dataset(datasetID).Table(tableID).Uploader()
	if err := u.Put(ctx, item); err != nil {
		return err
	}
	return nil
}
