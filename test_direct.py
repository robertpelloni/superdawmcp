import socket
import struct
import time

def _pad4(data):
    extra = (4 - len(data) % 4) % 4
    return data + b'\x00' * extra

# Create UDP socket and send to SuperDAW on 11000
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

# Send tempo message
addr = b'/superdaw/transport/tempo'
addr_padded = _pad4(addr + b'\x00')
tags = b',f'
tags_padded = _pad4(tags + b'\x00')
val = struct.pack('>f', 148.0)
msg = addr_padded + tags_padded + val
sock.sendto(msg, ('127.0.0.1', 11000))
print(f'Sent tempo: {msg.hex()}')

# Send track create
addr2 = b'/superdaw/track/create'
addr2_padded = _pad4(addr2 + b'\x00')
tags2 = b',ss'
tags2_padded = _pad4(tags2 + b'\x00')
s1 = b'Test Track'
s1_padded = _pad4(s1 + b'\x00')
s2 = b'midi'
s2_padded = _pad4(s2 + b'\x00')
msg2 = addr2_padded + tags2_padded + s1_padded + s2_padded
sock.sendto(msg2, ('127.0.0.1', 11000))
print(f'Sent track create: {msg2.hex()}')

sock.close()
print('Done sending messages')
