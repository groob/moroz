<p align="center">
<img src="moroz.png" alt="moroz"/><br/>
</p>

Moroz is a server for the [Santa](https://github.com/google/santa) project.

> Santa is a binary whitelisting/blacklisting system for macOS. It consists of a kernel extension that monitors for executions, a userland daemon that makes execution decisions based on the contents of a SQLite database, a GUI agent that notifies the user in case of a block decision and a command-line utility for managing the system and synchronizing the database with a server.

> Santa is a project of Google's Macintosh Operations Team.

# Configurations

Moroz uses [TOML](https://github.com/toml-lang/toml#example) rule files to specify configuration for Santa.
The path to the folder with the configurations can be specified with `-configs /path/to/configs`.

Moroz expects a `global.toml` file which contains a list of rules. The `global` config can be overriden by providing a machine specific config. 
To do so, name the file for each host with the santa `machine id` [configuration parameter](https://github.com/google/santa/wiki/Configuration#keys-to-be-used-with-a-tls-server). By default, this is the hardware UUID of the mac.

Below is a sample configuration file:

```
client_mode = "MONITOR"
#blacklist_regex = "^(?:/Users)/.*"
#whitelist_regex = "^(?:/Users)/.*"
batch_size = 100

[[rules]]
rule_type = "BINARY"
policy = "BLACKLIST"
sha256 = "2dc104631939b4bdf5d6bccab76e166e37fe5e1605340cf68dab919df58b8eda"
custom_msg = "blacklist firefox"

[[rules]]
rule_type = "CERTIFICATE"
policy = "BLACKLIST"
sha256 = "e7726cf87cba9e25139465df5bd1557c8a8feed5c7dd338342d8da0959b63c8d"
custom_msg = "blacklist dash app certificate"
```

# Creating rules

Acceptable values for client mode:
```
MONITOR | LOCKDOWN
```

Values for `rule_type`:
```
BINARY | CERTIFICATE
```

Values for `policy`:
```
BLACKLIST | WHITELIST
```

use the santactl command to get the sha256 value: 
```
santactl fileinfo /Applications/Firefox.app
```

# Build

The commands below assume you have `$GOPATH/bin` in your path.

```
go get -u github.com/golang/dep
dep ensure
cd cmd/moroz; go install; cd -
```

# Run

`moroz`  
See `moroz -h` for a full list of options.

```
Usage of moroz:
  -configs string
    	path to config folder (default "../../configs")
  -event-logfile string
    	path to file for saving uploaded events (default "/tmp/santa_events")
  -http-addr string
    	http address ex: -http-addr=:8080 (default ":8080")
  -tls-cert string
    	path to TLS certificate (default "server.crt")
  -tls-key string
    	path to TLS private key (default "server.key")
  -version
    	print version information
```

# Quickstart

Download the `moroz` binary from the [Releases](https://github.com/groob/moroz/releases) page.
Copy the `configs` folder from the repo somewhere locally. It must have the `global.toml` file.


Generate a self signed certificate which will be used by santa and the server for communication.

```
openssl genrsa -out server.key 2048
openssl rsa -in server.key -out server.key
openssl req -sha256 -new -key server.key -out server.csr -subj "/CN=santa"
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
rm -f server.csr
```

Add the santa CN to your hosts file.

```
sudo echo "127.0.0.1 santa" >> /etc/hosts
```


Install Santa
The latest version of santa is available on the github repo page: https://github.com/google/santa/releases

Configure Santa:
You will need to provide the `SyncBaseURL` and `ServerAuthRootsFile` settings.

```
sudo launchctl unload -w /Library/LaunchDaemons/com.google.santad.plist
sudo defaults write /var/db/santa/config.plist SyncBaseURL https://santa:8080/v1/santa/
sudo defaults write /var/db/santa/config.plist ServerAuthRootsFile $(pwd)/server.crt
sudo launchctl load -w /Library/LaunchDaemons/com.google.santad.plist
```

Start moroz:
Assumes you have the `./server.crt` and `./server.key` files.

```moroz -configs /path/to/configs/folder```
