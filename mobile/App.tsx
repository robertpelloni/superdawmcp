import type React from "react";
import { NavigationContainer } from "@react-navigation/native";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import HomeScreen from "./src/screens/HomeScreen";

type RootStackParamList = {
	Home: undefined;
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
			</Stack.Navigator>
		</NavigationContainer>
	);
}
