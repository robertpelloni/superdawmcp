module.exports = {
	preset: "react-native",
	moduleFileExtensions: ["ts", "tsx", "js", "jsx", "json", "node"],
	transformIgnorePatterns: [
		"node_modules/(?!(react-native|@react-native|@react-navigation|react-native-tcp-socket|@react-native-community/slider)/)",
	],
	testMatch: ["**/__tests__/**/*.test.(ts|tsx|js)"],
	collectCoverage: true,
	coverageDirectory: "coverage",
};
