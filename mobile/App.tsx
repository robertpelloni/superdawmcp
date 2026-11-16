import type React from "react";
import { NavigationContainer } from "@react-navigation/native";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import HomeScreen from "./src/screens/HomeScreen";
import MixerScreen from "./src/screens/MixerScreen";
import PluginScreen from "./src/screens/PluginScreen";

type RootStackParamList = {
	Home: undefined;
	Mixer: undefined;
	Plugins: undefined;
};

const Stack = createNativeStackNavigator<RootStackParamList>();

export default function App(): React.JSX.Element {
	return (
		<NavigationContainer>
			<Stack.Navigator initialRouteName="Home">
				<Stack.Screen
					name="Home"
					component={HomeScreen}
					options={{ title: "SuperDAW Remote" }}
				/>
				<Stack.Screen
					name="Mixer"
					component={MixerScreen}
					options={{ title: "Mixer" }}
				/>
				<Stack.Screen
					name="Plugins"
					component={PluginScreen}
					options={{ title: "Plugin Inspector" }}
				/>
			</Stack.Navigator>
		</NavigationContainer>
	);
}
