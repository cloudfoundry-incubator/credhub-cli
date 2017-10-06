package credhub_test

import (
	"net/http"

	. "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	version "github.com/hashicorp/go-version"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ServerVersion()", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/info"),
				ghttp.RespondWith(http.StatusOK, `{
							"auth-server": {
								"url": "https://uaa.example.com:8443"
							},
							"app": {
								"name": "CredHub",
								"version": "1.2.3"
							}
						}`)),
		)
	})

	AfterEach(func() {
		server.Close()
	})

	It("should obtain the server version from the /info endpoint for the first request", func() {
		expectedVersion, err := version.NewVersion("1.2.3")
		Expect(err).To(BeNil())

		ch, err := New(server.URL())
		Expect(err).To(BeNil())

		serverVersion, err := ch.ServerVersion()
		Expect(err).To(BeNil())
		Expect(serverVersion).To(Equal(expectedVersion))

		serverVersion, err = ch.ServerVersion()
		Expect(err).To(BeNil())
		Expect(serverVersion).To(Equal(expectedVersion))

		Expect(server.ReceivedRequests()).Should(HaveLen(1))
	})
})
