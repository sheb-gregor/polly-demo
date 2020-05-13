# Polly Demo

Polly is a simple pub/sub server and client for it.

## Build

- Pull Code

```shell script
git clone https://github.com/sheb-gregor/polly-demo.git

# OR

git clone git@github.com:sheb-gregor/polly-demo.git

cd polly-demo
```

- Test

```shell script
go test --race ./...
```

- Build

```shell script
go build .
```

## Usage

Usage of `./polly-demo`:

```text
./polly-demo -config  ./config.yaml

  -config string
        path to file with service configuration (default "config.yaml")
```

Example of `config.yaml`

```yaml
api:
  host: 127.0.0.1
  port: 3000
```

## API 

API Description provided in [polly.http](./polly.http)


## Client

**PollyClient** from the package `github.com/sheb-gregor/polly-demo/client` can be used to interact with the Polly Broker server from the Go application.
  
```go
package main

import (
	"encoding/json"
	"log"

	"github.com/sheb-gregor/polly-demo/client"
)

func main() {
	topic := "test_topic"
	name := "bobby"
	message := json.RawMessage(`{"my_key":"my_message"}`)

	pollyClient, err := client.NewClient("http://127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	err = pollyClient.Subscribe(topic, name)
	if err != nil {
		log.Fatal(err)
	}

	err = pollyClient.Publish(topic, message)
	if err != nil {
		log.Fatal(err)
	}
	msg, err := pollyClient.Poll(topic, name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New message: ", string(msg))

	err = pollyClient.Unsubscribe(topic, name)
	if err != nil {
		log.Fatal(err)
	}
}
```


## Q&A

Finally, provide the answer to the following questions: 

##### What is the message publish algorithm complexity in big-O notation? 

**O(N)** - where N is number of subscribers; because we add message ID into FIFO for each subscriber.

 
##### What is the message poll algorithm complexity in big-O notation?

**O(1)** - message retrieval is almost constant, because we take it by an identifier from the map.
 
##### What is the memory (space) complexity in big-O notation for the algorithm? 

**O(N)** -  where N is number of messages, because we store only one copy of message.

