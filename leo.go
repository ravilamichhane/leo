package leo

const version = "1.0.0"

type Leo struct {
	AppName string
	Debug   bool
	Version string
}

func (l *Leo) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{
			"handlers", "migrations", "src", "data", "public", "tmp", "logs", "middleware",
		},
	}

	if err := l.Init(pathConfig); err != nil {
		return err
	}

	return nil
}

func (l *Leo) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		if err := l.CreateDirIfNotExist(root + "/" + path); err != nil {
			return err
		}
	}
	return nil
}
