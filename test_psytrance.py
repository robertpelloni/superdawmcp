import sys, os, json, subprocess, time

sys.path.append('pkg/client/py')

# Start Go server
go_proc = subprocess.Popen(['./bin/superdaw-mcp'], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, bufsize=1)

# Start bridge
bridge_proc = subprocess.Popen(['python', 'scripts/osc_to_abletonmcp_bridge.py'], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)

time.sleep(3)

from superdaw_client import SuperDAWClient

client = SuperDAWClient('./bin/superdaw-mcp')
client.connect()
print('=== Connected ===')

# Test tempo
print('Transport:', client.transport_control(playing=True, bpm=148.0))
time.sleep(0.5)

# Create tracks
print('Bass:', client.create_track(name='Psy Bass', track_type='midi'))
time.sleep(0.5)
print('Lead:', client.create_track(name='Psy Lead', track_type='midi'))
time.sleep(0.5)
print('Pad:', client.create_track(name='Psy Pad', track_type='midi'))
time.sleep(0.5)
print('Perc:', client.create_track(name='Psy Perc', track_type='midi'))
time.sleep(0.5)

# Generate euclidean patterns
print('Kick euclid:', client.generate_euclidean(track_id='0', hits=1, steps=4, pitch=36))
time.sleep(0.5)
print('Hat euclid:', client.generate_euclidean(track_id='1', hits=1, steps=8, pitch=42))

# Print bridge output
print('\n=== Bridge output ===')
while True:
    line = bridge_proc.stdout.readline()
    if not line:
        break
    print(line.rstrip())
    if 'Psy Bass' in line or 'Success' in line:
        break

client.disconnect()
go_proc.terminate()
bridge_proc.terminate()
print('=== Done ===')
