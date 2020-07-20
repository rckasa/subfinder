package digicert

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/projectdiscovery/subfinder/pkg/subscraping"
)

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		resp, err := session.Get(ctx, fmt.Sprintf("https://ssltools.digicert.com/chainTester/webservice/ctsearch/search?keyword=%s", domain), "", nil)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHttpResponse(resp)
			close(results)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			resp.Body.Close()
			close(results)
			return
		}
		resp.Body.Close()

		src := string(body)

		for _, subdomain := range session.Extractor.FindAllString(src, -1) {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain}
		}

		close(results)
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "digicert"
}
