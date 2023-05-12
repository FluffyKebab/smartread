package text

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/exp/embeddings"
	"github.com/tmc/langchaingo/exp/schema"
	"github.com/tmc/langchaingo/exp/textSplitters"
	"github.com/tmc/langchaingo/exp/vectorStores/pinecone"
)

var pineconeEnv = "us-central1-gcp"
var indexName = "database"
var dimensions = 1536

type Handler struct {
	pineconeIndex pinecone.Pinecone
}

func NewHandler() (Handler, error) {
	e, err := embeddings.NewOpenAI()
	if err != nil {
		return Handler{}, err
	}

	index, err := pinecone.NewPinecone(e, pineconeEnv, indexName, dimensions)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		pineconeIndex: index,
	}, nil
}

func (h Handler) AddFile(fileData string) (fileId string, err error) {

	splitter := textSplitters.NewRecursiveCharactersSplitter()
	docs, err := textSplitters.SplitDocuments(splitter, []schema.Document{{PageContent: fileData}})
	if err != nil {
		return "", err
	}

	fmt.Println("docs: ", docs)

	fileId = uuid.New().String()
	h.pineconeIndex.AddDocuments(docs, []string{}, fileId)

	return fileId, nil
}

func (h Handler) QueryFile(fileId string, query string) (string, error) {
	return "wow dette er en vrldig god response", nil
}
