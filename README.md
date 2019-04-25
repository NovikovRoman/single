[half fork from marcsauter/single](https://github.com/marcsauter/single)


```go
package main

import (
	"github.com/NovikovRoman/single"
	"log"
	"time"
)

func main() {
	s := single.New("your-app-name")
	busy, err := s.CheckLock()
	if err != nil {
		log.Fatalf("failed to acquire exclusive app lock: %v", err)
	}
	defer func(s *Single) {
        if err := s.TryUnlock(); err != nil {
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