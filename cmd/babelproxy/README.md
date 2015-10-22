![logo](media/logo.png)

BABELPROXY
==========

babelproxy is a server that converts REST calls to Babel calls automatically.

	Usage of babelproxy:
	  babelproxy [options] <filepatterns>

	"For the faasti RESTafarians! (Feel no way, we maas.)"

	Example:
	  babelproxy *.babel

	Options include:
	  -babeladdr string
	    	HTTP service address of the remote Babel server (default "localhost")
	  -babelpath string
	    	Specifies the base path of the Babel endpoints, for example /foo/bar (default "/")
	  -babelproto string
	    	Babel service protocol - http or https (default "http")
	  -babelversion string
	    	Set the Babel service version number (default "1.0")
	  -cpu int
	    	Number of CPUs to use (default 1)
	  -cpuprofile string
	    	Write CPU profile to file
	  -ctl string
	    	Service control value - start, stop, restart, install, uninstall
	  -help
	    	Show help
	  -log
	    	Log requests
	  -mediapath string
	    	Specifies the base path of the media files, for example /media (default "/media")
	  -memprofile string
	    	Write memory profile to file
	  -pool int
	    	Size of HTTP connection pool (default 10)
	  -pubaddr string
	    	Published HTTP service address in swagger file
	  -restaddr string
	    	HTTP service address for the REST endpoints (default "localhost:9999")
	  -restpath string
	    	Specifies the base path of the REST endpoints, for example /foo/bar (default "/")
	  -restversion string
	    	Set the REST service version number (default "1.1")
	  -statusaddr string
	    	HTTP service address of the status site (default "localhost:9998")
	  -statuspath string
	    	Specifies the base path of the status page, for example / (default "/")
	  -swaggerpath string
	    	Specifies the base path of the swagger endpoint, for example /swagger (default "/swagger")
	  -timeout duration
	    	HTTP Timeout (default 10s)
	  -title string
	    	Service name (default "My Service")
	  -version
	    	Display the version of babelproxy
	  -wd string
	    	Set the working directory

	All of the options can be overridden in the configuration file. To override
	the location of the configuration file, set the BABELPROXY_CONFIG environment
	variable. When running as a service, all options must be specified in the
	configuration file or environment variables.
