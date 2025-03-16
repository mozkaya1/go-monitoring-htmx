package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/mozkaya1/go-monitoring-htmx/internal/api"
	"github.com/mozkaya1/go-monitoring-htmx/internal/hardware"
)

type server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscribersMu           sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type subscriber struct {
	msgs chan []byte
}

func NewServer() *server {
	s := &server{
		subscriberMessageBuffer: 20,
		subscribers:             make(map[*subscriber]struct{}),
	}
	s.mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *server) addSubscriber(subscriber *subscriber) {
	s.subscribersMu.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMu.Unlock()
	fmt.Println("Added subscriber", subscriber)
}

func (s *server) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var c *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)
	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *server) publishMsg(msg []byte) {
	cs.subscribersMu.Lock()
	defer cs.subscribersMu.Unlock()

	for s := range cs.subscribers {
		s.msgs <- msg
	}
}

func main() {
	fmt.Println("Starting monitor server on port 8000")
	fmt.Println("open browser at http://localhost:8000")
	s := NewServer()

	go func(srv *server) {
		for {
			// musti := time.Now().Format("2006-01-02 15:04:05")
			systemData, err := hardware.GetSystemSection()
			if err != nil {
				fmt.Println(err, "Get System Error")
			}
			diskData, err := hardware.GetDiskSection()
			if err != nil {
				fmt.Println(err, "Disk Error")
			}
			cpuData, err := hardware.GetCpuSection()
			if err != nil {
				fmt.Println(err, "CPU Error")
			}
			load, err := hardware.GetLoad()
			if err != nil {
				fmt.Println(err, "LOAD Error")
			}

			// My system temp has always 1 warning -- Disabled Error spamming
			systemp, _ := hardware.GetSensors()
			// if err != nil {
			// 	fmt.Println(err, "Temp Error")
			// }

			// If current user has not access the docker system you will get error. Run script with sudo to get docker monitoring if you have docker
			dock, err := hardware.GetDocker()
			if err != nil {
				fmt.Println(err, "Docker ERROR")

			}

			coin, err := api.GetApi()
			if err != nil {
				log.Println(err, "API Error")
			}

			timeStamp := time.Now().Format("2006-01-02 15:04:05")
			msg := []byte(`
      <div hx-swap-oob="innerHTML:#update-timestamp">
        <p><i style="color: green" class="fa fa-circle"></i> ` + timeStamp + `</p>
      </div>
      <div hx-swap-oob="innerHTML:#temp"> ` + coin.WeatherBucket.Temp + `</div>
      <div hx-swap-oob="innerHTML:#weather"> ` + coin.WeatherBucket.WeatherDesc + `</div>
      <div hx-swap-oob="innerHTML:#location"> ` + coin.WeatherBucket.Location + `</div>
      <div hx-swap-oob="innerHTML:#btc"> ` + coin.Crypto.Asset["BTCUSDT"].LastPrice + `</div>
      <div hx-swap-oob="innerHTML:#eth"> ` + coin.Crypto.Asset["ETHUSDT"].LastPrice + `</div>
      <div hx-swap-oob="innerHTML:#system-data">` + systemData + `</div>
      <div hx-swap-oob="innerHTML:#cpu-data">` + cpuData + `</div>
      <div hx-swap-oob="innerHTML:#load">` + load + `</div>
      <div hx-swap-oob="innerHTML:#systemp">` + systemp + `</div>
      <div hx-swap-oob="innerHTML:#dock">` + dock + `</div>
      <div hx-swap-oob="innerHTML:#disk-data">` + diskData + `</div>`)
			srv.publishMsg(msg)

			// Refreshing DATA inverval
			time.Sleep(5 * time.Second)
		}
	}(s)

	err := http.ListenAndServe(":8000", &s.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
