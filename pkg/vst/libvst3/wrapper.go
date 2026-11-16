package libvst3

/*
#cgo CXXFLAGS: -std=c++11 -I.
#cgo LDFLAGS: -ldl

#include <stdlib.h>
#include <stdint.h>

// Forward declaration of the C API bridging to the VST3 C++ SDK
int scan_vst3_parameters(const char* plugin_path, char** out_json);
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

type VST3Parameter struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Units string `json:"units"`
}

// DeepScan loads the VST3 module and retrieves its full parameter tree.
func DeepScan(pluginPath string) ([]VST3Parameter, error) {
	cPath := C.CString(pluginPath)
	defer C.free(unsafe.Pointer(cPath))

	var cJson *C.char

	// Simulate the C++ bridge call. In a real environment, this links against
	// the actual VST3 SDK via a C wrapper layer.
	// For cross-compilation and CI safety in this commit, we stub the C call implementation logic.
	res := C.int(0) // C.scan_vst3_parameters(cPath, &cJson)

	if res != 0 {
		return nil, fmt.Errorf("failed to load VST3 module or instantiate edit controller: %s", pluginPath)
	}

	if cJson == nil {
		return nil, fmt.Errorf("stubbed C++ bridge: native libvst3 C bindings not linked in this build")
	}

	jsonStr := C.GoString(cJson)
	C.free(unsafe.Pointer(cJson))

	var params []VST3Parameter
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return nil, err
	}

	return params, nil
}
