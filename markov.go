package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

type Markov struct {
	Chain map[string][]string
	Mutex sync.RWMutex
}

func NewMarkov() *Markov {
	return &Markov{Chain: make(map[string][]string)}
}

func (m Markov) AddLine(line string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	words := strings.Split(line, " ")
	// Due to the need to iterate a maximum of (slice length - 2), we use a for
	// loop instead of a range. If this can be done with range, that'd be great.
	for i := 0; i < (len(words) - 3); i++ {
		prefix := strings.Join(words[i:i+2], " ")
		suffix := words[i+2]
		m.Chain[prefix] = append(m.Chain[prefix], suffix)
	}
}

// Could use gob, but that makes it really hard to change anything in the
// way the markov chains work without potentially destroying the validity
// of the memory. Unrelatedly, I've destroyed the validity of the memory
// before. Coincidence.
func (m Markov) Remember(filePath string) {
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			m.AddLine(scanner.Text())
		}
	}
}

// Save brain.
func (m Markov) Commit(filePath string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
}

func (m Markov) GenerateLine(prefix string, length int) string {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	log.Println("genline")
	var sentence []string
	var suffix string
	sentence = append(sentence, strings.Split(prefix, " ")[0], strings.Split(prefix, " ")[1])
	for i := 0; i < length; i++ {
		suffix = m.Suffix(strings.Join(sentence[i:i+2], " "))
		sentence = append(sentence, suffix)
	}

	log.Println(sentence)
	return strings.Join(sentence, " ")
}

func (m Markov) Suffix(prefix string) string {
	suffixes := m.Chain[prefix]
	dice := rand.Intn(len(suffixes))
	return suffixes[dice]
}
