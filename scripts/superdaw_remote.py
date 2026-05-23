#!/usr/bin/env python3
import tkinter as tk
from tkinter import ttk, messagebox
import sys
import os
import threading

# Add client to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "pkg", "client", "py")))
from superdaw_client.client import SuperDAWClient

class SuperDAWRemote(tk.Tk):
    def __init__(self):
        super().__init__()
        self.title("SuperDAW Universal Remote")
        self.geometry("600x650")
        self.configure(bg="#1a1a1a")

        self.client = None
        self.server_path = tk.StringVar(value="./bin/superdaw-mcp")
        self.active_daw = tk.StringVar(value="ableton")

        self._build_ui()

    def _build_ui(self):
        style = ttk.Style()
        style.theme_use('clam')
        style.configure("TFrame", background="#1a1a1a")
        style.configure("TLabel", background="#1a1a1a", foreground="#eee")
        style.configure("TButton", padding=5)

        # Connection Frame
        conn_frame = ttk.LabelFrame(self, text="Connection", padding=10)
        conn_frame.pack(fill="x", padx=10, pady=5)
        ttk.Label(conn_frame, text="Server Path:").pack(side="left")
        ttk.Entry(conn_frame, textvariable=self.server_path, width=40).pack(side="left", padx=5)
        ttk.Button(conn_frame, text="Connect", command=self._connect).pack(side="left")

        # DAW Selection
        daw_frame = ttk.Frame(self, padding=10)
        daw_frame.pack(fill="x", padx=10)
        ttk.Label(daw_frame, text="Active DAW:").pack(side="left")
        daws = ["ableton", "reaper", "bitwig", "flstudio", "ardour"]
        ttk.Combobox(daw_frame, textvariable=self.active_daw, values=daws).pack(side="left", padx=5)

        # Transport Frame
        tp_frame = ttk.LabelFrame(self, text="Transport", padding=10)
        tp_frame.pack(fill="x", padx=10, pady=5)
        ttk.Button(tp_frame, text="PLAY", command=lambda: self._transport(True)).pack(side="left", padx=5)
        ttk.Button(tp_frame, text="STOP", command=lambda: self._transport(False)).pack(side="left", padx=5)

        ttk.Label(tp_frame, text="BPM:").pack(side="left", padx=(20, 5))
        self.bpm_val = tk.DoubleVar(value=120.0)
        ttk.Scale(tp_frame, from_=60, to=200, variable=self.bpm_val, orient="horizontal", command=self._on_bpm_change).pack(side="left", fill="x", expand=True)
        ttk.Label(tp_frame, textvariable=self.bpm_val).pack(side="left")

        # Mixer Frame
        mx_frame = ttk.LabelFrame(self, text="Mixer Control", padding=10)
        mx_frame.pack(fill="x", padx=10, pady=5)

        ttk.Label(mx_frame, text="Track ID:").grid(row=0, column=0, sticky="w")
        self.track_id = tk.StringVar(value="0")
        ttk.Entry(mx_frame, textvariable=self.track_id, width=5).grid(row=0, column=1, sticky="w", padx=5)

        ttk.Label(mx_frame, text="Volume:").grid(row=1, column=0, pady=10)
        self.vol_val = tk.DoubleVar(value=0.8)
        ttk.Scale(mx_frame, from_=0, to=1, variable=self.vol_val, orient="horizontal", command=self._on_mixer_change).grid(row=1, column=1, sticky="ew", padx=5)

        ttk.Label(mx_frame, text="Pan:").grid(row=2, column=0)
        self.pan_val = tk.DoubleVar(value=0.0)
        ttk.Scale(mx_frame, from_=-1, to=1, variable=self.pan_val, orient="horizontal", command=self._on_mixer_change).grid(row=2, column=1, sticky="ew", padx=5)

        # Routing Frame
        rt_frame = ttk.LabelFrame(self, text="Audio Routing Matrix", padding=10)
        rt_frame.pack(fill="x", padx=10, pady=5)

        ttk.Label(rt_frame, text="Source DAW:").grid(row=0, column=0)
        self.src_daw = tk.StringVar(value="ableton")
        ttk.Entry(rt_frame, textvariable=self.src_daw, width=10).grid(row=0, column=1, padx=5)

        ttk.Label(rt_frame, text="Track:").grid(row=0, column=2)
        self.src_track = tk.StringVar(value="1")
        ttk.Entry(rt_frame, textvariable=self.src_track, width=5).grid(row=0, column=3, padx=5)

        ttk.Label(rt_frame, text=" -> Dest DAW:").grid(row=1, column=0, pady=5)
        self.dst_daw = tk.StringVar(value="reaper")
        ttk.Entry(rt_frame, textvariable=self.dst_daw, width=10).grid(row=1, column=1, padx=5)

        ttk.Label(rt_frame, text="Track:").grid(row=1, column=2)
        self.dst_track = tk.StringVar(value="MASTER")
        ttk.Entry(rt_frame, textvariable=self.dst_track, width=5).grid(row=1, column=3, padx=5)

        ttk.Button(rt_frame, text="PATCH", command=self._patch).grid(row=2, column=0, columnspan=4, pady=10)

        # Console
        self.log = tk.Text(self, height=8, bg="#000", fg="#00ff00", font=("Consolas", 9))
        self.log.pack(fill="both", expand=True, padx=10, pady=10)

    def _connect(self):
        try:
            self.client = SuperDAWClient(self.server_path.get())
            self.client.connect()
            self._write_log("Connected to SuperDAW-MCP server.")
        except Exception as e:
            messagebox.showerror("Error", str(e))

    def _write_log(self, msg):
        self.log.insert("end", f"> {msg}\n")
        self.log.see("end")

    def _transport(self, playing):
        if not self.client: return
        self.client.transport_control(playing=playing, bpm=self.bpm_val.get(), daw=self.active_daw.get())
        self._write_log(f"Transport: {'PLAY' if playing else 'STOP'} on {self.active_daw.get()}")

    def _on_bpm_change(self, val):
        if not self.client: return
        # Throttle? Or just send
        self.client.transport_control(playing=True, bpm=float(val), daw=self.active_daw.get())

    def _on_mixer_change(self, _):
        if not self.client: return
        self.client.set_mixer(self.track_id.get(), self.vol_val.get(), self.pan_val.get(), daw=self.active_daw.get())

    def _patch(self):
        if not self.client: return
        # Using custom command for patching or specialized tool if available
        args = {
            "source_daw": self.src_daw.get(),
            "source_track": self.src_track.get(),
            "dest_daw": self.dst_daw.get(),
            "dest_track": self.dst_track.get()
        }
        self.client._call("tools/call", {"name": "superdaw_patch_audio", "arguments": args})
        self._write_log(f"Patched {self.src_daw.get()}:{self.src_track.get()} -> {self.dst_daw.get()}:{self.dst_track.get()}")

if __name__ == "__main__":
    app = SuperDAWRemote()
    app.mainloop()
