package config

import "fmt"

// NoteEntry maps a port to a descriptive note.
type NoteEntry struct {
	Port int    `yaml:"port"`
	Note string `yaml:"note"`
}

// NoteMapConfig holds operator-defined notes for specific ports.
type NoteMapConfig struct {
	Entries []NoteEntry `yaml:"entries"`
}

// defaultNoteMapConfig returns an empty NoteMapConfig.
func defaultNoteMapConfig() NoteMapConfig {
	return NoteMapConfig{}
}

// Validate checks that all entries have valid ports and non-empty notes.
func (c NoteMapConfig) Validate() error {
	for _, e := range c.Entries {
		if e.Port < 1 || e.Port > 65535 {
			return fmt.Errorf("notemap: port %d is out of range", e.Port)
		}
		if e.Note == "" {
			return fmt.Errorf("notemap: note for port %d must not be empty", e.Port)
		}
	}
	return nil
}
