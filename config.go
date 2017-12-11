package cupaloy

// Configurator is a functional option that can be passed to cupaloy.New() to change snapshotting behaviour.
type Configurator func(*config)

// EnvVariableName can be used to customize the environment variable that determines whether snapshots should be updated.
// e.g.
//  cupaloy.New(EnvVariableName("UPDATE"))
// Will create an instance where snapshots will be updated if the UPDATE environment variable is set,
// instead of the default of UPDATE_SNAPSHOTS.
func EnvVariableName(name string) Configurator {
	return func(c *config) {
		c.envVariable = name
	}
}

// SnapshotSubdirectory can be used to customize the location that snapshots are stored in.
// e.g.
//  cupaloy.New(SnapshotSubdirectory("testdata"))
// Will create an instance where snapshots are stored in testdata/ rather than the default .snapshots/
func SnapshotSubdirectory(name string) Configurator {
	return func(c *config) {
		c.subDirName = name
	}
}

// RawOutput can be used to enable dumping raw output to the snapshot file
// rather than using go-spew. This raw output is possible when snapshotting
// strings or io.Readers
func RawOutput(b bool) Configurator {
	return func(c *config) {
		c.rawOutput = b
	}
}

type config struct {
	envVariable       string
	subDirName        string
	snapshotExtension string
	rawOutput         bool
}

func defaultConfig() *config {
	return &config{
		envVariable:       "UPDATE_SNAPSHOTS",
		subDirName:        ".snapshots",
		snapshotExtension: "",
		rawOutput:         false,
	}
}

func (c *config) shouldUpdate() bool {
	return envVariableSet(c.envVariable)
}
