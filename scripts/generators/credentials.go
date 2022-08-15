package generators

var Generators = []CredentialGenerator{
	AWSCredentialGenerator{},
	GCPCredentialGenerator{},
	KubernetesCredentialGenerator{},
}

type CredentialGenerator interface {
	Generate() error
}
