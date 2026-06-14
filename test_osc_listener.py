import socket
import struct
import time

def _pad4(data):
    extra = (4 - len(data) % 4) % 4
    return data + b'\x00' * extra

def _decode_osc_packet(data):
    if data.startswith(b'#bundle'):
        pos = 8
        pos += 8
        while pos < len(data):
            size = struct.unpack('>I', data[pos:pos + 4])[0]
            pos += 4
            for msg in _decode_osc_packet(data[pos:pos + size]):
                yield msg
            pos += size
    else:
        addr_end = data.find(b'\x00')
        if addr_end < 0:
            return
        address = data[:addr_end].decode('utf-8', errors='replace')
        pos = addr_end + 1
        pos = ((pos + 3) // 4) * 4
        if pos >= len(data) or data[pos:pos + 1] != b',':
            yield (address, [])
            return
        tag_end = data.find(b'\x00', pos)
        if tag_end < 0:
            yield (address, [])
            return
        tags = data[pos + 1:tag_end]
        pos = tag_end + 1
        pos = ((pos + 3) // 4) * 4
        args = []
        for t in tags:
            if t == ord('i'):
                val = struct.unpack('>i', data[pos:pos + 4])[0]
                args.append(val)
                pos += 4
            elif t == ord('f'):
                val = struct.unpack('>f', data[pos:pos + 4])[0]
                args.append(val)
                pos += 4
            elif t == ord('s'):
                s_end = data.find(b'\x00', pos)
                val = data[pos:s_end].decode('utf-8', errors='replace')
                pos = ((s_end + 4) // 4) * 4
                args.append(val)
