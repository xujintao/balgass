package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

var (
	// ServerEnv env config for server
	ServerEnv configServerEnv

	ServerList []ServerEntry
)

type configServerEnv struct {
	Debug             bool     `envconfig:"DEBUG" default:"false"`
	LogLevel          string   `envconfig:"LOG_LEVEL" default:"info"`
	LogFile           []string `envconfig:"LOG_FILE" default:"-"`
	TCPPort           int      `envconfig:"TCP_PORT" required:"true"`
	UDPPort           int      `envconfig:"UDP_PORT" required:"true"`
	UpdateVersion     string   `envconfig:"UPDATE_VERSION" required:"true"`
	UpdateHostURL     string   `envconfig:"UPDATE_HOST_URL" required:"true"`
	UpdateFTPPort     int      `envconfig:"UPDATE_FTP_PORT" required:"true"`
	UpdateFTPLogin    string   `envconfig:"UPDATE_FTP_LOGIN" required:"true"`
	UpdateFTPPassword string   `envconfig:"UPDATE_FTP_PASSWORD" required:"true"`
	UpdateVersionFile string   `envconfig:"UPDATE_VERSION_FILE" required:"true"`
	ServerListJSON    string   `envconfig:"SERVER_LIST_JSON" required:"true"`
}

type ServerEntry struct {
	Code    int    `json:"code"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Visible bool   `json:"visible"`
	Name    string `json:"name"`
}

func init() {
	// Load environment variables
	ENV(&ServerEnv)
	// Validate port numbers
	for name, port := range map[string]int{"TCP_PORT": ServerEnv.TCPPort, "UDP_PORT": ServerEnv.UDPPort, "UPDATE_FTP_PORT": ServerEnv.UpdateFTPPort} {
		if port < 1 || port > 65535 {
			log.Fatalf("%s: port must be between 1 and 65535", name)
		}
	}
	// Configure logger
	configureLogger()
	// Parse server list JSON
	var err error
	ServerList, err = parseServerList(ServerEnv.ServerListJSON)
	if err != nil {
		log.Fatalf("SERVER_LIST_JSON: %v", err)
	}
}

func ENV(v any) {
	err := envconfig.Process("", v)
	if err != nil {
		log.Fatal(err)
	}
}

func configureLogger() {
	var writes []io.Writer
	for _, s := range ServerEnv.LogFile {
		switch s {
		case "-":
			writes = append(writes, os.Stdout)
		default:
			f, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal(err)
			}
			writes = append(writes, f)
		}
	}
	var l slog.Level
	switch ServerEnv.LogLevel {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	}
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				io.MultiWriter(writes...),
				&slog.HandlerOptions{
					Level: l,
					// AddSource: true,
				},
			),
		),
	)
}

func parseServerList(raw string) ([]ServerEntry, error) {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 2 && raw[0] == '\'' && raw[len(raw)-1] == '\'' {
		raw = raw[1 : len(raw)-1]
	}
	var entries []ServerEntry
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("server list is empty")
	}
	seen := make(map[int]bool, len(entries))
	for _, entry := range entries {
		if entry.Code < 0 || entry.Code > 65535 || entry.IP == "" || entry.Port < 1 || entry.Port > 65535 || entry.Name == "" || seen[entry.Code] {
			return nil, fmt.Errorf("invalid or duplicate server code %d", entry.Code)
		}
		seen[entry.Code] = true
	}
	return entries, nil
}
