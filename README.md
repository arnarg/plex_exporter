Plex exporter
-------------

```
NAME:
   plex_exporter - A Prometheus exporter that exports metrics on Plex Media Server.

USAGE:
   plex_exporter [global options] command [command options] [arguments...]

COMMANDS:
   token, t  Get authentication token from plex.tv
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config-path value, -c value     Path to config file (default: "/etc/plex_exporter/config.yaml")
   --listen-address value, -l value  Port for server (default: ":9594")
   --log-level value                 Verbosity level of logs (default: "info")
   --format value, -f value          Output format of logs (default: "text")
   --auto-discover, -a               Auto discover Plex servers from plex.tv
   --plex-server value, -p value     Address of Plex Media Server
   --token value, -t value           Authentication token for Plex Media Server
   --help, -h                        show help
   --version, -v                     print the version
```

## Metrics

```
# HELP plex_library_section_size_count Number of items in a library section
# TYPE plex_library_section_size_count gauge
plex_library_section_size_count{name="Movies",server_id="asdf1234",server_name="myplexserver",type="movie"} 598
plex_library_section_size_count{name="TV Shows",server_id="asdf1234",server_name="myplexserver",type="show"} 31
# HELP plex_server_info Information about Plex server
# TYPE plex_server_info counter
plex_server_info{platform="Linux",server_id="asdf1234",server_name="myplexserver",version="1.16.6.1592-b9d49bdb7"} 1
# HELP plex_sessions_active_count Number of active Plex sessions
# TYPE plex_sessions_active_count gauge
plex_sessions_active_count{server_id="asdf1234",server_name="myplexserver"} 1
```

## Running

Before application can be run an authentication token is needed from plex.tv. This can be acquired by running `plex_exporter token`.

Plex exporter can then be run with `plex_exporter -t <auth_token> --auto-discover` and Plex servers will be auto discovered from plex.tv (this will only include servers you own).

### Without auto discovery

If you don't want to auto discover from plex.tv you can provide a base url to the server with `--plex-server`.

### Config file

`--plex-server` flag only supports specifying a single Plex server. If you want to provide more you need to use the config file.

```yaml
address: ":9594"
logLevel: "info"
logFormat: "text"
autoDiscover: false
token: "asdf1234"
servers:
- baseUrl: https://my.plexserver.io:32400
  insecure: true
# My friend trusts me a lot
- baseUrl: https://myfriends.plexserver.io:32400
  token: my-friends-token
```

Servers without a token use the token from top-level config.

The insecure key turns off tls verify for that server.
