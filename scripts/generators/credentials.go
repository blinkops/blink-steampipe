package generators

var Generators = []CredentialGenerator{
	AWSCredentialGenerator{},
	GCPCredentialGenerator{},
	KubernetesCredentialGenerator{},
	StorageCredentialsGenerator{},
	FreePortGenerator{},
	OCICredentialGenerator{},
}

type CredentialGenerator interface {
	Generate() error
}
