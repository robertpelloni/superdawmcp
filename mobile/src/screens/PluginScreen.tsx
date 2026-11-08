import type React from "react";
import { useState, useCallback } from "react";
import { View, Text, StyleSheet, ScrollView, TouchableOpacity } from "react-native";
import Slider from "@react-native-community/slider";
import { useMcp } from "../hooks/useMcp";

interface PluginParam {
	name: string;
	index: number;
	defaultValue: number;
	minValue: number;
	maxValue: number;
	unit: string;
	category: string;
}

const DEFAULT_PARAMS: PluginParam[] = [
	{ name: "Filter Cutoff", index: 0, defaultValue: 0.55, minValue: 0.0, maxValue: 1.0, unit: "Hz", category: "Tone" },
	{ name: "Filter Resonance", index: 1, defaultValue: 0.25, minValue: 0.0, maxValue: 1.0, unit: "%", category: "Tone" },
	{ name: "Drive", index: 2, defaultValue: 0.20, minValue: 0.0, maxValue: 1.0, unit: "dB", category: "Tone" },
	{ name: "Attack", index: 3, defaultValue: 0.05, minValue: 0.0, maxValue: 1.0, unit: "ms", category: "Env" },
	{ name: "Release", index: 4, defaultValue: 0.30, minValue: 0.0, maxValue: 1.0, unit: "ms", category: "Env" },
	{ name: "Reverb Mix", index: 5, defaultValue: 0.35, minValue: 0.0, maxValue: 1.0, unit: "%", category: "FX" },
];

export default function PluginScreen(): React.JSX.Element {
	const { connected, setPluginParam } = useMcp();
	const [activePlugin, setActivePlugin] = useState("Surge XT");
	const [paramValues, setParamValues] = useState<Record<string, number>>({
		"Filter Cutoff": 0.55,
		"Filter Resonance": 0.25,
		Drive: 0.20,
		Attack: 0.05,
		Release: 0.30,
		"Reverb Mix": 0.35,
	});

	const handleParamChange = useCallback(
		(name: string, val: number) => {
			setParamValues((prev) => ({ ...prev, [name]: val }));
			if (setPluginParam) {
				setPluginParam(activePlugin, name, val);
			}
		},
		[activePlugin, setPluginParam],
	);

	const triggerMacro = useCallback(
		(type: string) => {
			if (type === "brighten") {
				handleParamChange("Filter Cutoff", 0.90);
			} else if (type === "warmth") {
				handleParamChange("Drive", 0.70);
			} else if (type === "space") {
				handleParamChange("Reverb Mix", 0.75);
			} else if (type === "punch") {
				handleParamChange("Attack", 0.01);
				handleParamChange("Release", 0.40);
			}
		},
		[handleParamChange],
	);

	return (
		<ScrollView style={styles.container}>
			<Text style={styles.title}>Plugin Touch Remote</Text>

			<View style={styles.section}>
				<Text style={styles.sectionTitle}>Active VST3 Plugin</Text>
				<View style={styles.pluginPillRow}>
					{["Surge XT", "Vital", "Pro-Q 3", "SuperDAW Reverb"].map((p) => (
						<TouchableOpacity
							key={p}
							style={[styles.pill, activePlugin === p && styles.pillActive]}
							onPress={() => setActivePlugin(p)}
						>
							<Text style={[styles.pillText, activePlugin === p && styles.pillTextActive]}>
								{p}
							</Text>
						</TouchableOpacity>
					))}
				</View>
			</View>

			<View style={styles.section}>
				<Text style={styles.sectionTitle}>Sound Design Macros</Text>
				<View style={styles.macroRow}>
					<TouchableOpacity style={styles.macroBtn} onPress={() => triggerMacro("brighten")}>
						<Text style={styles.macroBtnText}>✨ Brighten</Text>
					</TouchableOpacity>
					<TouchableOpacity style={styles.macroBtn} onPress={() => triggerMacro("warmth")}>
						<Text style={styles.macroBtnText}>🔥 Warmth</Text>
					</TouchableOpacity>
					<TouchableOpacity style={styles.macroBtn} onPress={() => triggerMacro("space")}>
						<Text style={styles.macroBtnText}>🌌 Space</Text>
					</TouchableOpacity>
					<TouchableOpacity style={styles.macroBtn} onPress={() => triggerMacro("punch")}>
						<Text style={styles.macroBtnText}>🥊 Punch</Text>
					</TouchableOpacity>
				</View>
			</View>

			<View style={styles.section}>
				<Text style={styles.sectionTitle}>Controllable Parameters</Text>
				{DEFAULT_PARAMS.map((param) => {
					const val = paramValues[param.name] ?? param.defaultValue;
					return (
						<View key={param.name} style={styles.paramBox}>
							<View style={styles.paramHeader}>
								<Text style={styles.paramName}>{param.name}</Text>
								<View style={styles.nudgeRow}>
									<TouchableOpacity
										style={styles.nudgeBtn}
										onPress={() => handleParamChange(param.name, Math.max(param.minValue, val - 0.05))}
									>
										<Text style={styles.nudgeBtnText}>-0.05</Text>
									</TouchableOpacity>
									<TouchableOpacity
										style={styles.nudgeBtn}
										onPress={() => handleParamChange(param.name, Math.min(param.maxValue, val + 0.05))}
									>
										<Text style={styles.nudgeBtnText}>+0.05</Text>
									</TouchableOpacity>
									<Text style={styles.paramVal}>
										{val.toFixed(2)} {param.unit}
									</Text>
								</View>
							</View>
							<Slider
								style={styles.slider}
								minimumValue={param.minValue}
								maximumValue={param.maxValue}
								step={0.01}
								value={val}
								onValueChange={(v) => handleParamChange(param.name, v)}
								disabled={!connected}
								minimumTrackTintColor="#38bdf8"
								maximumTrackTintColor="#334155"
								thumbTintColor="#38bdf8"
							/>
						</View>
					);
				})}
			</View>
		</ScrollView>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		padding: 16,
		backgroundColor: "#0f172a",
	},
	title: {
		fontSize: 24,
		fontWeight: "800",
		color: "#38bdf8",
		textAlign: "center",
		marginBottom: 20,
		marginTop: 8,
	},
	section: {
		backgroundColor: "#1e293b",
		borderRadius: 12,
		padding: 16,
		marginBottom: 16,
		borderWidth: 1,
		borderColor: "#334155",
	},
	sectionTitle: {
		fontSize: 13,
		fontWeight: "700",
		color: "#94a3b8",
		marginBottom: 12,
		textTransform: "uppercase",
		letterSpacing: 1,
	},
	pluginPillRow: {
		flexDirection: "row",
		flexWrap: "wrap",
		gap: 8,
	},
	pill: {
		paddingHorizontal: 12,
		paddingVertical: 8,
		borderRadius: 20,
		backgroundColor: "#0f172a",
		borderWidth: 1,
		borderColor: "#475569",
	},
	pillActive: {
		backgroundColor: "#38bdf8",
		borderColor: "#38bdf8",
	},
	pillText: {
		color: "#cbd5e1",
		fontSize: 12,
		fontWeight: "600",
	},
	pillTextActive: {
		color: "#0f172a",
		fontWeight: "700",
	},
	macroRow: {
		flexDirection: "row",
		gap: 8,
	},
	macroBtn: {
		flex: 1,
		paddingVertical: 10,
		backgroundColor: "#0f172a",
		borderRadius: 8,
		borderWidth: 1,
		borderColor: "#475569",
		alignItems: "center",
	},
	macroBtnText: {
		color: "#f8fafc",
		fontSize: 11,
		fontWeight: "700",
	},
	paramBox: {
		marginBottom: 14,
	},
	paramHeader: {
		flexDirection: "row",
		justifyContent: "space-between",
		marginBottom: 4,
	},
	paramName: {
		color: "#f8fafc",
		fontSize: 13,
		fontWeight: "600",
	},
	paramVal: {
		color: "#38bdf8",
		fontSize: 12,
		fontWeight: "700",
	},
	slider: {
		width: "100%",
		height: 36,
	},
	nudgeRow: {
		flexDirection: "row",
		alignItems: "center",
		gap: 6,
	},
	nudgeBtn: {
		backgroundColor: "#0f172a",
		paddingHorizontal: 6,
		paddingVertical: 2,
		borderRadius: 4,
		borderWidth: 1,
		borderColor: "#475569",
	},
	nudgeBtnText: {
		color: "#cbd5e1",
		fontSize: 10,
		fontWeight: "600",
	},
});
