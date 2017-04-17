package generate

type Config struct {
	Version string `yaml:"Version"`
	Datas   []Data `yaml:"Data"`
	// DALs        []DAL       `yaml:"DAL"`
	// APIs        []API       `yaml:"API"`
	// CommandLine CommandLine `yaml:"CommandLine"`
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
