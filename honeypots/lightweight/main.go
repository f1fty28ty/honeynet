package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp   string `json:"@timestamp"`
	Honeypot    string `json:"honeypot"`
	EventType   string `json:"event_type"`
	Command     string `json:"command,omitempty"`
	Response    string `json:"response,omitempty"`
	Username    string `json:"username,omitempty"`
	IP          string `json:"ip,omitempty"`
	EventSource string `json:"event_source,omitempty"`
}

type Profile struct {
	Name     string            `json:"name"`
	Port     int               `json:"port"`
	Banner   string            `json:"banner,omitempty"`
	Commands map[string]string `json:"commands"`
}

var (
	logPath  = "/app/logs/lightweight.log"
	profPath = "/app/profiles"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	os.MkdirAll(filepath.Dir(logPath), 0755)

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(f)
		log.SetFlags(0) // Disable timestamp prefix in logs
	} else {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
	}
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(213)+11,
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(254)+1)
}

func randomUser() string {
	return fmt.Sprintf("attacker%d", rand.Intn(9000)+1000)
}

func logJSON(entry LogEntry) {
	entry.Timestamp = now()
	entry.EventSource = "simulated"
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return
	}
	log.Println(string(jsonData)) // filebeat file logs
	fmt.Println(string(jsonData)) // docker stdout logs
}

func simulateProfile(profile Profile) {
	ip := randomIP()
	user := randomUser()

	logJSON(LogEntry{
		Honeypot:  profile.Name,
		EventType: "connection",
		Response:  profile.Banner,
		Username:  user,
		IP:        ip,
	})

	var keys []string
	for k := range profile.Commands {
		keys = append(keys, k)
	}
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	count := rand.Intn(4) + 2 // 2–5 commands
	for i := 0; i < count && i < len(keys); i++ {
		cmd := strings.ReplaceAll(keys[i], "*", " example")
		resp := profile.Commands[keys[i]]
		logJSON(LogEntry{
			Honeypot:  profile.Name,
			EventType: "command",
			Command:   cmd,
			Response:  resp,
			Username:  user,
			IP:        ip,
		})
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	}

	logJSON(LogEntry{
		Honeypot:  profile.Name,
		EventType: "disconnect",
		Response:  "client disconnected",
		Username:  user,
		IP:        ip,
	})
}

func loadProfiles(dir string) ([]Profile, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}

	var profiles []Profile
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		var profile Profile
		if err := json.Unmarshal(data, &profile); err != nil {
			continue
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func main() {
	for {
		profiles, err := loadProfiles(profPath)
		if err != nil || len(profiles) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}

		profile := profiles[rand.Intn(len(profiles))]
		simulateProfile(profile)

		time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)
	}
}
