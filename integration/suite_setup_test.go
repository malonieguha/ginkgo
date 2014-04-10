package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("SuiteSetup", func() {
	var pathToTest string

	Context("when the BeforeSuite and AfterSuite pass", func() {
		BeforeEach(func() {
			pathToTest = tmpPath("suite_setup")
			copyIn("passing_suite_setup", pathToTest)
		})

		It("should run the BeforeSuite once, then run all the tests", func() {
			output, err := runGinkgo(pathToTest, "--noColor")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(strings.Count(output, "BEFORE SUITE")).Should(Equal(1))
			Ω(strings.Count(output, "AFTER SUITE")).Should(Equal(1))
		})

		It("should run the BeforeSuite once per parallel node, then run all the tests", func() {
			output, err := runGinkgo(pathToTest, "--noColor", "--nodes=2")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(strings.Count(output, "BEFORE SUITE")).Should(Equal(2))
			Ω(strings.Count(output, "AFTER SUITE")).Should(Equal(2))
		})
	})

	Context("when the BeforeSuite fails", func() {
		BeforeEach(func() {
			pathToTest = tmpPath("suite_setup")
			copyIn("failing_before_suite", pathToTest)
		})

		It("should run the BeforeSuite once, none of the tests, but it should run the AfterSuite", func() {
			output, err := runGinkgo(pathToTest, "--noColor")
			Ω(err).Should(HaveOccurred())
			Ω(strings.Count(output, "BEFORE SUITE")).Should(Equal(1))
			Ω(strings.Count(output, "Test Panicked")).Should(Equal(1))
			Ω(strings.Count(output, "AFTER SUITE")).Should(Equal(1))
			Ω(output).ShouldNot(ContainSubstring("NEVER SEE THIS"))
		})

		It("should run the BeforeSuite once per parallel node, none of the tests, but it should run the AfterSuite for each node", func() {
			output, err := runGinkgo(pathToTest, "--noColor", "--nodes=2")
			Ω(err).Should(HaveOccurred())
			Ω(strings.Count(output, "BEFORE SUITE")).Should(Equal(2))
			Ω(strings.Count(output, "Test Panicked")).Should(Equal(2))
			Ω(strings.Count(output, "AFTER SUITE")).Should(Equal(2))
			Ω(output).ShouldNot(ContainSubstring("NEVER SEE THIS"))
		})
	})

	Context("when the AfterSuite fails", func() {
		BeforeEach(func() {
			pathToTest = tmpPath("suite_setup")
			copyIn("failing_after_suite", pathToTest)
		})

		It("should run the BeforeSuite once, none of the tests, but it should run the AfterSuite", func() {
			output, err := runGinkgo(pathToTest, "--noColor")
			Ω(err).Should(HaveOccurred())
			Ω(strings.Count(output, "BEFORE SUITE")).Should(Equal(1))
			Ω(strings.Count(output, "AFTER SUITE")).Should(Equal(1))
			Ω(strings.Count(output, "Test Panicked")).Should(Equal(1))
			Ω(strings.Count(output, "A TEST")).Should(Equal(2))
		})

		It("should run the BeforeSuite once per parallel node, none of the tests, but it should run the AfterSuite for each node", func() {
			output, err := runGinkgo(pathToTest, "--noColor", "--nodes=2")
			Ω(err).Should(HaveOccurred())
			Ω(strings.Count(output, "BEFORE SUITE")).Should(Equal(2))
			Ω(strings.Count(output, "AFTER SUITE")).Should(Equal(2))
			Ω(strings.Count(output, "Test Panicked")).Should(Equal(2))
			Ω(strings.Count(output, "A TEST")).Should(Equal(2))
		})
	})

	Context("With compound before and after suites", func() {
		BeforeEach(func() {
			pathToTest = tmpPath("suite_setup")
			copyIn("compound_setup_tests", pathToTest)
		})

		Context("when run with one node", func() {
			It("should do all the work on that one node", func() {
				output, err := runGinkgo(pathToTest, "--noColor")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(output).Should(ContainSubstring("BEFORE_A_1\nBEFORE_B_1: DATA"))
				Ω(output).Should(ContainSubstring("AFTER_A_1\nAFTER_B_1"))
			})
		})

		Context("when run across multiple nodes", func() {
			It("should run the first BeforeSuite function (BEFORE_A) on node 1, the second (BEFORE_B) on all the nodes, the first AfterSuite (AFTER_A) on all the nodes, and then the second (AFTER_B) on Node 1 *after* everything else is finished", func() {
				output, err := runGinkgo(pathToTest, "--noColor", "--nodes=3")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(output).Should(ContainSubstring("BEFORE_A_1"))
				Ω(output).Should(ContainSubstring("BEFORE_B_1: DATA"))
				Ω(output).Should(ContainSubstring("BEFORE_B_2: DATA"))
				Ω(output).Should(ContainSubstring("BEFORE_B_3: DATA"))

				Ω(output).ShouldNot(ContainSubstring("BEFORE_A_2"))
				Ω(output).ShouldNot(ContainSubstring("BEFORE_A_3"))

				Ω(output).Should(ContainSubstring("AFTER_A_1"))
				Ω(output).Should(ContainSubstring("AFTER_A_2"))
				Ω(output).Should(ContainSubstring("AFTER_A_3"))
				Ω(output).Should(ContainSubstring("AFTER_B_1"))

				Ω(output).ShouldNot(ContainSubstring("AFTER_B_2"))
				Ω(output).ShouldNot(ContainSubstring("AFTER_B_3"))
			})
		})

		Context("when streaming across multiple nodes", func() {
			It("should run the first BeforeSuite function (BEFORE_A) on node 1, the second (BEFORE_B) on all the nodes, the first AfterSuite (AFTER_A) on all the nodes, and then the second (AFTER_B) on Node 1 *after* everything else is finished", func() {
				output, err := runGinkgo(pathToTest, "--noColor", "--nodes=3", "--stream")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(output).Should(ContainSubstring("[1] BEFORE_A_1"))
				Ω(output).Should(ContainSubstring("[1] BEFORE_B_1: DATA"))
				Ω(output).Should(ContainSubstring("[2] BEFORE_B_2: DATA"))
				Ω(output).Should(ContainSubstring("[3] BEFORE_B_3: DATA"))

				Ω(output).ShouldNot(ContainSubstring("BEFORE_A_2"))
				Ω(output).ShouldNot(ContainSubstring("BEFORE_A_3"))

				Ω(output).Should(ContainSubstring("[1] AFTER_A_1"))
				Ω(output).Should(ContainSubstring("[2] AFTER_A_2"))
				Ω(output).Should(ContainSubstring("[3] AFTER_A_3"))
				Ω(output).Should(ContainSubstring("[1] AFTER_B_1"))

				Ω(output).ShouldNot(ContainSubstring("AFTER_B_2"))
				Ω(output).ShouldNot(ContainSubstring("AFTER_B_3"))
			})
		})
	})
})
