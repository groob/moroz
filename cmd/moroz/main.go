package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/groob/moroz/santa"

	stdlog "log"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/groob/moroz/moroz"
	"github.com/groob/moroz/santaconfig"
	"github.com/micromdm/go4/version"
)

const openSSLBash = `
Looks like you're missing a TLS certifacte and private key. You can quickly generate one 
by using the commands below:

openssl genrsa -out server.key 2048
openssl rsa -in server.key -out server.key
openssl req -sha256 -new -key server.key -out server.csr -subj "/CN=santa"
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
rm -f server.csr

Add the santa CN to your hosts file.

sudo echo "127.0.0.1 santa" >> /etc/hosts


You also will need to configure santa:

sudo launchctl unload -w /Library/LaunchDaemons/com.google.santad.plist 
sudo defaults write /var/db/santa/config.plist SyncBaseURL https://santa:8080/v1/santa/
sudo defaults write /var/db/santa/config.plist ServerAuthRootsFile $(pwd)/server.crt
sudo launchctl load -w /Library/LaunchDaemons/com.google.santad.plist


The latest version of santa is available on the github repo page:
	https://github.com/google/santa/releases
`

func main() {
	var (
		flTLSCert   = flag.String("tls-cert", envString("MOROZ_TLS_CERT", "server.crt"), "path to TLS certificate")
		flTLSKey    = flag.String("tls-key", envString("MOROZ_TLS_KEY", "server.key"), "path to TLS private key")
		flAddr      = flag.String("http-addr", envString("MOROZ_HTTP_ADDRESS", ":8080"), "http address ex: -http-addr=:8080")
		flConfigs   = flag.String("configs", envString("MOROZ_CONFIGS", "../../configs"), "path to config folder")
		flEvents    = flag.String("event-logfile", envString("MOROZ_EVENTLOG_FILE", "/tmp/santa_events"), "path to file for saving uploaded events")
		flVersion   = flag.Bool("version", false, "print version information")
		flHTTPDebug = flag.Bool("http-debug", false, "enable debug for http(dumps full request)")
		flNoTLS     = flag.Bool("tls-handled-elsewhere", false, "I promise I terminated TLS elsewhere")
	)
	flag.Parse()

	if *flVersion {
		version.PrintFull()
		return
	}

	if _, err := os.Stat(*flTLSCert); !*flNoTLS && os.IsNotExist(err) {
		fmt.Println(openSSLBash)
		os.Exit(2)
	}

	if !validateConfigExists(*flConfigs) {
		fmt.Println("you need to provide at least a 'global.toml' configuration file in the configs folder. See the configs folder in the git repo for an example")
		os.Exit(2)
	}

	logger := log.NewLogfmtLogger(os.Stderr)

	repo := santaconfig.NewFileRepo(*flConfigs)
	var svc santa.Service
	{
		s, err := moroz.NewService(repo, *flEvents)
		if err != nil {
			stdlog.Fatal(err)
		}
		svc = s
		svc = moroz.LoggingMiddleware(logger)(svc)
	}

	endpoints := moroz.MakeServerEndpoints(svc)

	var h http.Handler
	{
		r := mux.NewRouter()
		h = r
		moroz.AddHTTPRoutes(r, endpoints, logger)
		if *flHTTPDebug {
			h = debugHTTPmiddleware(h)
		}
	}

	go func() { fmt.Println("started server") }()

	if *flNoTLS {
		stdlog.Fatal(http.ListenAndServe(*flAddr, h))
	} else {
		stdlog.Fatal(http.ListenAndServeTLS(*flAddr,
			*flTLSCert,
			*flTLSKey,
			h))
	}
}

func validateConfigExists(configsPath string) bool {
	var hasConfig = true
	if _, err := os.Stat(configsPath); os.IsNotExist(err) {
		hasConfig = false
	}
	if _, err := os.Stat(configsPath + "/global.toml"); os.IsNotExist(err) {
		hasConfig = false
	}
	if !hasConfig {
	}
	return hasConfig
}

func envString(key, def string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return def
}

func debugHTTPmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := io.TeeReader(r.Body, os.Stderr)
		r.Body = ioutil.NopCloser(body)
		out, err := httputil.DumpRequest(r, true)
		if err != nil {
			stdlog.Println(err)
		}
		fmt.Println("")
		fmt.Println(string(out))
		fmt.Println("")
		next.ServeHTTP(w, r)
	})
}
