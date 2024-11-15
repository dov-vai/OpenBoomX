# Earfun UBoom X Protocol

RFCOMM protocol, port 2

```python
# Simple python example to test the commands
import socket

mac = "FF:FF:FF:FF:FF:FF" # replace with your device MAC address
command = 0xef000000fe # Replace with yours after 0x

s = socket.socket(socket.AF_BLUETOOTH, socket.SOCK_STREAM, socket.BTPROTO_RFCOMM)
s.connect((mac, 2))
byte_length = (command.bit_length() + 7) // 8
s.send(command.to_bytes(byte_length, byteorder='big'))
s.close()
```

# Oluv's EQ

**7 byte packets** to configure Oluv's EQ modes.

| Mode       | Data            |
|------------|------------------|
| Studio     | `efb046010102fe` |
| Indoor     | `efb046010203fe` | 
| Indoor+    | `efb046010304fe` |
| Outdoor    | `efb046010405fe` |
| Outdoor+   | `efb046010506fe` |
| Boom XXX   | `efb046010607fe` |
| Ground O   | `efb046010708fe` |

# Lights

**10 byte packets** control the RGB lighting settings.

- **Note**: The byte before `fe` at the end doesn’t appear to serve any purpose (not a checksum). Changing it has no effect, so `00` is used as the default.

| Action            | Packet                      |
|-------------------|-----------------------------|
| Default Dancing   | `efb095040000000000fe`      |
| Turn Off          | `efb095040100000000fe`      |

### Crafting RGB Packets

RGB packet format:
- Prefix: `efb09504`
- Type: `01` for solid or `02` for dancing
- RGB value in hex
- End: `00fe`

#### Example
To set a dancing white color:
- Packet: `efb0950402ffffff00fe`

# EQ

**17 byte packets** for EQ configuration.

### Crafting EQ Packets

- Prefix: `efb0450b01`
- EQ Bands: `00000000000000000000` (10 bands, 1 byte each)
- End: `00fe`

#### Band Values
- **+10 dB**: `0x78` (120 decimal)
- **0 dB**: `0x3c` (60 decimal)
- **-10 dB**: `0x00` (0 decimal)

Going above `120` (decimal) has no effect; **+10 dB** is the maximum.

# Automatic Shutdown

**7 byte packets** that set the speaker timeout when it's disconnected from a Bluetooth device.

| Mode         | Data             |
|---------------|------------------|
| 5 minutes     | `efb075010102fe` |
| 10 minutes    | `efb075010203fe` | 
| 30 minutes    | `efb075010304fe` |
| 60 minutes    | `efb075010405fe` |
| 90 minutes    | `efb075010506fe` |
| 120 minutes   | `efb075010607fe` |
| No shutdown   | `efb07501ff00fe` |

# Speaker Poweroff

**7 byte packet** that powers off the speaker.

- Packet: `efb025010102fe`

# Bluetooth Pairing

**7 byte packets** that control the Bluetooth pairing modes.

| Mode | Data            |
| -----|------------------|
| On   | `efb035010102fe` |
| Off  | `efb035010001fe` |

# Speaker Beep Volume Control

| Volume | Data            |
|--------|------------------|
| 0%     | `efb065010102fe` |
| 25%    | `efb065010203fe` | 
| 50%    | `efb065010304fe` |
| 75%    | `efb065010405fe` |
| 100%   | `efb065010506fe` |
