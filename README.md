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

