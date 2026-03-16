package ccsp

type CcspComponentCfg struct {
	ComponentId      string
	ComponentName    string
	Version          uint32
	DbusPath         string
	DmXmlCfgFileName string
}

type CCSPComponent struct {
	cfg CcspComponentCfg
}

func newCCSPComponent(cfg CcspComponentCfg) *CCSPComponent {
	return &CCSPComponent{
		cfg: cfg,
	}
}

func (c *CCSPComponent) Initialize() error {
	return nil
}
