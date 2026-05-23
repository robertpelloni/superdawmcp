/**
 * SuperDAW-MCP VST Bridge Example
 * Demonstrates how a native VST plugin can use the SuperDAW C++ SDK
 * to send state updates back to the orchestration core.
 */

#include "../../pkg/client/cpp/SuperDAWClient.hpp"
#include <iostream>
#include <thread>
#include <chrono>

class MockVSTPlugin {
public:
    MockVSTPlugin() : client("./bin/superdaw-mcp") {
        try {
            client.connect();
            std::cout << "VST Bridge: Connected to SuperDAW-MCP core." << std::endl;
        } catch (const std::exception& e) {
            std::cerr << "VST Bridge Connection Failed: " << e.what() << std::endl;
        }
    }

    // Called when a parameter in the VST UI is changed
    void onParameterChanged(int index, float value) {
        std::cout << "VST Param " << index << " -> " << value << std::endl;

        // Example: Map VST knob 0 to DAW track 1 volume
        if (index == 0) {
            client.setMixer("1", value, 0.0f, "ableton");
        }
    }

    void simulate() {
        // Simulate a knob movement
        onParameterChanged(0, 0.75f);
        std::this_thread::sleep_for(std::chrono::milliseconds(500));
        onParameterChanged(0, 0.90f);
    }

private:
    superdaw::Client client;
};

int main() {
    MockVSTPlugin plugin;
    plugin.simulate();
    return 0;
}
