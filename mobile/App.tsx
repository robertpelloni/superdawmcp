import type React from "react";
import { NavigationContainer } from "@react-navigation/native";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import HomeScreen from "./src/screens/HomeScreen";

import PluginScreen from "./src/screens/PluginScreen";

type RootStackParamList = {
	Home: undefined;
	Plugin: undefined;
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
					name="Plugin"
					component={PluginScreen}
					options={{ title: "VST3 Touch Inspector" }}
				/>
			</Stack.Navigator>
		</NavigationContainer>
	);
}
