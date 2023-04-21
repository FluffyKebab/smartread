package text

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/exp/chains"
	"github.com/tmc/langchaingo/exp/documentLoaders"
	"github.com/tmc/langchaingo/exp/embeddings"
	"github.com/tmc/langchaingo/exp/textSplitters"
	"github.com/tmc/langchaingo/exp/vectorStores/pinecone"
	"github.com/tmc/langchaingo/llms/openai"
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
	docs, err := documentLoaders.NewTextLoaderFromFile(fileData).LoadAndSplit(splitter)
	if err != nil {
		return "", err
	}

	fileId = uuid.New().String()
	h.pineconeIndex.AddDocuments(docs, []string{}, fileId)

	return fileId, nil
}

func (h Handler) QueryFile(fileId string, query string) (string, error) {
	llm, err := openai.New()
	if err != nil {
		return "", err
	}

	chain := chains.NewRetrievalQAChainFromLLM(llm, h.pineconeIndex.ToRetriever(5, fileId))

	resultMap, err := chains.Call(chain, map[string]any{
		"query": query,
	})

	resultAny, ok := resultMap["text"]
	if !ok {
		return "", fmt.Errorf("Missing text field in QAchain result")
	}

	result, ok := resultAny.(string)
	if !ok {
		return "", fmt.Errorf("Text field of wrong type")
	}

	return result, err
}
