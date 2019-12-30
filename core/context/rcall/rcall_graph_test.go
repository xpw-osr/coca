package rcall

import (
	"encoding/json"
	"github.com/phodal/coca/core/domain"
	"github.com/phodal/coca/core/infrastructure"
	"log"
	"testing"

	. "github.com/onsi/gomega"
)

func TestRCallGraph_Analysis(t *testing.T) {
	g := NewGomegaWithT(t)

	var parsedDeps []domain.JClassNode
	analyser := NewRCallGraph()
	file := infrastructure.ReadFile("../../../_fixtures/call/call_api_test.json")
	if file == nil {
		log.Fatal("lost file")
	}

	_ = json.Unmarshal(file, &parsedDeps)

	content := analyser.Analysis("com.phodal.pholedge.book.BookFactory.create", *&parsedDeps)

	g.Expect(content).To(Equal(`digraph G {
rankdir = LR;
edge [dir="back"];

}
`))
}