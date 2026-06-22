import TcpSocket from "react-native-tcp-socket";

interface RPCRequest {
	jsonrpc: "2.0";
	method: string;
	params?: Record<string, unknown>;
	id: number;
}

interface RPCError {
	code: number;
	message: string;
}

interface RPCResponse {
	jsonrpc: "2.0";
	result?: unknown;
	error?: RPCError;
	id: number;
}

type ResponseHandler = (resp: RPCResponse) => void;

/**
 * Singleton TCP socket wrapper that communicates with the
 * SuperDAW-MCP JSON-RPC gateway on port 12002.
 */
class TcpGateway {
	private socket: ReturnType<typeof TcpSocket.createConnection> | null = null;
	private nextId = 1;
	private pending = new Map<number, ResponseHandler>();
	private buffer = "";
	private _connected = false;

	/** Whether the socket is currently connected. */
	get connected(): boolean {
		return this._connected;
	}

	/**
	 * Open a TCP connection to the MCP gateway.
	 * @param host  Defaults to 127.0.0.1.
	 * @param port  Defaults to 12002.
	 */
	connect(host = "127.0.0.1", port = 12002): Promise<void> {
		return new Promise((resolve, reject) => {
			if (this._connected) {
				resolve();
				return;
			}

			this.socket = TcpSocket.createConnection({ host, port }, () => {
				this._connected = true;
				resolve();
			});

			this.socket.on("data", (data: Buffer) => {
				this.buffer += data.toString();
				this.processBuffer();
			});

			this.socket.on("error", (err: Error) => {
				this._connected = false;
				reject(err);
			});

			this.socket.on("close", () => {
				this._connected = false;
				// Reject all pending requests
				this.pending.forEach((handler, id) => {
					handler({
						jsonrpc: "2.0",
						error: { code: -32000, message: "Connection closed" },
						id,
					});
				});
				this.pending.clear();
			});
		});
	}

	/** Close the TCP socket. */
	disconnect(): void {
		if (this.socket) {
			this.socket.destroy();
			this.socket = null;
		}
		this._connected = false;
		this.pending.clear();
		this.buffer = "";
	}

	/**
	 * Send a JSON-RPC 2.0 request and return the result.
	 * Rejects with the RPC error object if the server returns an error.
	 */
	sendRequest(
		method: string,
		params?: Record<string, unknown>,
	): Promise<unknown> {
		if (!this._connected || !this.socket) {
			return Promise.reject(new Error("Not connected to MCP gateway"));
		}

		const id = this.nextId++;
		const request: RPCRequest = { jsonrpc: "2.0", method, params, id };

		return new Promise((resolve, reject) => {
			this.pending.set(id, (resp: RPCResponse) => {
				if (resp.error) {
					reject(
						new Error(`RPC error ${resp.error.code}: ${resp.error.message}`),
					);
				} else {
					resolve(resp.result);
				}
			});

			try {
				this.socket!.write(JSON.stringify(request) + '\n');
			} catch (err) {
				this.pending.delete(id);
				reject(err);
			}
		});
	}

	/**
	 * Parse complete JSON lines from the receive buffer and dispatch responses.
	 */
	private processBuffer(): void {
		const lines = this.buffer.split("\n");
		// Keep the (possibly partial) last line in the buffer
		this.buffer = lines.pop() || "";

		for (const line of lines) {
			const trimmed = line.trim();
			if (!trimmed) {
				continue;
			}

			try {
				const resp: RPCResponse = JSON.parse(trimmed);
				const handler = this.pending.get(resp.id);
				if (handler) {
					handler(resp);
					this.pending.delete(resp.id);
				}
			} catch {
				// Malformed JSON – ignore and continue
			}
		}
	}
}

const gateway = new TcpGateway();
export default gateway;
