class MockTcpSocket {
  constructor() {
    this.callbacks = {};
  }
  on(event, cb) {
    this.callbacks[event] = cb;
  }
  write(data) {
     const req = JSON.parse(data);
     if (req.method === "superdaw_transport_control" || req.method === "tools/call") {
         setTimeout(() => {
             if (this.callbacks['data']) {
                 this.callbacks['data'](JSON.stringify({
                     jsonrpc: "2.0",
                     id: req.id,
                     result: "Success"
                 }) + "\n");
             }
         }, 0);
     }
  }
  destroy() {
      if (this.callbacks['close']) this.callbacks['close']();
  }
}

module.exports = {
  createConnection: jest.fn((opts, cb) => {
    const s = new MockTcpSocket();
    if (cb) setTimeout(cb, 0);
    return s;
  }),
};
