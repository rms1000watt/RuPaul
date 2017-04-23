package generate

type Config struct {
	Version         string `yaml:"Version"`
	MainImportPath  string `yaml:"MainImportPath"`
	CopyrightHolder string `yaml:"CopyrightHolder"`
	// Datas          map[string]Data `yaml:"Data"`
	// Connectors     map[string]Connector `yaml:"Connector"`
	APIs        map[string]API `yaml:"API"`
	CommandLine CommandLine    `yaml:"CommandLine"`
}

type Data struct {
	Name          string `yaml:"Name"`
	Type          string `yaml:"Type"`
	Default       string `yaml:"Default"`
	Required      bool   `yaml:"Required"`
	Encrypted     bool   `yaml:"Encrypted"`
	Hashed        string `yaml:"Hashed"`
	MaxLength     int    `yaml:"MaxLength"`
	MinLength     int    `yaml:"MinLength"`
	MustHaveChars string `yaml:"MustHaveChars"`
	CantHaveChars string `yaml:"CantHaveChars"`
	OnlyHaveChars string `yaml:"OnlyHaveChars"`
}

type API struct {
	Name string `yaml:"Name"`
}

type CommandLine struct {
	AppName             string             `yaml:"AppName"`
	AppShortDescription string             `yaml:"AppShortDescription"`
	AppLongDescription  string             `yaml:"AppLongDescription"`
	GlobalArgs          map[string]Arg     `yaml:"GlobalArgs"`
	Commands            map[string]Command `yaml:"Commands"`
}

type Arg struct {
	Name        string `yaml:"Name"`
	Description string `yaml:"Description"`
	ShortName   string `yaml:"ShortName"`
	Type        string `yaml:"Type"`
	Default     string `yaml:"Default"`
}

type Command struct {
	Name             string         `yaml:"Name"`
	ShortDescription string         `yaml:"ShortDescription"`
	LongDescription  string         `yaml:"LongDescription"`
	Args             map[string]Arg `yaml:"Args"`
	API              string         `yaml:"API"`
}

type SingleConfig struct {
	Version         string
	MainImportPath  string
	CopyrightHolder string
	APIs            map[string]API
	CommandLine     SingleCommandLine
}

type SingleCommandLine struct {
	AppName             string
	AppShortDescription string
	AppLongDescription  string
	GlobalArgs          map[string]Arg
	Command             Command
}
