package etsi119612

import (
	"testing"

	"github.com/SUNET/g119612/pkg/cache"
	"github.com/h2non/gock"
)

func TestGraph(t *testing.T) {

	defer gock.Off()
	gock.New("https://trustedlist.pts.se").
		Get("/SE-TL.xml").
		Reply(200).
		File("./testdata/SE-TL.xml")

	defer gock.Off()
	gock.New("https://trustedlist.pts.se").
		Get("/SE-TL.xml").
		Reply(200).
		File("./testdata/SE-TL.xml")

	defer gock.Off()
	gock.New("https://ec.europa.eu/tools/lotl").
		Get("/eu-lotl.xml").
		Reply(200).
		File("./testdata/EWC-TL.xml")

	defer gock.Off()
	gock.New("https://ec.europa.eu/tools/lotl").
		Get("/eu-lotl.xml").
		Reply(200).
		File("./testdata/EWC-TL.xml")

	cachesets := cache.CacheSettings{Backend: cache.BackendGoCache}
	deps := cache.Dependencies{}
	deps.DefaultOutput()
	//probably mock here better
	newCache := cache.NewCache[[]byte](cachesets, deps)

	graph, error := GraphSearch("https://trustedlist.pts.se/SE-TL.xml", newCache)

	if error != nil {
		t.Fatal(error)
	}
	t.Logf("%v", graph.adj)
}
