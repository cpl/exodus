![Exodus](./logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/cpl/exodus)](https://goreportcard.com/report/github.com/cpl/exodus)

DNS Exfiltration tool. By setting up a remote server listening for "DNS queries", we can bypass most firewall rules and monitoring.

This does not guarantee perfect stealth! A good network admin and setup will easily spot such traffic and flag it.

> Do not use for malicious or illegal purposes! Use at your own risk.

## Installing

### From source

1. Download the source code `git clone https://github.com/cpl/exodus.git`
2. Run the following: `make`, this will test everything, get dependencies and build the executables
3. Inside the folder `./out/`, you will see the server and client executables

### Releases page

> May not contain latest changes (some of which are bug fixes)!

1. Go to [Exodus Releases](https://github.com/cpl/exodus/releases) on GitHub
2. Download the binaries for your systems

### Using go get

1. `go get cpl.li/go/exodus/cmd/exodus-client`
2. `go get cpl.li/go/exodus/cmd/exodus-server`

## Usage

### Server

You can simply run it after installing it: `exodus-server`

This will run with the default port `53` and default data directory (tmp).

You can set `--port 1453` to something custom and `--data /var/exodus` to whatever you want. There is also a `-v` flag to display logging messages.

When the server will receive a "query" it will store it in the following structure:

`DATADIR/{token}/{count}.out`

For example sending a file in **4** chunks using the token **example** and the server configured with the default temp dir, will result in the following:

```text
/tmp/exodus/example/00000000.out
/tmp/exodus/example/00000001.out
/tmp/exodus/example/00000002.out
/tmp/exodus/example/00000003.out
```

From here you could do something as simple as: `cat * > full.out` to assemble it back.

### Client

The Exodus client **needs** the following flags:

* `--server dns.example.com`, this will be the address where YOU installed the Exodus Server
* `--target normaldomain.com`, this domain will be the "cover up", so set it to something realistic

Other optional flags are:

* `--file something.txt`, by default Exodus Client will use *stdin* as the input source
* `--size 16`, this is how many bytes to send per DNS query, the default is the max
* `--port 1453`, if you set the server to use something other than 53
* `--token example`, use different tokens for different "uploads", this will separate them server side
* `-v`, enables verbose logging





