# name=SuperDAW-MCP FL Studio Agent
# SuperDAW FL Studio MIDI Script
# Bridges SuperDAW Go Core to FL Studio API.

import transport
import mixer
import general
import ui
import device
import playlist

# Note: FL Studio scripts are event-driven based on MIDI.
# This script expects a MIDI-to-OSC bridge to send MIDI CC/Note messages.

def OnMidiMsg(event):
    # Mapping MIDI CC to SuperDAW Commands
    # CC 102: Transport Play/Stop
    # CC 103: Tempo (BPM)
    # CC 104: Track Volume

    if event.status == 176: # CC
        if event.data1 == 102:
            if event.data2 > 0: transport.start()
            else: transport.stop()
            event.handled = True
        elif event.data1 == 103:
            # Map 0-127 to 60-187 BPM
            new_bpm = 60 + event.data2
            transport.setTempo(new_bpm)
            event.handled = True
        elif event.data1 == 104:
            # CC 104 Val = Track Index, next message Val = Volume
            pass

def OnIdle():
    # Polling logic for more complex state sync if needed
    pass

def OnUpdateBeat():
    # Send transport state back via a dummy MIDI message or log
    # for an external bridge to pick up
    pass

def OnInit():
    print("SuperDAW-MCP FL Studio Agent Initialized")

def OnDeInit():
    print("SuperDAW-MCP FL Studio Agent Terminated")
