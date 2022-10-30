package generators

var Generators = []CredentialGenerator{
	AWSCredentialGenerator{},
	GCPCredentialGenerator{},
	KubernetesCredentialGenerator{},
	StorageCredentialsGenerator{},
	FreePortGenerator{},
}

type CredentialGenerator interface {
	Generate() error
}
