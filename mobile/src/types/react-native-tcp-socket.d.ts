declare module "@react-native-community/slider" {
	import { Component } from "react";
	import { ViewStyle } from "react-native";

	interface SliderProps {
		style?: ViewStyle;
		value?: number;
		minimumValue?: number;
		maximumValue?: number;
		step?: number;
		disabled?: boolean;
		minimumTrackTintColor?: string;
		maximumTrackTintColor?: string;
		thumbTintColor?: string;
		onValueChange?: (value: number) => void;
	}

	export default class Slider extends Component<SliderProps> {}
}

declare module "react-native-tcp-socket" {
	import { EventEmitter } from "events";

	interface TcpSocketOptions {
		host: string;
		port: number;
	}

	interface TcpSocketInstance extends EventEmitter {
		write(data: string | Buffer): boolean;
		destroy(): void;
		on(event: "data", listener: (data: Buffer) => void): this;
		on(event: "error", listener: (err: Error) => void): this;
		on(event: "close", listener: () => void): this;
		on(event: string, listener: (...args: unknown[]) => void): this;
	}

	export function createConnection(
		options: TcpSocketOptions,
		callback?: () => void,
	): TcpSocketInstance;
}
