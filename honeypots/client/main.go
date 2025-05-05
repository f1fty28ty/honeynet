package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
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

type SimUser struct {
	Username string
	Password string
	SourceIP string
}

var simUsers = []SimUser{
	{"alice", "apple", "192.0.2.10"},
	{"bob", "banana", "198.51.100.5"},
	{"carol", "cherry", "203.0.113.7"},
	{"eve", "password123", "145.100.123.89"},
	{"dave", "hunter2", "212.10.33.51"},
	{"oscar", "letmein", "83.55.111.19"},
}

var ftpCmds = []string{
	"USER %s", "PASS %s", "LIST", "RETR flag.txt", "QUIT",
}
var httpCmds = []string{
	"GET /", "GET /admin", "POST /login", "GET /flag", "HEAD /",
}
var smtpCmds = []string{
	"HELO mail.client.local", "MAIL FROM:<%s@example.com>", "RCPT TO:<target@victim.com>", "DATA", "QUIT",
}
var dnsCmds = []string{
	"QUERY A attacker1.ru", "QUERY TXT bad-domain.xyz", "QUERY MX suspicious.biz",
}

var targets = []string{
	"attacker1.ru:21", "192.168.57.23:80", "malicious.cloud:2525", "evilnode.net:53",
	"leakyhost.org:8080", "phishgate.biz:25",
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func logEntry(entry LogEntry) {
	entry.Timestamp = now()
	entry.Honeypot = "client"
	entry.EventSource = "simulated"

	data, err := json.Marshal(entry)
	if err == nil {
		fmt.Println(string(data)) // stdout for Docker
		log.Println(string(data)) // log file for Filebeat
	}
}

func simulateEvent(user SimUser) {
	serviceType := []string{"ftp", "http", "smtp", "dns"}[rand.Intn(4)]
	var cmdList []string

	switch serviceType {
	case "ftp":
		cmdList = ftpCmds
	case "http":
		cmdList = httpCmds
	case "smtp":
		cmdList = smtpCmds
	case "dns":
		cmdList = dnsCmds
	}

	used := map[int]bool{}
	commands := make([]string, 0, 3)
	for len(commands) < rand.Intn(3)+2 {
		i := rand.Intn(len(cmdList))
		if used[i] {
			continue
		}
		used[i] = true
		cmd := cmdList[i]
		if strings.Contains(cmd, "%s") {
			cmd = fmt.Sprintf(cmd, user.Username)
		}
		commands = append(commands, cmd)
	}

	target := targets[rand.Intn(len(targets))]

	for _, cmd := range commands {
		resp := "OK"
		if rand.Float64() < 0.2 {
			resp = "Error"
		}
		logEntry(LogEntry{
			EventType: serviceType,
			Command:   cmd,
			Response:  resp,
			Username:  user.Username,
			IP:        user.SourceIP,
		})
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	}

	logEntry(LogEntry{
		EventType: serviceType,
		Command:   "disconnect",
		Response:  fmt.Sprintf("Disconnected from %s", target),
		Username:  user.Username,
		IP:        user.SourceIP,
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())
	os.MkdirAll("/app/logs", 0755)
	logFile, err := os.OpenFile("/app/logs/client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(logFile)
		log.SetFlags(0) // Disable default timestamp and prefix
		defer logFile.Close()
	}

	for {
		user := simUsers[rand.Intn(len(simUsers))]
		simulateEvent(user)
		time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
	}
}
