<h1 align="center">Open Boom X</h1>

Earfun UBoom X speaker controls for the desktop.

# Protocol

The protocol documentation can be viewed [here](protocol.md). 

This was reverse engineered from the Earfun App so there might be some innacuracies on the description of it.

# GUI

<p align="center">
  <img src="https://github.com/user-attachments/assets/dbde3bed-186f-4e5f-b6e6-de0eed93a2b0" width="45%"></img>
  <img src="https://github.com/user-attachments/assets/993e820c-8711-406b-9c30-3c2996717f1e" width="45%"></img>
  <img src="https://github.com/user-attachments/assets/9b14267e-e9ab-4584-9d15-e3b1aae703c4" width="45%"></img>
  <img src="https://github.com/user-attachments/assets/db8503af-f5eb-4942-8644-2c3930b2eac1" width="45%"></img>
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

