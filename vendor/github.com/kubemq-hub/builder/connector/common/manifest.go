package common

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Manifest struct {
	Schema  string     `json:"schema"`
	Version string     `json:"version"`
	Sources Connectors `json:"sources"`
	Targets Connectors `json:"targets"`
}

func NewManifest() *Manifest {
	return &Manifest{}
}
func LoadManifest(data []byte) (*Manifest, error) {
	m := &Manifest{}
	err := json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
func LoadManifestFromFile(filename string) (*Manifest, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadManifest(b)
}

func LoadFromUrl(url string) (*Manifest, error) {
	file, err := ioutil.TempFile("./", "mfx")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove(file.Name())
		fmt.Println(err)
	}()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	return LoadManifestFromFile(file.Name())
}
func (m *Manifest) Save(filename string) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}
func (m *Manifest) Marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

func (m *Manifest) SetSchema(value string) *Manifest {
	m.Schema = value
	return m
}

func (m *Manifest) SetVersion(value string) *Manifest {
	m.Version = value
	return m
}
func (m *Manifest) SetSourceConnectors(value Connectors) *Manifest {
	m.Sources = value
	return m
}
func (m *Manifest) SetTargetConnectors(value Connectors) *Manifest {
	m.Targets = value
	return m
}
func (m *Manifest) AddConnector(value *Connector) *Manifest {
	m.Sources = append(m.Sources, value)
	return m
}
