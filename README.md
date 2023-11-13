# MUECKE - MQTT Bruecke

MUECKE is a MQTT bridge, or "Bruecke" (German for bridge), designed to connect to any MQTT broker and forward predefined topics to a local
MQTT broker.

## Configuration

MUECKE uses a YAML configuration file for its settings. Here's a sample configuration:

```yaml
remote_broker:
  broker: tcp://127.0.0.1:1883
  client_id: my-client
  username: my-username
  password: my-password
  bridge_topic: remote-bridge

app_configs:
  - app_name: App1
  - app_name: App2
  - app_name: App3
```

| Field Name   | Description                                                                                                                                      |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| broker       | The address of the MQTT broker to connect to. For example, `tcp://127.0.0.1:1883` means the broker is running on the local machine on port 1883. |
| client_id    | The unique identifier for this client.                                                                                                           |
| username     | The username for authenticating with the MQTT broker.                                                                                            |
| password     | The password for authenticating with the MQTT broker.                                                                                            |
| bridge_topic | The topic that this bridge will use for communicating with the MQTT broker.                                                                      |
| app_configs  | A list of applications that will be connected to this broker. Each application is identified by its name.                                        |

## Running MUECKE

To build the executable using Go, follow the steps below:

1. Navigate to the directory containing your Go source code.
2. Run the following command to compile the source code into an executable file:

   ```bash
   go build

   ```

## Systemv Service

Here's an example of a System V service file for MUECKE (`muecke.service`).
Replace the path to the executable with the actual path to your executable. Then
copy and paste the contents of the file into `/etc/systemd/system/muecke.service`.
This will start MUECKE automatically on boot.

```ini
[Unit]
Description=MUECKE Service
After=network.target

[Service]
ExecStart=/path/to/muecke
Restart=always

[Install]
WantedBy=multi-user.target
```

## External Libraries

We used the following libraries for this project:

- [Mochi-MQTT](https://github.com/mochi-mqtt/server): A MQTT client library designed using the Mochi framework, for building network applications with Python 3.6+ type hints.
- [Paho](https://github.com/eclipse/paho.mqtt.python): The Paho Python MQTT Client is a client library designed to provide some helper functions and objects for developers who are writing applications aimed at integrating devices into an MQTT environment.
