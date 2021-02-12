# Port scanner

What can be funnier and easier than pentesting your new smart home device? (Mine've got 5 open ports)

Usage _(we want to find out what ports on IP address 3.14.15.92 are available for TCP (on by default) and UPD (-u flag))_:
```bash
go run main.go -i "3.14.15.92" -u
```