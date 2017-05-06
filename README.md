# mqtt-bigquery
Stores mqtt messages into bigquery

This is a PoC to store messages from Lora Sensors on the TTN network into bigquery

It requires the following ENV variables 

| ENV Variable  | description   |
| ------------- | ------------- |
| GOOGLE_APPLICATION_CREDENTIALS | Path to saved key after creating service account in console.cloud.google.com  |
| TTN_APP | Name of the application in console.thethingsnetwork.org  |
| TTN_KEY | Access key of TTN application |
| GCP_PROJECT | id of the Google Cloud Project |
| MQTT_CLIENT_ID | client id that subscribes to the MQTT topic |

It expexts the payload of the TTN message to contain a json with a list of name/values

It will try to map these to a struct in devices.go and insert that struct into a table in bigquery

## Docker

the makefile will create a Dockerfile (change the organisation)

NOTE: the GOOGLE_APPLICATION_CREDENTIALS is set to /credentials in the Dockerfile

copy the service account json to a location on the docker host and add that as a volume
```
docker run -d \
  -e TTN_APP="app name" \ 
  -e TTN_KEY="app key" \
  -e GCP_PROJECT="google cloud project" \
  -e MQTT_CLIENT_ID="my-mqtt-client" \
  -v /local/path/to/credentials:/credentials \
  --name mqtt-bq \
  organisation/mqtt-bq
```


## Data model

The application will attempt to create a Dataset for the TTN Application and Tables in there for the devices if they do not already exist.
It will infer the schema based on the Struct that is used (see devices.go)

## Add a Device

At this point in time to add a device recompilation is needed
to add a device you need to (in devices.go)
 * Add your device to the the devices list
 * If necessary create a DeviceType and a corresponding Struct

Based on the struct the bigquery part will infer the schema

