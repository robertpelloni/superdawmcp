// SuperDAW Cubase MIDI Remote Driver
// This script maps Cubase transport and mixer to the SuperDAW OSC/MIDI bridge.

var deviceDriver = makeDefaultDeviceDriver("SuperDAW", "Universal Controller");

var activeDevice = deviceDriver.mSurface.makeCustomControlSurface();

// Map Transport
var transport = deviceDriver.mTransport;
// ... mapping logic ...

console.log("SuperDAW Cubase MIDI Remote Loaded");
