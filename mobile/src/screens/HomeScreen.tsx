import type React from "react";
import { View, Text, StyleSheet } from "react-native";
import { useMcp } from "../hooks/useMcp";
import PlayButton from "../components/PlayButton";
import TempoSlider from "../components/TempoSlider";

export default function HomeScreen(): React.JSX.Element {
	const {
		connected,
		connecting,
		error,
		tempo,
		playing,
		connect,
		disconnect,
		togglePlay,
		setTempo,
	} = useMcp();

	return (
		<View style={styles.container}>
			<Text style={styles.title}>SuperDAW Remote</Text>

			{/* Connection section */}
			<View style={styles.section}>
				<Text style={styles.sectionTitle}>Connection</Text>
				<View style={styles.statusRow}>
					<View
						style={[
							styles.statusDot,
							{ backgroundColor: connected ? "#4caf50" : "#f44336" },
						]}
					/>
					<Text style={styles.statusText}>
						{connecting
							? "Connecting..."
							: connected
								? "Connected"
								: "Disconnected"}
					</Text>
				</View>
				{error ? <Text style={styles.errorText}>{error}</Text> : null}
				{connected ? (
					<Text style={styles.disconnectLink} onPress={disconnect}>
						Disconnect
					</Text>
				) : (
					<PlayButton
						label={connecting ? "Connecting..." : "Connect"}
						onPress={connect}
						disabled={connecting}
						color="#2196f3"
					/>
				)}
			</View>

			{/* Transport section */}
			<View style={styles.section}>
				<Text style={styles.sectionTitle}>Transport</Text>
				<PlayButton
					playing={playing}
					onPress={togglePlay}
					disabled={!connected}
				/>
			</View>

			{/* Tempo section */}
			<View style={styles.section}>
				<Text style={styles.sectionTitle}>Tempo</Text>
				<TempoSlider
					value={tempo}
					onValueChange={setTempo}
					disabled={!connected}
					minimumValue={80}
					maximumValue={200}
					step={1}
				/>
			</View>
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		padding: 24,
		backgroundColor: "#1a1a2e",
	},
	title: {
		fontSize: 28,
		fontWeight: "bold",
		color: "#e0e0e0",
		textAlign: "center",
		marginBottom: 32,
		marginTop: 16,
	},
	section: {
		backgroundColor: "#16213e",
		borderRadius: 12,
		padding: 20,
		marginBottom: 20,
		borderWidth: 1,
		borderColor: "#0f3460",
	},
	sectionTitle: {
		fontSize: 16,
		fontWeight: "600",
		color: "#a0a0b0",
		marginBottom: 12,
		textTransform: "uppercase",
		letterSpacing: 1,
	},
	statusRow: {
		flexDirection: "row",
		alignItems: "center",
		marginBottom: 12,
	},
	statusDot: {
		width: 10,
		height: 10,
		borderRadius: 5,
		marginRight: 8,
	},
	statusText: {
		color: "#c0c0d0",
		fontSize: 15,
	},
	errorText: {
		color: "#f44336",
		fontSize: 13,
		marginBottom: 8,
	},
	disconnectLink: {
		color: "#ff7043",
		fontSize: 14,
		textDecorationLine: "underline",
		marginTop: 8,
	},
});
