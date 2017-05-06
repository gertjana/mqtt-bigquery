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

It expexts the payload of the TTN message to contain a json with a list of name/values

It will try to map these to a struct in devices.go and insert that struct into a table in bigquery

## Data model

The application will attempt to create a Dataset for the TTN Application and Tables in there for the devices if they do not already exist.
It will infer the schema based on the Struct that is used (see devices.go)

## Add a Device

At this point in time to add a device recompilation is needed
to add a device you need to (in devices.go)
 * Add your device to the the devices list
 * If necessary create a DeviceType and a corresponding Struct

Based on the struct the bigquery part will infer the schame 

