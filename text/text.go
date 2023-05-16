package text

import (
	"context"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
)

var pineconeEnv = "us-central1-gcp"
var indexName = "database"
var dimensions = 1536

type Handler struct {
	store vectorstores.VectorStore
}

func NewHandler() (Handler, error) {
	e, err := embeddings.NewOpenAI()
	if err != nil {
		return Handler{}, err
	}

	store, err := pinecone.New(
		context.TODO(),
		pinecone.WithEmbedder(e),
		pinecone.WithIndexName("database"),
		pinecone.WithEnvironment("us-central1-gcp"),
		pinecone.WithProjectName("8fb28ba"),
		pinecone.WithNameSpace("temp"),
	)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		store: store,
	}, nil

}

func (h Handler) AddFile(fileData string) (fileId string, err error) {
	e, err := embeddings.NewOpenAI()
	if err != nil {
		return "", err
	}

	fileId = uuid.New().String()

	store, err := pinecone.New(
		context.TODO(),
		pinecone.WithEmbedder(e),
		pinecone.WithIndexName("database"),
		pinecone.WithEnvironment("us-central1-gcp"),
		pinecone.WithProjectName("8fb28ba"),
		pinecone.WithNameSpace(fileId),
	)
	if err != nil {
		return "", err
	}

	splitter := textsplitter.NewRecursiveCharacter()
	docs, err := textsplitter.CreateDocuments(splitter, []string{fileData}, nil)
	if err != nil {
		return "", err
	}

	err = store.AddDocuments(context.TODO(), docs)

	return fileId, err
}

func (h Handler) QueryFile(fileId string, query string) (string, error) {
	e, err := embeddings.NewOpenAI()
	if err != nil {
		return "", err
	}

	store, err := pinecone.New(
		context.TODO(),
		pinecone.WithEmbedder(e),
		pinecone.WithIndexName("database"),
		pinecone.WithEnvironment("us-central1-gcp"),
		pinecone.WithProjectName("8fb28ba"),
		pinecone.WithNameSpace(fileId),
	)
	if err != nil {
		return "", err
	}

	llm, err := openai.New()
	if err != nil {
		return "", err
	}

	return chains.Run(
		context.TODO(),
		chains.NewRetrievalQAFromLLM(llm, vectorstores.ToRetriever(store, 2)),
		query,
	)
}
