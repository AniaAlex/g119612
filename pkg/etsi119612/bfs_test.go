package etsi119612

import (
	"testing"

	"github.com/SUNET/g119612/pkg/cache"
	"github.com/h2non/gock"
)

func TestGraph(t *testing.T) {
	gock.New("https://eidas.agid.gov.it").
		Get("/TL/TSL-IT.xml").
		Reply(200).
		File("./testdata/testdata_walker/signed_it-tsl.xml")

	defer gock.Off()
	gock.New("https://trustedlist.pts.se").
		Get("/SE-TL.xml").
		Reply(200).
		File("./testdata/testdata_walker/signed_se-tl.xml")

	defer gock.Off()
	gock.New("https://ec.europa.eu").
		Get("/tools/lotl/eu-lotl.xml").
		Reply(200).
		File("./testdata/testdata_walker/signed_lotl.xml")

	cachesets := cache.CacheSettings{Backend: cache.BackendGoCache}
	deps := cache.Dependencies{}
	deps.DefaultOutput()
	//probably mock here better
	newCache := cache.NewCache[[]byte](cachesets, deps)

	graph, error := GraphSearch("https://ec.europa.eu/tools/lotl/eu-lotl.xml", newCache)

	if error != nil {
		t.Fatal(error)
	}
	t.Logf("%v", graph.adj)
}
