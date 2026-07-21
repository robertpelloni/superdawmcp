import React, { useState } from "react";
import { View, Text, StyleSheet, Button, ScrollView } from "react-native";
import { useMcp } from "../hooks/useMcp";
import TcpGateway from "../services/TcpGateway";

export default function PluginScreen(): React.JSX.Element {
	const { connected } = useMcp();
	const [params, setParams] = useState<string>("No plugin data loaded.");

	const fetchParams = async () => {
		if (!connected) return;
		try {
			const res: any = await TcpGateway.sendRequest("superdaw_get_plugin_params", { plugin_name: "Serum" });
			setParams(JSON.stringify(res, null, 2));
		} catch (e) {
			setParams("Error fetching parameters from daemon.");
		}
	};

	return (
		<ScrollView style={styles.container}>
			<Text style={styles.title}>VST3 Inspector</Text>
			<Button title="Introspect Target Plugin" onPress={fetchParams} disabled={!connected} color="#4caf50" />
			<Text style={styles.output}>{params}</Text>
		</ScrollView>
	);
}

const styles = StyleSheet.create({
	container: { flex: 1, padding: 24, backgroundColor: "#1a1a2e" },
	title: { fontSize: 28, color: "#e0e0e0", textAlign: "center", marginBottom: 32 },
	output: { color: "#81c784", marginTop: 20, fontFamily: "monospace" },
});
