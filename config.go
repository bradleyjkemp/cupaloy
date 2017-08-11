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

type config struct {
	envVariable       string
	subDirName        string
	snapshotExtension string
}

func defaultConfig() *config {
	return &config{
		envVariable:       "UPDATE_SNAPSHOTS",
		subDirName:        ".snapshots",
		snapshotExtension: "",
	}
}

func (c *config) shouldUpdate() bool {
	return envVariableSet(c.envVariable)
}
