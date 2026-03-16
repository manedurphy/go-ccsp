package ccsp

type CcspComponentCfg struct {
	ComponentID      string `xml:"ID"`
	ComponentName    string `xml:"Name"`
	Version          uint32 `xml:"Version"`
	DbusPath         string `xml:"DbusPath"`
	DmXMLCfgFileName string `xml:"DataModelXmlCfg"`
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
