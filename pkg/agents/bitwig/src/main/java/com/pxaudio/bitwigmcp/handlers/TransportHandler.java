package com.pxaudio.bitwigmcp.handlers;

import com.bitwig.extension.controller.api.Transport;
import com.bitwig.extension.controller.api.ControllerHost;
import com.google.gson.*;
import com.pxaudio.bitwigmcp.BitwigMCPExtension;

public class TransportHandler {
    private final Transport transport;

    public TransportHandler(BitwigMCPExtension extension, ControllerHost host) {
        this.transport = host.createTransport();
    }

    public JsonElement handle(String action, JsonObject params) {
        switch (action) {
            case "play":
                transport.play();
                break;
            case "stop":
                transport.stop();
                break;
            case "record":
                transport.record();
                break;
            case "set_tempo":
                if (params.has("bpm")) {
                    transport.tempo().setRaw(params.get("bpm").getAsDouble());
                }
                break;
            default:
                throw new IllegalArgumentException("Unknown transport action: " + action);
        }
        JsonObject result = new JsonObject();
        result.addProperty("success", true);
        return result;
    }
}
