package main

import (
	"context"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v9/typedapi/esdsl"
	elasticsearch_integration "github.com/hjwalt/platform/integration/elasticsearch"
	"github.com/hjwalt/platform/type/optional"
)

type Document struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {

	client, err := elasticsearch_integration.Create(elasticsearch_integration.Configuration{
		Username: "elastic",
		Password: "Elastic123",
		Address:  "http://localhost:9200",
		CertFile: optional.Of("./tmp/http_ca.crt"),
	})

	if err != nil {
		panic(err)
	}

	slog.Info("connected")

	index, err := client.Indices.Create("test-index").
		Mappings(
			esdsl.NewTypeMapping().
				AddProperty("Id", esdsl.NewIntegerNumberProperty()).
				AddProperty("Name", esdsl.NewTextProperty()).
				AddProperty("price", esdsl.NewIntegerNumberProperty()),
		).
		Settings(
			esdsl.NewIndexSettings().
				Mode("lookup"),
		).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	slog.Info("index", "res", index)

	toAdd := Document{
		Id:    1,
		Name:  "foo",
		Price: 10,
	}

	added, err := client.Index("test-index").
		Id("1").
		Document(toAdd).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	slog.Info("added", "res", added)

	search, err := client.Search().
		Index("test-index").
		Query(esdsl.NewMatchQuery("name", "foo")).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	slog.Info("searched", "res", search)

	removed, err := client.Delete("test-index", "1").
		Do(context.Background())

	slog.Info("removed", "res", removed)

	del, err := client.Indices.Delete("test-index").
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	slog.Info("del", "res", del)
}
