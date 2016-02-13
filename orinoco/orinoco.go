package orinoco

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	"github.com/vail130/orinoco/litmus"
	"github.com/vail130/orinoco/pump"
	"github.com/vail130/orinoco/sieve"
	"github.com/vail130/orinoco/stringutils"
)

type Config struct {
	IsTestEnv    bool
	Port          string                `yaml:"port"`
	MinBatchSize  int                   `yaml:"min_batch_size"`
	MaxBatchDelay int                   `yaml:"max_batch_delay"`
	Streams       [](map[string]string) `yaml:"streams"`
	Triggers      []litmus.Trigger      `yaml:"triggers"`
}

func makeConfig(configPath string) Config {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	yaml.Unmarshal(configData, &config)
	config.IsTestEnv = stringutils.StringToBool(os.Getenv("TEST"))

	return config
}

func startServices(config Config) {
	go litmus.Litmus(config.Triggers)
	eventChannel := make(chan *sieve.Event)
	go pump.Pump(config.MinBatchSize, config.MaxBatchDelay, config.Streams, eventChannel)
	sieve.Sieve(eventChannel)
}

func runWebServer(config Config) {
	r := mux.NewRouter()

	s := r.PathPrefix("/streams").Subrouter()
	s.HandleFunc("/", sieve.GetAllStreamsHandler).Methods("GET")
	s.HandleFunc("/", sieve.DeleteAllStreamsHandler).Methods("DELETE")
	s.HandleFunc("/{stream}/", sieve.GetStreamHandler).Methods("GET")
	s.HandleFunc("/{stream}/", sieve.PostStreamHandler).Methods("POST")

	if config.IsTestEnv {
		r.HandleFunc("/litmus/triggers/evaluate/", litmus.PutEvaluateLitmusTriggersHandler).Methods("PUT")
	}

	http.Handle("/", r)

	port := stringutils.Concat(":", config.Port)

	err := http.ListenAndServe(port, nil)
	if err != nil && err != syscall.EPIPE {
		log.Fatalln(err)
	}
}

func Orinoco(configPath string) {
	config := makeConfig(configPath)
	startServices(config)
	runWebServer(config)
}
