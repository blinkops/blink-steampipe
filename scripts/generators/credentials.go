package generators

var Generators = []CredentialGenerator{
	AWSCredentialGenerator{},
	GCPCredentialGenerator{},
	KubernetesCredentialGenerator{},
	StorageCredentialsGenerator{},
}

type CredentialGenerator interface {
	Generate() error
}
