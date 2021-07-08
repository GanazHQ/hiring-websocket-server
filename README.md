# Customer Service Phone Event Simulation _(Over Websockets!)_

This server simulates a third-party system which manages inbound customer service phone calls.

Periodically, this server sends a message to any connected client. The majority of messages are _JSON events_ representing a pending customer call in the following format:

```typescript
type CustomerCall = {
  first_name: string;
  last_name: string;
  sip: string;
  city: string;
  state: string;
  phone_number: string;
  
  // 1 - 100, with 100 being the most important
  priority: number;
  
  // ISO 8601 format
  timestamp: string;
}
```

Occasionally, there are other events sent--those should be ignored.

## Running the server

The easiest way is with Docker by getting a pre-built image from [Docker Hub](https://hub.docker.com/r/kyleemail/websocket-server):

```shell
docker pull kyleemail/websocket-server
docker run --rm -p 7777:7777 kyleemail/websocket-server
```

The default port of 7777 can be changed with an environment variable:
```shell
docker run --rm -p 8811:8811 -e PORT=8811 kyleemail/websocket-server
```

## Testing the server

Download [example.htm](./example.htm) and open it in your local browser. It connects to the server on localhost:7777 and displays any messages:

![example browser run](example%20run.gif)
