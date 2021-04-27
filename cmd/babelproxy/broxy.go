package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/ancientlore/flagcfg"
	"github.com/facebookgo/flagenv"
	"github.com/kardianos/service"
)

var (
	conf     config          // global configuration
	svcBroxy service.Service // service object
	//svcLogger service.Logger  // service logger
	ver     bool
	control string
	help    bool
)

// init it called before main
func init() {
	conf.Init()
	cfgFile := LocateConfigFile("BABELPROXY_CONFIG", "babelproxy.config")
	//log.Print("Using config file: ", cfgFile)
	err := conf.Load(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	flag.DurationVar((*time.Duration)(&conf.Timeout), "timeout", time.Duration(conf.Timeout), "HTTP Timeout")
	flag.IntVar(&conf.Cpus, "cpu", conf.Cpus, "Number of CPUs to use")
	flag.IntVar(&conf.PoolSize, "pool", conf.PoolSize, "Size of HTTP connection pool")
	flag.StringVar(&conf.PubAddr, "pubaddr", conf.PubAddr, "Published HTTP service address in swagger file")
	flag.StringVar(&conf.RestAddr, "restaddr", conf.RestAddr, "HTTP service address for the REST endpoints")
	flag.StringVar(&conf.StatusAddr, "statusaddr", conf.StatusAddr, "HTTP service address of the status site")
	flag.StringVar(&conf.BabelAddr, "babeladdr", conf.BabelAddr, "HTTP service address of the remote Babel server")
	flag.StringVar(&conf.BabelProto, "babelproto", conf.BabelProto, "Babel service protocol - http or https")
	flag.StringVar(&conf.CpuProfile, "cpuprofile", conf.CpuProfile, "Write CPU profile to file")
	flag.StringVar(&conf.MemProfile, "memprofile", conf.MemProfile, "Write memory profile to file")
	flag.StringVar(&conf.RestVersion, "restversion", conf.RestVersion, "Set the REST service version number")
	flag.StringVar(&conf.BabelVersion, "babelversion", conf.BabelVersion, "Set the Babel service version number")
	flag.StringVar(&conf.Title, "title", conf.Title, "Service name")
	flag.BoolVar(&ver, "version", false, "Display the version of babelproxy")
	flag.BoolVar(&conf.Log, "log", conf.Log, "Log requests")
	flag.StringVar(&control, "ctl", "", "Service control value - "+strings.Join(service.ControlAction[:], ", "))
	flag.StringVar(&conf.WorkingDir, "wd", conf.WorkingDir, "Set the working directory")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.StringVar(&conf.RestPath, "restpath", conf.RestPath, "Specifies the base path of the REST endpoints, for example /foo/bar")
	flag.StringVar(&conf.BabelPath, "babelpath", conf.BabelPath, "Specifies the base path of the Babel endpoints, for example /foo/bar")
	flag.StringVar(&conf.SwaggerPath, "swaggerpath", conf.SwaggerPath, "Specifies the base path of the swagger endpoint, for example /swagger")
	flag.StringVar(&conf.MediaPath, "mediapath", conf.MediaPath, "Specifies the base path of the media files, for example /media")
	flag.StringVar(&conf.StatusPath, "statuspath", conf.StatusPath, "Specifies the base path of the status page, for example /")
}

// program start
func main() {
	var err error

	sconf := service.Config{
		Name:        "BabelPoxy",
		DisplayName: "Babel Proxy",
		Description: "babelproxy converts RESTful calls into Babel calls.",
	}

	// Parse flags from command-line
	flag.Parse()

	// Parser flags from config
	flagcfg.AddDefaults()
	flagcfg.Parse()

	// Parse flags from environment (using github.com/facebookgo/flagenv)
	flagenv.Prefix = "BABELPROXY_"
	flagenv.Parse()

	// show help if asked
	if help {
		fmt.Fprintln(os.Stderr, `Usage of babelproxy:
  babelproxy [options] <filepatterns>

"For the faasti RESTafarians! (Feel no way, we maas.)"

Example:
  babelproxy *.babel

Options include:`)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
All of the options can be overridden in the configuration file. To override
the location of the configuration file, set the BABELPROXY_CONFIG environment
variable. When running as a service, all options must be specified in the
configuration file or environment variables.`)
		return
	}

	// show version if asked
	if ver {
		fmt.Printf("babelproxy Version %s\n", BABELPROXY_VERSION)
		return
	}

	// default the pubaddr
	if strings.TrimSpace(conf.PubAddr) == "" {
		conf.PubAddr = conf.RestAddr
	}

	// handle kill signals
	go func() {
		// setup cpu profiling if desired
		var f io.WriteCloser
		var err error
		if conf.CpuProfile != "" {
			f, err = os.Create(conf.CpuProfile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
		}

		// Set up channel on which to send signal notifications.
		// We must use a buffered channel or risk missing the signal
		// if we're not ready to receive when the signal is sent.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		// Block until a signal is received.
		s := <-c
		log.Print("Got signal ", s, ", canceling work")

		if conf.CpuProfile != "" {
			log.Print("Writing CPU profile to ", conf.CpuProfile)
			pprof.StopCPUProfile()
			f.Close()
		}

		// write memory profile if configured
		if conf.MemProfile != "" {
			f, err := os.Create(conf.MemProfile)
			if err != nil {
				log.Print(err)
			} else {
				log.Print("Writing memory profile to ", conf.MemProfile)
				pprof.WriteHeapProfile(f)
				f.Close()
			}
		}
		os.Exit(0)
	}()

	// setup number of CPUs
	runtime.GOMAXPROCS(conf.Cpus)

	// set working directory
	if conf.WorkingDir != "" {
		err = os.Chdir(conf.WorkingDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// override configured arguments (file list) with command line
	if flag.NArg() > 0 {
		conf.Args = flag.Args()
	}

	var i svcImpl
	svcBroxy, err = service.New(&i, &sconf)
	if err != nil {
		log.Fatal(err)
	}
	/*
		svcLogger, err = svcBroxy.Logger(nil)
		if err != nil {
			log.Fatal(err)
		}
	*/
	if control != "" {
		err = service.Control(svcBroxy, control)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if service.Interactive() {
		i.doWork()
	} else {
		err = svcBroxy.Run()
		if err != nil {
			//svcLogger.Error(err.Error())
			log.Println(err)
		}
	}
}
