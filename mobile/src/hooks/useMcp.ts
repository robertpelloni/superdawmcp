import { useState, useCallback, useEffect } from "react";
import TcpGateway from "../services/TcpGateway";

export interface UseMcpReturn {
	connected: boolean;
	connecting: boolean;
	error: string | null;
	tempo: number;
	playing: boolean;
	connect: () => Promise<void>;
	disconnect: () => void;
	play: (bpm?: number) => Promise<void>;
	stop: () => Promise<void>;
	setTempo: (bpm: number) => Promise<void>;
	togglePlay: () => Promise<void>;
}

export function useMcp(): UseMcpReturn {
	const [connected, setConnected] = useState(false);
	const [connecting, setConnecting] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [tempo, setTempoState] = useState(145);
	const [playing, setPlaying] = useState(false);


	useEffect(() => {
		TcpGateway.onStateUpdate = (state: any) => {
			if (state.playing !== undefined) {
				setPlaying(state.playing);
			}
			if (state.bpm !== undefined) {
				setTempoState(state.bpm);
			}
		};
		return () => {
			TcpGateway.onStateUpdate = undefined;
		};
	}, []);

	const connect = useCallback(async () => {
		setConnecting(true);
		setError(null);
		try {
			await TcpGateway.connect();
			setConnected(true);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : String(err);
			setError(`Connection failed: ${msg}`);
			setConnected(false);
		} finally {
			setConnecting(false);
		}
	}, []);

	const disconnect = useCallback(() => {
		TcpGateway.disconnect();
		setConnected(false);
		setPlaying(false);
	}, []);

	const transport = useCallback(async (play: boolean, bpm: number) => {
		if (!TcpGateway.connected) {
			return;
		}
		await TcpGateway.sendRequest("superdaw_transport_control", {
			playing: play,
			bpm,
		});
	}, []);

	const play = useCallback(
		async (bpm?: number) => {
			const targetBpm = bpm ?? tempo;
			await transport(true, targetBpm);
			setPlaying(true);
		},
		[tempo, transport],
	);

	const stop = useCallback(async () => {
		await transport(false, tempo);
		setPlaying(false);
	}, [tempo, transport]);

	const togglePlay = useCallback(async () => {
		if (playing) {
			await stop();
		} else {
			await play();
		}
	}, [playing, play, stop]);

	const setTempo = useCallback(
		async (bpm: number) => {
			setTempoState(bpm);
			if (playing) {
				await transport(true, bpm);
			}
		},
		[playing, transport],
	);

	return {
		connected,
		connecting,
		error,
		tempo,
		playing,
		connect,
		disconnect,
		play,
		stop,
		setTempo,
		togglePlay,
	};
}
