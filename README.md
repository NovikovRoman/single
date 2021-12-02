[half fork from marcsauter/single](https://github.com/marcsauter/single)

```go
package main

import (
	"github.com/NovikovRoman/single"
	"log"
	"time"
)

func main() {
	var (
		busy bool
		err  error
	)

	s := single.New("your-app-name")
	busy, err = s.Lock()
	if err != nil {
		log.Fatalf("failed to acquire exclusive app lock: %v", err)
	}
	defer func(s *single.Single) {
		if err = s.Unlock(); err != nil {
			log.Fatal(err)
		}
	}(s)

	if busy {
		log.Fatal("another instance of the app is already running, exiting")
	}

	log.Println("working")
	time.Sleep(60 * time.Second)
	log.Println("finished")
}
```