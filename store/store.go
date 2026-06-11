package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

// ScoreRecord is one completed game's data, persisted to JSON.
type ScoreRecord struct {
	Mode        string    `json:"mode"`
	Score       int       `json:"score"`
	HighStreak  int       `json:"highStreak"`
	BallsCaught int       `json:"ballsCaught"`
	BallsMissed int       `json:"ballsMissed"`
	MaxPhase    int       `json:"maxPhase"`
	DurationSec int       `json:"durationSec"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
}

// AggStats are computed live from all score records.
type AggStats struct {
	TotalCaught  int
	TotalTimeSec int
	BestStreak   int
	GamesPlayed  int
}

// Store handles all disk I/O for score history and config.
type Store struct {
	dir       string
	scoresFile string
	configFile string
}

// New creates a Store, creating the data directory if needed.
func New() *Store {
	dir := dataDir()
	_ = os.MkdirAll(dir, 0755)
	return &Store{
		dir:        dir,
		scoresFile: filepath.Join(dir, "scores.json"),
		configFile: filepath.Join(dir, "config.json"),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Scores
// ─────────────────────────────────────────────────────────────────────────────

// LoadAll reads all score records, sorted descending by score.
func (s *Store) LoadAll() ([]ScoreRecord, error) {
	data, err := os.ReadFile(s.scoresFile)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var records []ScoreRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Score > records[j].Score
	})
	return records, nil
}

// Save appends a new record atomically (write tmp → rename).
func (s *Store) Save(r ScoreRecord) error {
	existing, _ := s.LoadAll()
	existing = append(existing, r)

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.scoresFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	// keep a backup of the previous file
	if _, err := os.Stat(s.scoresFile); err == nil {
		_ = os.Rename(s.scoresFile, s.scoresFile+".bak")
	}
	return os.Rename(tmp, s.scoresFile)
}

// Reset deletes all score data.
func (s *Store) Reset() error {
	_ = os.Remove(s.scoresFile + ".bak")
	return os.Remove(s.scoresFile)
}

// HiScore returns the best score for a given mode code (or overall if empty).
func (s *Store) HiScore(mode string) int {
	records, err := s.LoadAll()
	if err != nil || len(records) == 0 {
		return 0
	}
	best := 0
	for _, r := range records {
		if mode == "" || r.Mode == mode {
			if r.Score > best {
				best = r.Score
			}
		}
	}
	return best
}

// Aggregate computes summary statistics across the given records.
func (s *Store) Aggregate(records []ScoreRecord) AggStats {
	var a AggStats
	a.GamesPlayed = len(records)
	for _, r := range records {
		a.TotalCaught += r.BallsCaught
		a.TotalTimeSec += r.DurationSec
		if r.HighStreak > a.BestStreak {
			a.BestStreak = r.HighStreak
		}
	}
	return a
}

// ─────────────────────────────────────────────────────────────────────────────
// Config
// ─────────────────────────────────────────────────────────────────────────────

type Config struct {
	ThemeIndex int `json:"themeIndex"`
}

// LoadConfig returns the saved config or sensible defaults.
func (s *Store) LoadConfig() Config {
	data, err := os.ReadFile(s.configFile)
	if err != nil {
		return Config{}
	}
	var c Config
	_ = json.Unmarshal(data, &c)
	return c
}

// SaveTheme persists the theme index.
func (s *Store) SaveTheme(idx int) {
	c := s.LoadConfig()
	c.ThemeIndex = idx
	data, _ := json.MarshalIndent(c, "", "  ")
	_ = os.WriteFile(s.configFile, data, 0644)
}

// ─────────────────────────────────────────────────────────────────────────────
// Platform helpers
// ─────────────────────────────────────────────────────────────────────────────

func dataDir() string {
	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "paddle-ball")
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".paddle-ball")
	}
	return filepath.Join(home, ".paddle-ball")
}