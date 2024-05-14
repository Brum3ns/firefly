
<h1 align="center">
  <img src="assets/fireflyLogo.png" alt="firefly" width="220px">
  <br>
</h1>
 
<p align="center">&lt/
  <a href="#advantages">Advantages</a> |
  <a href="#features">Features</a> |
  <a href="#installation">Installation</a> |
  <a href="#usage">Usage</a> |
  <a href="#community">Community</a> &gt
</p>

Firefly is an advanced black-box fuzzer and not just a standard asset discovery tool. Firefly provides the advantage of testing a target with a large number of built-in checks to detect behaviors in the target.

# Advantages
- [x] Hevy use of gorutines and internal hardware for great preformance
- [x] Built-in engine that handles each task for "x" response results inductively
- [x] Highly cusomized to handle more complex fuzzing
- [x] Filter options and request verifications to avoid junk results
- [x] Friendly error and debug output
- [x] Build in payloads (default list are mixed with the wordlist from [seclists](https://github.com/danielmiessler/SecLists))
- [x] Payload tampering and encoding functionality

# Features
<h1 align="center">
  <img src="assets/fireflyOptions.png" alt="fireflyOptions" width="100%">
  <br>
</h1>

# Installation

## Golang

```bash
go install -v github.com/Brum3ns/firefly/cmd/firefly@latest
# or
go get -v github.com/Brum3ns/firefly/cmd/firefly
```

## Docker

```bash
docker build -t firefly:latest .
alias firefly='docker run --rm -it --net=host -v "$PWD:/firefly/" firefly:latest'
firefly -u 'http://example.com/?query=FUZZ/'
```

<!--
If the above install method do not work try the following:
```
git clone https://github.com/Brum3ns/firefly.git
cd firefly/
go build cmd/firefly/firefly.go
./firefly -h
```
-->


# Usage

### Simple

```bash
firefly -h
```

```bash
firefly -u 'http://example.com/?query=FUZZ'
```

---

## Advanced usage

### Request
Different types of request input that can be used

Basic
```bash
firefly -u 'http://example.com/?query=FUZZ' --timeout 7000
```

Request with different methods and protocols
```bash
firefly -u 'http://example.com/?query=FUZZ' -m GET,POST,PUT -p https,http,ws
```

#### Pipeline
```bash
echo 'http://example.com/?query=FUZZ' | firefly 
```

#### HTTP Raw
```bash
firefly -r '
GET /?query=FUZZ HTTP/1.1
Host: example.com
User-Agent: FireFly'
```

This will send the HTTP Raw  and auto detect all GET and/or POST parameters to fuzz.
```bash
firefly -r '
POST /?A=1 HTTP/1.1
Host: example.com
User-Agent: Firefly
X-Host: FUZZ

B=2&C=3' -au replace
```

### Request Verifier
Request verifier is the most important part. This feature let Firefly know the core behavior of the target your fuzz. It's important to do quality over quantity. More verfiy requests will lead to better quality at the cost of internal hardware preformance (*depending on your hardware*)

```bash
firefly -u 'http://example.com/?query=FUZZ' -e 
```

### Payloads
Payload can be highly customized and with a good core wordlist it's possible to be able to fully adapt the payload wordlist within Firefly itself.

#### Payload debug
> Display the format of all payloads and exit
```bash
firefly -show-payload
```

#### Tampers 
> List of all Tampers avalible
```bash
firefly -list-tamper
```

Tamper all paylodas with given type (*More than one can be used separated by comma*)
```bash
firefly -u 'http://example.com/?query=FUZZ' -e s2c
```

#### Encode
```bash
firefly -u 'http://example.com/?query=FUZZ' -e hex
```
Hex then URL encode all payloads
```bash
firefly -u 'http://example.com/?query=FUZZ' -e hex,url
```

#### Payload regex replace
```bash
firefly -u 'http://example.com/?query=FUZZ' -pr '\([0-9]+=[0-9]+\) => (13=(37-24))'
```
>The Payloads: `' or (1=1)-- -` and `" or(20=20)or "` 
> Will result in: `' or (13=(37-24))-- -`  and `" or(13=(37-24))or "`
> Where the ` => ` (with spaces) inducate the "*replace to*".


### Filters
> Filter options to filter/match requests that include a given rule.

Filter response to **ignore** (filter) `status code 302` and `line count 0`
```bash
firefly -u 'http://example.com/?query=FUZZ' -fc 302 -fl 0
```

Filter responses to **include** (match) `regex`, and `status code 200`
```bash
firefly -u 'http://example.com/?query=FUZZ' -mr '[Ee]rror (at|on) line \d' -mc 200
```

```bash
firefly -u 'http://example.com/?query=FUZZ' -mr 'MySQL' -mc 200
```


### Preformance
> Preformance and time delays to use for the request process

Threads / Concurrency 
```bash
firefly -u 'http://example.com/?query=FUZZ' -t 35
```

Time Delay in millisecounds (ms) for each Concurrency
```bash
FireFly -u 'http://example.com/?query=FUZZ' -t 35 -dl 2000
```

### Wordlists
> Wordlist that contains the paylaods can be added separatly or extracted from a given folder

Single Wordlist with its attack type
```bash
firefly -u 'http://example.com/?query=FUZZ' -w wordlist.txt:fuzz
```

Extract all wordlists inside a folder. Attack type is depended on the suffix `<type>_wordlist.txt`
```bash
firefly -u 'http://example.com/?query=FUZZ' -w wl/
```
Example
> Wordlists names inside folder `wl` :
> 1. fuzz_wordlist.txt
> 2. time_wordlist.txt


### Output
> JSON output is **strongly recommended**. This is because you can benefit from the `jq` tool to navigate throw the result and compare it.

(*If Firefly is pipeline chained with other tools, standard plaintext may be a better choice.*)

Simple plaintext output format
```bash
firefly -u 'http://example.com/?query=FUZZ' -o file.txt
```

JSON output format (*recommended*)
```bash
firefly -u 'http://example.com/?query=FUZZ' -oJ file.json
```

# Community

Everyone in the community are allowed to suggest new features, improvements and/or add new payloads to Firefly just make a pull request or add a comment with your suggestions!
