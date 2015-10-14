package main

import (
	"log"
	"os"
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

func RemoveDataFiles() {
	dirName := os.ExpandEnv("${GOPATH}/src/github.com/duskhacker/cqrsnu/data")
	dir, err := os.Open(dirName)
	if err != nil {
		log.Fatalf("error opening %s: %s", dirName, err)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatalf("error reading dir %s: %s\n", dir.Name(), err)
	}

	for _, file := range files {
		os.Remove(dir.Name() + "/" + file.Name())
	}
}

var _ = BeforeSuite(func() {
	var err error

	RemoveDataFiles()

	command := exec.Command("forego", "start")
	serverSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	Eventually(serverSession, "10s").Should(gbytes.Say(`peer info`))
})

var _ = AfterSuite(func() {
	serverSession.Interrupt()
	gexec.CleanupBuildArtifacts()
})
