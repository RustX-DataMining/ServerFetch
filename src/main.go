package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	Status     bool     `json:"status"`
	State      string   `json:"state"`
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Online     int      `json:"online"`
	MaxPlayers int      `json:"maxPlayers"`
	Map        string   `json:"map"`
	Mod        string   `json:"mod"`
	Tags       []string `json:"tags"`
	Address    string   `json:"address"`
}

type Config struct {
	Url      string
	Interval time.Duration
}

func loadConfig() Config {
	seconds, err := strconv.Atoi(os.Getenv("SERVER_FETCH_INTERVAL"))
	if err != nil || seconds <= 0 {
		seconds = 30
	}

	return Config{
		Url:      os.Getenv("SERVER_FETCH_API_URL"),
		Interval: time.Duration(seconds) * time.Second,
	}
}

func fetchServers(cfg *Config) ([]Server, error) {
	resp, err := http.Get(cfg.Url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var servers []Server
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func loop(cfg *Config) {
	for {
		startTime := time.Now()
		servers, err := fetchServers(cfg)
		latency := time.Since(startTime)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Servers (%v):\n", latency)
			fmt.Println(servers)
		}

		time.Sleep(cfg.Interval)
	}
}

func main() {
	cfg := loadConfig()
	go loop(&cfg)
	select {}
}
