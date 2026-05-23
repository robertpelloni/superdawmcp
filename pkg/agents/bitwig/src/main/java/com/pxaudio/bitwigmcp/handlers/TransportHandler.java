package com.pxaudio.bitwigmcp.handlers;

import com.bitwig.extension.controller.api.Transport;
import com.bitwig.extension.controller.api.TrackBank;

public class TransportHandler {
    private final Transport transport;

    public TransportHandler(Transport transport) {
        this.transport = transport;
    }

    public void handle(String method, Object params) {
        if ("transport.set_playing".equals(method)) {
            // Simplified handling logic for Bitwig
            transport.play();
        }
    }
}
