package main

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	serverSession *gexec.Session
	suite         = "main"
)

func TestMainSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, suite+" Suite")
}

var _ = BeforeSuite(func() {
	var err error
	command := exec.Command("forego", "start")
	serverSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	//	Eventually(serverSession.Out, "10s").Should(gbytes.Say(`LOOKUPD\(mnementh.dev:4160\): peer info \{TcpPort:4160 HttpPort:4161 Version:0.3.2 BroadcastAddress:Mnementh.local\}`))
	Eventually(serverSession.Out, "10s").Should(gbytes.Say(`peer info`))
})

var _ = AfterSuite(func() {
	serverSession.Interrupt()
	gexec.CleanupBuildArtifacts()
})
