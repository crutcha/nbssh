# nbssh
Concurrent SSH runner backed by NetBox inventory

## Setup

nbssh relies on 2 environment variables:
* `NETBOX_HOST`: netbox server IE: https://netboxdemo.com
* `NETBOX_API_TOKEN`: API token for auth

Set these in either `~/.bashrc` (linux) or `~/.bash_profile` (macOS)

#### Source
To build from source:

```
cd ~
git clone https://github.com/crutcha/nbssh.git
cd nbssh
go build
```

#### Download
Pre-compiled binaries located here: [Releases](https://github.com/crutcha/nbssh/releases)

Once downloaded, unzip and place into a folder on $PATH

```
wget https://github.com/crutcha/nbssh/releases/download/v0.2/nbssh-v{tag}-{os}-{arch}.tar.gz
tar -zxvf nbssh-v{tag}-{os}-{arch}.tar.gz
sudo mv nbssh /usr/bin
```

#### Install via HomeBrew (MacOS)

todo

## Usage
```
drew@test-vm[~/nbssh] (master) ×  ❯ nbssh
usage: nbssh [<flags>] <command>

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose            Verbose mode.
      --site=SITE          Site
      --tenant=TENANT      Tenant
      --role=ROLE          Role
      --status=STATUS      Status
      --manufacturer=MANUFACTURER
                           Vendor
      --customfield=CUSTOMFIELD ...
                           Custom Field definition as key-value pair IE: core=something
      --concurrency=10     Concurrent SSH runners
  -c, --confirm            Confirm device list before execution
      --username=USERNAME  Username. Defaults to logged in user
      --password=PASSWORD  Password. Defaults to SSH key

Args:
  <command>  Command
```

Multiple parameters can be passed in as comma-seperated value, IE: `--site=site-1,site-2`. Mutliple custom field values can also be used. Custom fields flag expects key/value pair, for example: `--customfield first=this --customfield second=that`.

## Example

```
drew@test-vm[~/nbssh] (master) ×  ❯ nbssh --site testsite --role testrole 'show system uptime'
Executing against:  [test-device-1 test-device-2]
#########################################################################################
test-device-1
#########################################################################################
fpc0:
--------------------------------------------------------------------------
Current time: 2020-05-06 17:05:39 UTC
System booted: 2019-10-29 20:46:00 UTC (27w0d 20:19 ago)
Protocols started: 2019-10-29 20:51:34 UTC (27w0d 20:14 ago)
Last configured: 2020-03-26 21:47:38 UTC (5w5d 19:18 ago) by testuser
 5:05PM  up 189 days, 20:20, 0 users, load averages: 0.99, 0.79, 0.67

#########################################################################################
test-device-2
#########################################################################################
fpc0:
--------------------------------------------------------------------------
Current time: 2020-05-06 17:05:40 UTC
System booted: 2019-10-25 22:41:49 UTC (27w4d 18:23 ago)
Protocols started: 2019-10-25 22:47:24 UTC (27w4d 18:18 ago)
Last configured: 2020-03-25 21:47:07 UTC (5w6d 19:18 ago) by testuser
 5:05PM  up 193 days, 18:24, 0 users, load averages: 0.58, 0.53, 0.49
```
