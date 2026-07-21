import React, { useState, useEffect } from "react";
import { View, Text, StyleSheet, ScrollView } from "react-native";
import Slider from "@react-native-community/slider";
import { useMcp } from "../hooks/useMcp";
import TcpGateway from "../services/TcpGateway";

interface Track {
    id: string;
    name: string;
    volume: number;
}

export default function MixerScreen(): React.JSX.Element {
	const { connected } = useMcp();
	const [tracks, setTracks] = useState<Track[]>([]);

    useEffect(() => {
        if (!connected) return;
        TcpGateway.sendRequest("superdaw_get_tracks", {})
            .then((res: any) => {
                if (Array.isArray(res)) {
                    setTracks(res.map(t => ({id: t.id || t.Name, name: t.name || t.Name, volume: 0.8})));
                } else if (res.content && res.content[0] && res.content[0].text) {
                    try {
                        const parsed = JSON.parse(res.content[0].text);
                        setTracks(parsed.map((t: any) => ({id: t.id || t.Name, name: t.name || t.Name, volume: 0.8})));
                    } catch(e) {}
                }
            })
            .catch(() => {});
    }, [connected]);

    const handleVolumeChange = (id: string, vol: number) => {
        if (!connected) return;
        TcpGateway.sendRequest("superdaw_set_mixer", { track_id: id, volume: vol });
        setTracks(prev => prev.map(t => t.id === id ? { ...t, volume: vol } : t));
    };

	return (
		<ScrollView style={styles.container}>
			<Text style={styles.title}>Mixer</Text>
			{!connected && <Text style={styles.warning}>Connect via Home first</Text>}
            {tracks.length === 0 && connected && <Text style={styles.warning}>No tracks found or loading...</Text>}
			{tracks.map((t) => (
				<View key={t.id} style={styles.trackRow}>
					<Text style={styles.trackName}>{t.name}</Text>
					<Slider
						style={styles.slider}
						minimumValue={0}
						maximumValue={1}
						value={t.volume}
						disabled={!connected}
						onSlidingComplete={(v) => handleVolumeChange(t.id, v)}
					/>
				</View>
			))}
		</ScrollView>
	);
}

const styles = StyleSheet.create({
	container: { flex: 1, padding: 24, backgroundColor: "#1a1a2e" },
	title: { fontSize: 28, color: "#e0e0e0", textAlign: "center", marginBottom: 32 },
	warning: { color: "#f44336", textAlign: "center", marginBottom: 16 },
	trackRow: { flexDirection: "row", alignItems: "center", marginBottom: 16 },
	trackName: { color: "#fff", width: 80, fontSize: 16 },
	slider: { flex: 1, height: 40 },
});
