package etsi119612

import (
	"context"
	"crypto/x509"
	"fmt"
)

// todo: add validation
// todo: make a big test for the walker --- need a lotl signed mock that contains pointers to lets say existing se-tl

type GraphUrls struct {
	adj map[string][]Edge
}

type Edge struct {
	URL   string
	Depth int
}

func NewGraph() *GraphUrls {
	return &GraphUrls{adj: make(map[string][]Edge)}
}

func (g *GraphUrls) AddEdge(parentURL string, edge Edge) {
	g.adj[parentURL] = append(g.adj[parentURL], edge)
}

type fetchFunctionType func(string) ([]byte, error)

var FetchTSLSwapBytesFunction fetchFunctionType = FetchTSLBytes




// better to have multiple options of cache
func GraphSearch(rootURL string, fetcher *cachedTSLFetcher, leafCert *x509.Certificate, intermediatePool *x509.CertPool) (*GraphUrls, error) {
	graph := NewGraph()
	visited := map[string]bool{}

	queue := []Edge{{URL: rootURL, Depth: 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.URL] {
			continue
		}
		visited[current.URL] = true
		tsl, bodyBytes, fetchErr := fetcher.Fetcher(context.Background(), current.URL)
		if fetchErr != nil {
			// no error return
			fmt.Errorf("error on fetch %w, %d", fetchErr, len(bodyBytes))
			continue
		}
		switch {
		case tsl.IsLOTL():
			fmt.Println("is lotl")
			links, fetchError := FetchPointersToOtherListTSL(tsl)
			if fetchError != nil {
				return nil, fmt.Errorf("error on fetch pointers to other tsl %v", fetchError)

			}
			for _, link := range links {
				child := Edge{URL: link, Depth: current.Depth + 1}
				graph.AddEdge(current.URL, child)
				if !visited[link] {
					queue = append(queue, child)
				}
			}
		case tsl.IsNationalTSL():
			policy := *PolicyAll
			pool := tsl.ToCertPool(&policy)
			_, verifyErr := leafCert.Verify(x509.VerifyOptions{
				Roots: pool, Intermediates: intermediatePool})
			if verifyErr != nil {
				links, err := FetchPointersToOtherListTSL(tsl)
				if err != nil {
					panic(err)
				}
				for _, link := range links {
					child := Edge{URL: link, Depth: current.Depth + 1}
					graph.AddEdge(current.URL, child)
					if !visited[link] {
						queue = append(queue, child)
					}
				}
			} else {
				return graph, nil
			}
		default:
			fmt.Printf("unknown TSLType %s\n",
				current.URL)

		}
	}

	return graph, fmt.Errorf("no national TSL verified the leaf")
}
