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
	isBipolar?: boolean;
	smoothingMs?: number;
	isFavorite?: boolean;
	tags?: string[];
	macroAssignments?: Record<string, number>;
	midiCc?: number;
}

const DEFAULT_PARAMS: PluginParam[] = [
	{ name: "Filter Cutoff", index: 0, defaultValue: 0.55, minValue: 0.0, maxValue: 1.0, unit: "Hz", category: "Tone", smoothingMs: 8.0, isFavorite: true, tags: ["Filter", "Tone"], macroAssignments: { brighten: 1.0 }, midiCc: 74 },
	{ name: "Filter Resonance", index: 1, defaultValue: 0.25, minValue: 0.0, maxValue: 1.0, unit: "%", category: "Tone", smoothingMs: 10.0, tags: ["Filter", "Resonance"], macroAssignments: { brighten: 0.3 }, midiCc: 71 },
	{ name: "Drive", index: 2, defaultValue: 0.20, minValue: 0.0, maxValue: 1.0, unit: "dB", category: "Tone", smoothingMs: 12.0, tags: ["Distortion", "Drive"], macroAssignments: { warmth: 0.85 } },
	{ name: "Attack", index: 3, defaultValue: 0.05, minValue: 0.0, maxValue: 1.0, unit: "ms", category: "Env", smoothingMs: 5.0, tags: ["Envelope", "Dynamics"], macroAssignments: { punch: -0.8 }, midiCc: 73 },
	{ name: "Release", index: 4, defaultValue: 0.30, minValue: 0.0, maxValue: 1.0, unit: "ms", category: "Env", smoothingMs: 5.0, tags: ["Envelope", "Dynamics"], macroAssignments: { punch: 0.6 }, midiCc: 72 },
	{ name: "Reverb Mix", index: 5, defaultValue: 0.35, minValue: 0.0, maxValue: 1.0, unit: "%", category: "FX", smoothingMs: 15.0, isFavorite: true, tags: ["Space", "FX"], macroAssignments: { space: 1.0 }, midiCc: 91 },
];

export default function PluginScreen(): React.JSX.Element {
	const { connected, setPluginParam } = useMcp();
	const [activePlugin, setActivePlugin] = useState("Surge XT");
	const [favorites, setFavorites] = useState<Record<string, boolean>>({
		"Filter Cutoff": true,
		"Reverb Mix": true,
	});
	const [onlyFavorites, setOnlyFavorites] = useState(false);
	const [paramValues, setParamValues] = useState<Record<string, number>>({
		"Filter Cutoff": 0.55,
		"Filter Resonance": 0.25,
		Drive: 0.20,
		Attack: 0.05,
		Release: 0.30,
		"Reverb Mix": 0.35,
	});
	const [snapshotBanks, setSnapshotBanks] = useState<Record<string, Record<string, number>>>({
		A: { "Filter Cutoff": 0.55, "Filter Resonance": 0.25, Drive: 0.20, Attack: 0.05, Release: 0.30, "Reverb Mix": 0.35 },
		B: { "Filter Cutoff": 0.85, "Filter Resonance": 0.60, Drive: 0.50, Attack: 0.01, Release: 0.60, "Reverb Mix": 0.70 },
		C: {},
		D: {},
	});
	const [morphVal, setMorphVal] = useState(0.0);

	const handleParamChange = useCallback(
		(name: string, val: number) => {
			setParamValues((prev) => ({ ...prev, [name]: val }));
			if (setPluginParam) {
				setPluginParam(activePlugin, name, val);
			}
		},
		[activePlugin, setPluginParam],
	);

	const handleResetParam = useCallback(
		(name: string, defVal: number) => {
			handleParamChange(name, defVal);
		},
		[handleParamChange],
	);

	const toggleFavorite = useCallback(
		(name: string) => {
			setFavorites((prev) => ({ ...prev, [name]: !prev[name] }));
		},
		[],
	);

	const saveSnapshot = useCallback(
		(slot: string) => {
			setSnapshotBanks((prev) => ({ ...prev, [slot]: { ...paramValues } }));
		},
		[paramValues],
	);

	const handleMorph = useCallback(
		(alpha: number) => {
			setMorphVal(alpha);
			const bankA = snapshotBanks.A;
			const bankB = snapshotBanks.B;
			if (!bankA || !bankB || Object.keys(bankA).length === 0 || Object.keys(bankB).length === 0) return;

			DEFAULT_PARAMS.forEach((param) => {
				const a = bankA[param.name];
				const b = bankB[param.name];
				if (a !== undefined && b !== undefined) {
					const morphed = a + (b - a) * alpha;
					handleParamChange(param.name, morphed);
				}
			});
		},
		[snapshotBanks, handleParamChange],
	);

	const triggerMacro = useCallback(
		(type: string) => {
			DEFAULT_PARAMS.forEach((param) => {
				const curVal = paramValues[param.name] ?? param.defaultValue;
				if (param.macroAssignments && param.macroAssignments[type] !== undefined) {
					const weight = param.macroAssignments[type];
					const delta = (param.maxValue - param.minValue) * 0.35 * weight;
					const targetVal = Math.min(param.maxValue, Math.max(param.minValue, curVal + delta));
					handleParamChange(param.name, targetVal);
				} else {
					if (type === "brighten" && param.name.includes("Cutoff")) {
						handleParamChange(param.name, 0.90);
					} else if (type === "warmth" && param.name.includes("Drive")) {
						handleParamChange(param.name, 0.70);
					} else if (type === "space" && param.name.includes("Reverb")) {
						handleParamChange(param.name, 0.75);
					} else if (type === "punch" && param.name.includes("Attack")) {
						handleParamChange(param.name, 0.01);
					}
				}
			});
		},
		[handleParamChange, paramValues],
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
				<View style={{ flexDirection: "row", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
					<Text style={styles.sectionTitle}>Sound Design Macros & History</Text>
					<View style={{ flexDirection: "row", gap: 6 }}>
						<TouchableOpacity
							style={styles.bankBtn}
							onPress={() => {
								// Undo touch simulation
							}}
						>
							<Text style={styles.bankBtnText}>↶ Undo</Text>
						</TouchableOpacity>
						<TouchableOpacity
							style={styles.bankBtn}
							onPress={() => {
								// Redo touch simulation
							}}
						>
							<Text style={styles.bankBtnText}>↷ Redo</Text>
						</TouchableOpacity>
						<TouchableOpacity
							style={styles.bankBtn}
							onPress={() => {
								DEFAULT_PARAMS.forEach((param) => {
									const cur = paramValues[param.name] ?? param.defaultValue;
									const span = param.maxValue - param.minValue;
									const delta = (Math.random() * 2.0 - 1.0) * span * 0.125;
									const mutated = Math.min(param.maxValue, Math.max(param.minValue, cur + delta));
									handleParamChange(param.name, mutated);
								});
							}}
						>
							<Text style={styles.bankBtnText}>🧬 Mutate</Text>
						</TouchableOpacity>
						<TouchableOpacity
							style={styles.bankBtn}
							onPress={() => {
								DEFAULT_PARAMS.forEach((param) => {
									const rnd = param.minValue + Math.random() * (param.maxValue - param.minValue);
									handleParamChange(param.name, rnd);
								});
							}}
						>
							<Text style={styles.bankBtnText}>🎲 Rnd</Text>
						</TouchableOpacity>
					</View>
				</View>
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

				{/* Snapshot Banks & Continuous Morphing */}
				<View style={styles.morphHeader}>
					<Text style={styles.morphLabel}>BANKS: </Text>
					{["A", "B", "C", "D"].map((slot) => (
						<TouchableOpacity
							key={slot}
							style={styles.bankBtn}
							onPress={() => saveSnapshot(slot)}
						>
							<Text style={styles.bankBtnText}>Slot {slot}</Text>
						</TouchableOpacity>
					))}
				</View>
				<View style={styles.morphSliderRow}>
					<Text style={styles.morphLabel}>A ⟷ B Morph ({Math.round(morphVal * 100)}%)</Text>
					<Slider
						style={styles.morphSlider}
						minimumValue={0.0}
						maximumValue={1.0}
						step={0.01}
						value={morphVal}
						onValueChange={handleMorph}
						disabled={!connected}
						minimumTrackTintColor="#10b981"
						maximumTrackTintColor="#334155"
						thumbTintColor="#10b981"
					/>
				</View>
			</View>

			<View style={styles.section}>
				<View style={{ flexDirection: "row", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
					<Text style={styles.sectionTitle}>Controllable Parameters</Text>
					<TouchableOpacity
						style={[styles.bankBtn, onlyFavorites && { backgroundColor: "#f59e0b" }]}
						onPress={() => setOnlyFavorites((prev) => !prev)}
					>
						<Text style={[styles.bankBtnText, onlyFavorites && { color: "#000", fontWeight: "700" }]}>⭐ Favorites</Text>
					</TouchableOpacity>
				</View>
				{DEFAULT_PARAMS.filter((param) => (onlyFavorites ? !!favorites[param.name] : true)).map((param) => {
					const val = paramValues[param.name] ?? param.defaultValue;
					const isFav = !!favorites[param.name];
					return (
						<View key={param.name} style={styles.paramBox}>
							<View style={styles.paramHeader}>
								<View style={{ flexDirection: "row", alignItems: "center", gap: 6 }}>
									<TouchableOpacity onPress={() => toggleFavorite(param.name)}>
										<Text style={{ fontSize: 13 }}>{isFav ? "⭐" : "☆"}</Text>
									</TouchableOpacity>
									{param.midiCc !== undefined && param.midiCc >= 0 && (
										<View style={{ backgroundColor: "#581c87", paddingHorizontal: 4, paddingVertical: 1, borderRadius: 3 }}>
											<Text style={{ color: "#d8b4fe", fontSize: 9, fontWeight: "700" }}>CC {param.midiCc}</Text>
										</View>
									)}
									<Text style={styles.paramName}>{param.name}</Text>
								</View>
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
									<TouchableOpacity
										style={styles.resetBtn}
										onPress={() => handleResetParam(param.name, param.defaultValue)}
									>
										<Text style={styles.resetBtnText}>↺</Text>
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
	resetBtn: {
		backgroundColor: "#0f172a",
		paddingHorizontal: 6,
		paddingVertical: 2,
		borderRadius: 4,
		borderWidth: 1,
		borderColor: "#475569",
	},
	resetBtnText: {
		color: "#38bdf8",
		fontSize: 11,
		fontWeight: "700",
	},
	morphHeader: {
		flexDirection: "row",
		alignItems: "center",
		gap: 6,
		marginTop: 12,
		paddingTop: 8,
		borderTopWidth: 1,
		borderTopColor: "#334155",
	},
	morphLabel: {
		color: "#94a3b8",
		fontSize: 11,
		fontWeight: "700",
	},
	bankBtn: {
		backgroundColor: "#0f172a",
		paddingHorizontal: 8,
		paddingVertical: 4,
		borderRadius: 4,
		borderWidth: 1,
		borderColor: "#475569",
	},
	bankBtnText: {
		color: "#cbd5e1",
		fontSize: 10,
		fontWeight: "600",
	},
	morphSliderRow: {
		marginTop: 8,
	},
	morphSlider: {
		width: "100%",
		height: 30,
	},
});
