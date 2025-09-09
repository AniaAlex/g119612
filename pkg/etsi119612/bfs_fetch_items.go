package etsi119612

import (
	"context"
	"fmt"

	go_cache "github.com/eko/gocache/lib/v4/cache"
)

// todo: add validation
// todo: make a big test for the walker --- need a lotl signed mock that contains pointers to lets say existing se-tl

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
	//to check the graph algorythm only, otherwise multiple middle step
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.URL] {
			continue
		}
		ctx := context.Background()
		bodyBytes, err := cache.Get(ctx, current.URL)
		fmt.Printf("there %v", bodyBytes)

		//try to fetch
		if err != nil {
			fmt.Errorf("%v", err)
			bodyBytes, err := FetchTSLBytes(current.URL)
			if err == nil {
				//here add to cache
				err = cache.Set(ctx, current.URL, bodyBytes)
				visited[current.URL] = true
				if err != nil {
					panic(err)
				}
				tsl, err := UnmarshalCleanCerts(bodyBytes, current.URL)
				_ = tsl
				_ = err

				if err == nil {
					switch {
					case tsl.IsLOTL():
						fmt.Println("is lotl")
						links, err := FetchPontersToOtherListTSL(tsl)
						_ = links
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
					case tsl.IsNationalTSL():
						//here should be validation
						break
					}
				}
			}

			// 	// need to add signer somehow other way it later
			// 	tsl, err := UnmarshalCleanCerts(bodyBytes, current.URL)
			// 	if err == nil {
			// 		fmt.Printf("Here should be verification")
			// 		//if verification is successful then break
			// 		//here should be break
			// 	}

			// 	links, err := FetchPontersToOtherListTSL(tsl)
			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	for _, link := range links {
			// 		child := Edge{URL: link, Depth: current.Depth + 1}
			// 		graph.AddEdge(current.URL, child)
			// 		if !visited[link] {
			// 			queue = append(queue, child)
			// 		}
			// 	}

			// 	if current.Depth == 0 {
			// 		graph.AddEdge("", Edge{URL: rootURL, Depth: 0})
			// 	}

		}

	}
	fmt.Printf("%v", graph)
	return graph, nil
}
