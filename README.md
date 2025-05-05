Honeynet Simulation (CSE3660 Final Project)

This project is a comprehensive Docker-based honeynet simulation designed to demonstrate and study various types of honeypots using Go-based services. It includes multiple honeypot types, an EFK logging stack (Filebeat + Elasticsearch + Kibana), and supports both autonomous and manual attacker simulations.

⸻

🧪 Honeypots Implemented

Low/Medium Interaction
	•	honeypot_low: Simulates basic FTP and Telnet services.
	•	honeypot_medium: Handles partial FTP commands.

Email
	•	honeypot_spam: Accepts SMTP-like spam payloads.
	•	honeypot_email: Simulates full SMTP exchange.

Server / Database
	•	honeypot_server: Fake HTTP server that returns content.
	•	honeypot_db: Fake SQL server with dummy responses.

Special Purpose
	•	honeypot_malware: Streams rotating fake payloads.
	•	honeypot_production: Mimics an internal admin panel.
	•	honeypot_pure: Simulated interactive shell with fake command output.

Client/Attacker
	•	honeypot_client: Auto-scans and connects to targets.

Targets
	•	malware_server: Fake malware payload source.
	•	exploit_server: Fake exploit trigger host.

⸻

📂 Directory Layout

honeynet/
├── honeypots/
│   ├── client/             # Simulated attacking agent
│   ├── database/           # SQL interaction trap
│   ├── email/              # SMTP honeypot
│   ├── low_interaction/    # Basic TCP traps
│   ├── malware/            # Rotating payloads
│   ├── medium_interaction/ # Limited FTP command support
│   ├── production/         # Fake admin dashboard
│   ├── pure/               # Full interactive shell
│   ├── server/             # HTTP trap
│   └── spam/               # SMTP spam receiver
├── logging/
│   └── filebeat/           # Filebeat setup to forward logs
├── docker-compose.yml      # Main orchestration file
├── network.env             # IP definitions for all services
└── README.md               # Project documentation



⸻

🚀 Getting Started
	1.	Build and Launch the Cluster

docker compose build
docker compose up -d

This will:
	•	Launch all honeypots, clients, and targets.
	•	Start Filebeat, Elasticsearch, and Kibana.
	•	Hook all logs under /honeypots/*/logs to the EFK stack.

	2.	Access Kibana Dashboard
Visit http://localhost:5601 to visualize logs from all honeypots.

⸻

📊 Kibana Dashboards

With the honeynet running and logs flowing:

Step 1: Create a Data View
	1.	Navigate to Stack Management → Data Views.
	2.	Click Create data view.
	3.	In the Name field, enter a descriptive name (e.g., Honeypot Logs).
	4.	In the Index pattern field, enter:

filebeat-*


	5.	Select @timestamp as the time field.
	6.	Click Save data view to Kibana.

Note: In Kibana 8.18, “Index Patterns” have been renamed to “Data Views”.

Step 2: Build Visualizations

Use the following setups:
	•	Pie Chart:
	•	Split slices by honeypot to show traffic per trap.
	•	Bar Chart:
	•	X-axis: Date histogram (@timestamp).
	•	Split series: event_type.
	•	Data Table:
	•	Show top IP addresses by interaction count.

Step 3: Assemble Dashboard
	1.	Navigate to Dashboard → Create new dashboard.
	2.	Add your visualizations.
	3.	Save it as “Honeypot Monitoring Overview”.

This setup allows you to monitor IP activity per honeypot, commands issued, and session behaviors over time.

⸻

🧾 Logging Behavior
	•	Logs are written to /logs within each honeypot container.
	•	All logs are structured JSON and include entries like:

{
  "timestamp": "2025-05-04T17:10:40Z",
  "honeypot": "pure",
  "event_type": "command",
  "command": "whoami",
  "response": "executed",
  "username": "nginx9549",
  "ip": "92.63.197.153",
  "source": "autonomous"
}


	•	Filebeat forwards these logs to Elasticsearch, making them visible via Kibana.

⸻

💡 Notes
	•	All orchestration is handled via Docker Compose; no additional scripts are needed.
	•	Static IPs are assigned per honeypot and target for reliable test routing.
	•	All honeypots run in isolated containers on the honeynet_net bridge network.

⸻

🛡️ What Each Honeypot Tests

Honeypot	Primary Focus
honeypot_low	FTP login attempts, command usage (USER, PASS, ls, get)
honeypot_medium	FTP traversal and file access with partial validation
honeypot_email	Credential stuffing and injection via SMTP commands
honeypot_spam	Bulk SMTP activity and spam-like content detection
honeypot_db	SQL enumeration, fake login, table dumping
honeypot_malware	Simulated C2 commands, payload downloads, and beaconing
honeypot_production	Unauthorized config access (cat /etc/shadow, env)
honeypot_pure	Full terminal interaction, privilege abuse, lateral movement
honeypot_server	Web probing (GET /admin, DELETE /backup)
honeypot_client	Simulated attacker that connects to honeypots

Each honeypot logs structured JSON events to /logs and forwards them via Filebeat for centralized monitoring.
