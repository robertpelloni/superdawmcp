import type React from "react";
import { View, Text, StyleSheet } from "react-native";
import Slider from "@react-native-community/slider";

interface TempoSliderProps {
	value: number;
	onValueChange: (value: number) => void;
	disabled?: boolean;
	minimumValue?: number;
	maximumValue?: number;
	step?: number;
}

export default function TempoSlider({
	value,
	onValueChange,
	disabled = false,
	minimumValue = 80,
	maximumValue = 200,
	step = 1,
}: TempoSliderProps): React.JSX.Element {
	return (
		<View style={styles.container}>
			<Text style={styles.bpmLabel}>{Math.round(value)} BPM</Text>
			<Slider
				style={styles.slider}
				minimumValue={minimumValue}
				maximumValue={maximumValue}
				step={step}
				value={value}
				onValueChange={onValueChange}
				disabled={disabled}
				minimumTrackTintColor="#4fc3f7"
				maximumTrackTintColor="#37474f"
				thumbTintColor="#81d4fa"
			/>
			<View style={styles.rangeRow}>
				<Text style={styles.rangeText}>{minimumValue}</Text>
				<Text style={styles.rangeText}>{maximumValue}</Text>
			</View>
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		marginVertical: 8,
	},
	bpmLabel: {
		color: "#e0e0e0",
		fontSize: 32,
		fontWeight: "300",
		textAlign: "center",
		marginBottom: 12,
	},
	slider: {
		width: "100%",
		height: 40,
	},
	rangeRow: {
		flexDirection: "row",
		justifyContent: "space-between",
		marginTop: 4,
	},
	rangeText: {
		color: "#607d8b",
		fontSize: 12,
	},
});
