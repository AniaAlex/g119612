package etsi119612

import (
	"context"
	"crypto/x509"
	"fmt"

	go_cache "github.com/eko/gocache/lib/v4/cache"
)

// pointer here insteda of value?
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

// better to have multiple options of cache
func GraphSearch(rootURL string, cache *go_cache.Cache[[]byte]) (*GraphUrls, error) {
	graph := NewGraph()
	visited := map[string]bool{}

	queue := []Edge{{URL: rootURL, Depth: 0}}

	//should be in loop
	//to check the graph algorythm only, otherwise multiple middle steps

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.URL] {
			continue
		}
		ctx := context.Background()
		bodyBytes, err := cache.Get(ctx, current.URL)

		//try to fetch
		if err != nil {
			fmt.Errorf("%v", err)
			bodyBytes, signer, err := FetchTSLBytes(current.URL)
			if err == nil {
				//here add to cache
				err = cache.Set(ctx, "https://ewc-consortium.github.io", bodyBytes)
				visited[current.URL] = true
				if err != nil {
					panic(err)
				}
				tsl, err := UnmarshalCleanCerts(bodyBytes, signer, current.URL)

				if err == nil {
					fmt.Println("Here should be verification:%v", tsl)
					break
				}

			}

		}

		var signer x509.Certificate
		// need to add signer somehow other way it later
		tsl, err := UnmarshalCleanCerts(bodyBytes, signer, current.URL)
		if err == nil {
			fmt.Println("Here should be verification:%v", tsl)
			//if verification is successful then break
			break
		}

		//if verification failed but we have tsl, we add children to continue with the loop
		// here should list the pointers in the tsl
		links := []string{"https://trustedlist.pts.se/SE-TL.xml", "https://trustedlist.pts.se/NL-TL.xml"}

		for _, link := range links {

			child := Edge{URL: link, Depth: current.Depth + 1}
			graph.AddEdge(current.URL, child)
			if !visited[link] {
				queue = append(queue, child)
			}

		}

		if current.Depth == 0 {
			graph.AddEdge("", Edge{URL: rootURL, Depth: 0})
		}

	}

	return graph, nil
}
