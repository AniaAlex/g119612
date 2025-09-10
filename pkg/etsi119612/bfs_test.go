package etsi119612

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/SUNET/g119612/pkg/cache"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

// can import it later change package to another package
type JWTCertBundle struct {
	Alg string   `json:"alg"`
	Typ string   `json:"typ"`
	X5c []string `json:"x5c"`
}

func TestGraph(t *testing.T) {
	header_mock, err := os.ReadFile("./testdata/x5c-test-root-leaf.json")
	if err != nil {
		t.Fatalf("Failed while reading json: %v", err)
	}
	assert.NotEmpty(t, header_mock)
	var jwt JWTCertBundle
	err = json.Unmarshal(header_mock, &jwt)
	if err != nil {
		t.Fatalf("Failed updating jwt bundle")
	}
	leafDER, err := base64.StdEncoding.DecodeString(jwt.X5c[0])
	assert.NoError(t, err)
	leafCert, err := x509.ParseCertificate(leafDER)
	assert.NoError(t, err)
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

	graph, error := GraphSearch("https://ec.europa.eu/tools/lotl/eu-lotl.xml", newCache, leafCert)

	if error != nil {
		t.Fatal(error)
	}
	t.Logf("%v", graph.adj)
}

func TestGraphLOTLInCache(t *testing.T) {
	header_mock, err := os.ReadFile("./testdata/x5c-test-root-leaf.json")
	if err != nil {
		t.Fatalf("Failed while reading json: %v", err)
	}
	assert.NotEmpty(t, header_mock)
	var jwt JWTCertBundle
	err = json.Unmarshal(header_mock, &jwt)
	if err != nil {
		t.Fatalf("Failed updating jwt bundle")
	}
	leafDER, err := base64.StdEncoding.DecodeString(jwt.X5c[0])
	assert.NoError(t, err)
	leafCert, err := x509.ParseCertificate(leafDER)
	assert.NoError(t, err)
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
	bodyBytes, err := FetchTSLBytes("https://ec.europa.eu/tools/lotl/eu-lotl.xml")
	assert.NoError(t, err)
	ctx := context.Background()
	err = newCache.Set(ctx, "https://ec.europa.eu/tools/lotl/eu-lotl.xml", bodyBytes)
	assert.NoError(t, err)
	graph, error := GraphSearch("https://ec.europa.eu/tools/lotl/eu-lotl.xml", newCache, leafCert)

	if error != nil {
		t.Fatal(error)
	}
	t.Logf("%v", graph.adj)
}

func TestGraphSELOTLInCache(t *testing.T) {
	header_mock, err := os.ReadFile("./testdata/x5c-test-root-leaf.json")
	if err != nil {
		t.Fatalf("Failed while reading json: %v", err)
	}
	assert.NotEmpty(t, header_mock)
	var jwt JWTCertBundle
	err = json.Unmarshal(header_mock, &jwt)
	if err != nil {
		t.Fatalf("Failed updating jwt bundle")
	}
	leafDER, err := base64.StdEncoding.DecodeString(jwt.X5c[0])
	assert.NoError(t, err)
	leafCert, err := x509.ParseCertificate(leafDER)
	assert.NoError(t, err)
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
	lotlbodyBytes, err := FetchTSLBytes("https://ec.europa.eu/tools/lotl/eu-lotl.xml")
	assert.NoError(t, err)
	ctx := context.Background()
	err = newCache.Set(ctx, "https://ec.europa.eu/tools/lotl/eu-lotl.xml", lotlbodyBytes)
	assert.NoError(t, err)
	slbodyBytes, err := FetchTSLBytes("https://trustedlist.pts.se/SE-TL.xml")
	assert.NoError(t, err)
	err = newCache.Set(ctx, "https://trustedlist.pts.se/SE-TL.xml", slbodyBytes)
	assert.NoError(t, err)
	graph, error := GraphSearch("https://ec.europa.eu/tools/lotl/eu-lotl.xml", newCache, leafCert)

	if error != nil {
		t.Fatal(error)
	}
	t.Logf("%v", graph.adj)
}
