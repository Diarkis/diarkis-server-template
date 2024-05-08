package matching

/*
Setup sets up MatchMaker on the server.
You must call this at the start of the server process.

	confpath - Absolute path of the configuration file to be loaded.
*/
func Setup(confpath string) {
}

/*
ClearDefinition clears already defined match making definition

	matchingID - Matching profile ID to clear the definition.
*/
func ClearDefinition(matchingID string) {
}

/*
Define defines a match making search schema:

	[IMPORTANT] This must be defined on the server that is specified by "targetNodeType" configuration
	            because all match making data is stored on those servers.
	[IMPORTANT] Profile ID must not contain "\t".

Parameters

	profileID  - Unique matching profile ID.
	props      - Matching profile condition properties.

In order to perform matchmaking, you must define the rules for matchmakings.

The matchmaking rules are called profiles.

These rules will dictate how matchmaking should be conditioned.

You may combine multiple matchmaking rules and create more complex matchmaking conditions as well.

You must define matchmaking rule profiles before invoking diarkis.Start.

The example below shows a matchmaking rule that uses level and creates buckets of matchmaking pools by the range of 10.

With this profile, each level bucket will pool users with level 0 to 10, 11 to 20, 21 to 30 and so forth...

The string name given as LevelMatch is the unique ID to represents the profile.

	levelMatchProfile := make(map[string]int)

	levelMatchProfile["level"] = 10

	matching.Define("LevelMatch", levelMatchProfile)

You may define as many matching definition as you require as long as each profileID is unique.
*/
func Define(profileID string, props map[string]int) {
}

/*
DefineByJSON defines multiple match making definitions from JSON string.

$(...) represents a variable.

	{
		"$(profile ID)": {
			"$(property name)": $(property value as int)
			"$(property name)": $(property value as int)
			"$(property name)": $(property value as int)
		}
	}

	jsonBytes - Matching profile byte array data to be used to define the profile.
*/
func DefineByJSON(jsonBytes []byte) {
}
