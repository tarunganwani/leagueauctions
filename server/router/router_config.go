package router

//Config - router configuration
type Config struct{
	HostAddress string
	PortNo int
	Secure bool
	CertFilePath string
	KeyPath string
}