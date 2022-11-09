# Aerospike Configuration File Parser cli tool

Download the binaries from the Releases page, or `go build` this branch to get the binary.

## Usage

```bash
% ./parser

Usage: ./parser command path [set-value1] [set-value2] [...set-valueX] filename

Commands:
        delete - delete configuration/stanza
        set    - set configuration parameter
        create - create a new stanza

Path: .path.to.item or .path.to.stanza, e.g. .network.heartbeat

Set-value: for the 'set' command - used to specify value of parameter; leave empty to crete no-value param

Example:
        touch new.conf
        ./parser create network.heartbeat new.conf
        ./parser set network.heartbeat.mode mesh new.conf
        ./parser set network.heartbeat.mesh-seed-address-port "172.17.0.2 3000" "172.17.0.3 3000" new.conf
        ./parser create service new.conf
        ./parser set service.proto-fd-max 3000 new.conf
```
