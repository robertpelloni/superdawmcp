use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use std::io::{BufRead, BufReader, Write};
use std::process::{Child, ChildStdin, ChildStdout, Command, Stdio};

pub struct SuperDAWClient {
    child: Child,
    stdin: ChildStdin,
    reader: BufReader<ChildStdout>,
    request_id: i32,
}

impl SuperDAWClient {
    pub fn new(server_path: &str) -> std::io.Result<Self> {
        let mut child = Command::new(server_path)
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .spawn()?;

        let stdin = child.stdin.take().ok_or(std::io::ErrorKind::Other)?;
        let stdout = child.stdout.take().ok_or(std::io::ErrorKind::Other)?;
        let reader = BufReader::new(stdout);

        Ok(Self {
            child,
            stdin,
            reader,
            request_id: 1,
        })
    }

    pub fn call(&mut self, method: &str, params: Value) -> std::io.Result<Value> {
        let request = json!({
            "jsonrpc": "2.0",
            "method": method,
            "params": params,
            "id": self.request_id,
        });
        self.request_id += 1;

        let json = serde_json::to_string(&request)?;
        writeln!(self.stdin, "{}", json)?;
        self.stdin.flush()?;

        let mut line = String::new();
        self.reader.read_line(&mut line)?;
        let response: Value = serde_json::from_str(&line)?;

        Ok(response["result"].clone())
    }

    pub fn set_mixer(&mut self, track_id: &str, volume: f32, pan: f32, daw: Option<&str>) -> std::io.Result<()> {
        let mut args = json!({
            "track_id": track_id,
            "volume": volume,
            "pan": pan,
        });
        if let Some(d) = daw {
            args["daw"] = json!(d);
        }

        self.call("tools/call", json!({ "name": "superdaw_set_mixer", "arguments": args }))?;
        Ok(())
    }

    pub fn transport_control(&mut self, playing: bool, bpm: Option<f64>, daw: Option<&str>) -> std::io.Result<()> {
        let mut args = json!({ "playing": playing });
        if let Some(b) = bpm { args["bpm"] = json!(b); }
        if let Some(d) = daw { args["daw"] = json!(d); }

        self.call("tools/call", json!({ "name": "superdaw_transport_control", "arguments": args }))?;
        Ok(())
    }
}

impl Drop for SuperDAWClient {
    fn drop(&mut self) {
        let _ = self.child.kill();
    }
}
