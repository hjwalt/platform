package elasticsearch_integration

import (
	"os"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/hjwalt/platform/type/optional"
)

type Configuration struct {
	Username string
	Password string
	Address  string
	CertFile optional.Optional[string]
}

type Client = elasticsearch.TypedClient

func Create(config Configuration) (*Client, error) {
	opts := make([]elasticsearch.Option, 0)

	if config.CertFile != nil && config.CertFile.IsPresent() {
		cert, certErr := os.ReadFile(config.CertFile.Get())
		if certErr != nil {
			return nil, certErr
		}
		opts = append(opts, elasticsearch.WithCACert(cert))
	}

	opts = append(opts, elasticsearch.WithAddresses(config.Address))
	opts = append(opts, elasticsearch.WithBasicAuth(config.Username, config.Password))

	client, err := elasticsearch.NewTyped(opts...)

	return client, err
}
