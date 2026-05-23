#pragma once
#include <iostream>
#include <string>
#include <vector>
#include <cstdio>
#include <memory>
#include <stdexcept>

#ifdef _WIN32
#define POPEN _popen
#define PCLOSE _pclose
#else
#define POPEN popen
#define PCLOSE pclose
#endif

namespace superdaw {

struct MIDINote {
    int pitch;
    int velocity;
    float start_beat;
    float duration;
};

class Client {
public:
    Client(const std::string& serverPath) : serverPath(serverPath) {}

    void connect() {
        fp = POPEN(serverPath.c_str(), "w");
        if (!fp) {
            throw std::runtime_error("Failed to start SuperDAW server process");
        }
    }

    void disconnect() {
        if (fp) {
            PCLOSE(fp);
            fp = nullptr;
        }
    }

    void setMixer(const std::string& trackId, float volume, float pan = 0.0f, const std::string& daw = "") {
        std::string json = "{\"jsonrpc\":\"2.0\",\"method\":\"tools/call\",\"params\":{\"name\":\"superdaw_set_mixer\",\"arguments\":{\"track_id\":\"" + trackId + "\",\"volume\":" + std::to_string(volume) + ",\"pan\":" + std::to_string(pan);
        if (!daw.empty()) json += ",\"daw\":\"" + daw + "\"";
        json += "}},\"id\":1}\n";
        send(json);
    }

    void transportControl(bool playing, double bpm = 120.0, const std::string& daw = "") {
        std::string json = "{\"jsonrpc\":\"2.0\",\"method\":\"tools/call\",\"params\":{\"name\":\"superdaw_transport_control\",\"arguments\":{\"playing\":" + std::string(playing ? "true" : "false") + ",\"bpm\":" + std::to_string(bpm);
        if (!daw.empty()) json += ",\"daw\":\"" + daw + "\"";
        json += "}},\"id\":1}\n";
        send(json);
    }

    ~Client() {
        disconnect();
    }

private:
    void send(const std::string& data) {
        if (fp) {
            fputs(data.c_str(), fp);
            fflush(fp);
        }
    }

    std::string serverPath;
    FILE* fp = nullptr;
};

} // namespace superdaw
