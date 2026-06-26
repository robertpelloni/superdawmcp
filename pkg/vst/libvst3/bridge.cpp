#include <stdlib.h>
#include <string.h>

extern "C" {
    // This is the stub implementation of the C++ bridge to the VST3 SDK.
    // In a real integration, this includes <pluginterfaces/vst/ivsteditcontroller.h>
    // and uses GetPluginFactory()->createInstance(...)
    int scan_vst3_parameters(const char* plugin_path, char** out_json) {
        // Stub: do not actually link the heavy Steinberg SDK for this commit.
        *out_json = nullptr;
        return 0;
    }
}
