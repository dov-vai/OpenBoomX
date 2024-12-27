<h1 align="center">Open Boom X</h1>

Earfun UBoom X speaker controls for the desktop.

# Protocol

The protocol documentation can be viewed [here](protocol.md). 

This was reverse engineered from the Earfun App so there might be some innacuracies on the description of it.

# GUI

<p align="center">
  <img src="https://github.com/user-attachments/assets/7c788c17-0188-4547-b61e-394c149febbd" width="32%"></img>
  <img src="https://github.com/user-attachments/assets/dc96bd0d-57df-4493-a8d3-89d4fed1d034" width="32%"></img>
  <img src="https://github.com/user-attachments/assets/20ec7b4b-3c57-4c41-8abd-4f453c5bcdae" width="32%"></img>
  <img src="https://github.com/user-attachments/assets/67ce438d-b656-41b2-af23-a19a3e05178e" width="32%"></img>
  <img src="https://github.com/user-attachments/assets/da0e2bd0-9d50-439a-ba32-61101cab605f" width="32%"></img>
</p>

# CLI

Command line interface is also available:
```
Usage of ./OpenBoomX:
  -custom string
        Send custom hex message (advanced)
  -eq string
        Set custom eq bands: 10 comma separated values from 0 (-10 dB) to 120 (+10dB). E.g. 0,0,0,0,0,0,0,0,0,0
  -light string
        Set light action: 'default', 'off', or RGB hex value
  -oluv string
        Set EQ mode: 'studio', 'indoor', 'indoor+', 'outdoor', 'outdoor+', 'boom', 'ground'
  -pairing string
        Enable or disable Bluetooth pairing: 'on' or 'off'
  -poweroff
        Power off the speaker
  -shutdown string
        Set shutdown timeout: '5m', '10m', '30m', '60m', '90m', '120m', 'no'
  -solid
        Set if the light should be solid. Otherwise it will dance. Must be used with -light.
  -volume int
        Set beep volume: 0, 25, 50, 75, 100 (default -1)
```

# Building

Install Golang, inside [OpenBoomX/gui](OpenBoomX/gui) or [OpenBoomX/cli](OpenBoomX/cli) run:
```
go build
```

# Contributing

Contributions welcome. No special requirements yet, just open a pull request.

# LICENSE

GNU General Public License 3.0 or later.

See [LICENSE](LICENSE) for the full text.

