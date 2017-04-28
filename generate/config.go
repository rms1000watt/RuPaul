package generate

// Yaml Config
type Config struct {
	Version         string          `yaml:"Version"`
	MainImportPath  string          `yaml:"MainImportPath"`
	CopyrightHolder string          `yaml:"CopyrightHolder"`
	Datas           map[string]Data `yaml:"Data"`
	APIs            map[string]API  `yaml:"APIs"`
	CommandLine     CommandLine     `yaml:"CommandLine"`
	// Connectors      map[string]Connector `yaml:"Connector"`
}

type Data struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
	// Validations
	Required      bool   `yaml:"Required"`
	MaxLength     int    `yaml:"MaxLength"`
	MinLength     int    `yaml:"MinLength"`
	MustHaveChars string `yaml:"MustHaveChars"`
	CantHaveChars string `yaml:"CantHaveChars"`
	OnlyHaveChars string `yaml:"OnlyHaveChars"`
	GreaterThan   *int   `yaml:"GreaterThan"`
	LessThan      *int   `yaml:"LessThan"`
	// Transforms
	TrimChars    string `yaml:"TrimChars"`
	TrimSpace    bool   `yaml:"TrimSpace"`
	Truncate     int    `yaml:"Truncate"`
	Encrypt      bool   `yaml:"Encrypt"`
	Decrypt      bool   `yaml:"Decrypt"`
	PasswordHash bool   `yaml:"PasswordHash"`
	Hash         bool   `yaml:"Hash"`
	Default      string `yaml:"Default"`
}

type API struct {
	Name          string   `yaml:"Name"`
	Type          string   `yaml:"Type"`
	Serialization string   `yaml:"Serialization"`
	Middlewares   []string `yaml:"Midlewares"`
	Paths         []Path   `yaml:"Paths"`
}

type Path struct {
	Name      string   `yaml:"Name"`
	Pattern   string   `yaml:"Pattern"`
	Inputs    []string `yaml:"Inputs"`
	Connector string   `yaml:"Connector"`
	Methods   []string `yaml:"Methods"`
	Outputs   []string `yaml:"Outputs"`
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

// Template/Generator Config
type TemplateConfig struct {
	Version         string
	MainImportPath  string
	CopyrightHolder string
	API             TemplateAPI
	CommandLine     TemplateCommandLine
}

type TemplateCommandLine struct {
	AppName             string
	AppShortDescription string
	AppLongDescription  string
	GlobalArgs          map[string]Arg
	Command             Command
}

type TemplateAPI struct {
	Name          string
	Type          string
	Serialization string
	Middlewares   []string
	Paths         []TemplatePath
}

type TemplatePath struct {
	Name      string
	Pattern   string
	Inputs    []Data
	Connector string
	Methods   []string
	Outputs   []Data
}
