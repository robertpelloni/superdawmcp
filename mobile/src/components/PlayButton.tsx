import type React from "react";
import {
	TouchableOpacity,
	Text,
	StyleSheet,
	type ViewStyle,
	type TextStyle,
} from "react-native";

interface PlayButtonProps {
	playing?: boolean;
	label?: string;
	onPress: () => void;
	disabled?: boolean;
	color?: string;
}

export default function PlayButton({
	playing = false,
	label,
	onPress,
	disabled = false,
	color,
}: PlayButtonProps): React.JSX.Element {
	const displayLabel = label ?? (playing ? "Stop" : "Play");
	const bgColor = color ?? (playing ? "#e53935" : "#43a047");

	return (
		<TouchableOpacity
			style={[
				styles.button,
				{ backgroundColor: disabled ? "#555" : bgColor },
				playing && styles.playingButton,
			]}
			onPress={onPress}
			disabled={disabled}
			activeOpacity={0.7}
		>
			<Text style={styles.label}>{displayLabel}</Text>
		</TouchableOpacity>
	);
}

const styles = StyleSheet.create({
	button: {
		paddingVertical: 16,
		paddingHorizontal: 48,
		borderRadius: 12,
		alignItems: "center",
		justifyContent: "center",
		minWidth: 160,
	} as ViewStyle,
	playingButton: {
		shadowColor: "#e53935",
		shadowOffset: { width: 0, height: 4 },
		shadowOpacity: 0.4,
		shadowRadius: 8,
		elevation: 6,
	} as ViewStyle,
	label: {
		color: "#fff",
		fontSize: 20,
		fontWeight: "700",
		letterSpacing: 1,
	} as TextStyle,
});
