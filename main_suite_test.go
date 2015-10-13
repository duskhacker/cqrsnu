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
	pathToCLI, pathToServer string
	//	thriftAddress           = "-thrift-address=localhost:9080"
	serverSession *gexec.Session
	suite         = "main"

//	serverConfig            string
//	rng                     *rand.Rand
)

func TestMainSuite(t *testing.T) {
	RegisterFailHandler(Fail)

	//	baseDir := os.ExpandEnv("${GOPATH}/src/github.com/duskhacker/cqrsnu")

	//	ci := os.Getenv("CI")
	//	if len(ci) > 0 {
	//		iniflags.SetConfigFile(fmt.Sprintf("%s/config/ci.ini", baseDir))
	//		serverConfig = fmt.Sprintf("-config=%s/config/ci.ini", baseDir)
	//		junitReporter := reporters.NewJUnitReporter(fmt.Sprintf("%s/%s_junit.xml", baseDir, strings.ToLower(suite)))
	//		RunSpecsWithDefaultAndCustomReporters(t, suite+" Suite", []Reporter{junitReporter})
	//	} else {
	//		serverConfig = fmt.Sprintf("-config=%s/config/test.ini", baseDir)
	//		iniflags.SetConfigFile(fmt.Sprintf("%s/config/test.ini", baseDir))
	RunSpecs(t, suite+" Suite")
	//	}
}

var _ = BeforeSuite(func() {
	var err error
	//	cfg := config.Parse()

	//	dbsetup.ResetTestDatabase(*cfg)
	//	dbsetup.LoadSchemaInTest(*cfg)

	//	pathToCLI, err = gexec.Build("dcx.rax.io/layer3/cmd/l3-client")
	//	Expect(err).ToNot(HaveOccurred())

	//	pathToServer, err = gexec.Build("dcx.rax.io/layer3/cmd/l3-server")
	//	Expect(err).ToNot(HaveOccurred())

	//	command := exec.Command(pathToServer, thriftAddress, serverConfig)
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
