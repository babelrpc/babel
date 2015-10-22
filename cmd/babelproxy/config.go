package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// duration is used to allow us to use custom TOML marshaling.
type duration time.Duration

// UnmarshalText parses a time duration from the config file, so we can write times like "10s"
func (d *duration) UnmarshalText(text []byte) error {
	dur, err := time.ParseDuration(string(text))
	if err == nil {
		*d = duration(dur)
	}
	return err
}

// config holds configuration settings for broxy
type config struct {
	Log          bool     `toml:"-"`    // whether to log
	RestAddr     string   `toml:"-"`    // listen address
	PubAddr      string   `toml:"-"`    // published address
	StatusAddr   string   `toml:"-"`    // status site address
	BabelProto   string   `toml:"-"`    // protocol (http, https)
	BabelAddr    string   `toml:"-"`    // address of babel server
	Cpus         int      `toml:"-"`    // number of CPUs to use
	Timeout      duration `toml:"-"`    // HTTP timeout
	CpuProfile   string   `toml:"-"`    // write cpu profile to file
	MemProfile   string   `toml:"-"`    // write memory profile to file
	WorkingDir   string   `toml:"-"`    // change to this working fir
	PoolSize     int      `toml:"-"`    // size of connection pool
	RestPath     string   `toml:"-"`    // base path of the REST endpoints
	BabelPath    string   `toml:"-"`    // base path of the Babel service
	SwaggerPath  string   `toml:"-"`    // base path of swagger-ui
	MediaPath    string   `toml:"-"`    // base path of media files
	StatusPath   string   `toml:"-"`    // base path of status page
	BabelVersion string   `toml:"-"`    // service version
	RestVersion  string   `toml:"-"`    // service version
	Title        string   `toml:"-"`    // service title
	Args         []string `toml:"args"` // list of file patterns to process
}

// Init sets the defaults for config values that haven't been set
func (c *config) Init() {
	c.Log = false
	c.Args = nil
	c.Timeout = duration(time.Second * 10)
	c.Cpus = 1
	c.RestAddr = "localhost:9999"
	c.PubAddr = ""
	c.StatusAddr = "localhost:9998"
	c.PoolSize = 10
	c.BabelAddr = "localhost"
	c.BabelProto = "http"
	c.RestPath = "/"
	c.BabelPath = "/"
	c.SwaggerPath = "/swagger"
	c.MediaPath = "/media"
	c.StatusPath = "/"
	c.RestVersion = "1.1"
	c.BabelVersion = "1.0"
	c.Title = "My Service"
}

// Load reads the configuration from a file
func (c *config) Load(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("Unable to stat configuration file %s: %s", filename, err))
		} else {
			// no file - okay go on
			log.Printf("No configuration file at \"%s\"", filename)
		}
	} else {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to read configuration file %s: %s", filename, err))
		}
		if _, err = toml.Decode(string(data), c); err != nil {
			return errors.New(fmt.Sprintf("Unable to parse configuration file %s: %s", filename, err))
		}
	}
	return nil
}

// LocateConfigFile returns the assumed location of the configuration file named cfgName, using the standard locations
// and checking the environment variable named envName.
func LocateConfigFile(envName, cfgName string) string {
	var cfgFile string

	if envName != "" {
		cfgFile = os.Getenv(envName)
	}

	if cfgFile == "" {
		var exePath, err = exec.LookPath(os.Args[0])
		if err != nil {
			log.Print("Warning: ", err)
			exePath = os.Args[0]
		}
		s, err := filepath.Abs(exePath)
		if err != nil {
			log.Print("Warning: ", err)
		} else {
			exePath = s
		}
		exePath, _ = filepath.Split(exePath)
		exePath = filepath.ToSlash(exePath)

		if strings.HasPrefix(exePath, "/usr/bin/") {
			exePath = strings.Replace(exePath, "/usr/bin/", "/etc/", 1)
		} else {
			exePath = strings.Replace(exePath, "/bin/", "/etc/", 1)
		}

		cfgFile = filepath.Join(filepath.FromSlash(exePath), cfgName)
	}

	return cfgFile
}
