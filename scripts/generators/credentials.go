package generators

var Generators = []CredentialGenerator{
	AWSCredentialGenerator{},
	GCPCredentialGenerator{},
	KubernetesCredentialGenerator{},
	StorageCredentialsGenerator{},
	FreePortGenerator{},
	OCICredentialGenerator{},
	CrowdstrikeCredentialGenerator{},
}

type CredentialGenerator interface {
	Generate() error
}
