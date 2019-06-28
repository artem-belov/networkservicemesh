// +build interdomain

package nsmd_integration_tests

import (
	"fmt"
	"github.com/networkservicemesh/networkservicemesh/test/kubetest/pods"
	v1 "k8s.io/api/core/v1"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/networkservicemesh/networkservicemesh/test/kubetest"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestInterdomainNSCDies(t *testing.T) {
	RegisterTestingT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	testInterdomainNSMDies(t, 2, true)
}

func TestInterdomainNSEDies(t *testing.T) {
	RegisterTestingT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	testInterdomainNSMDies(t, 2, false)
}

/**
If passed 1 both will be on same node, if not on different.
*/
func testInterdomainNSMDies(t *testing.T, clustersCount int, killSrc bool) {
	k8ss := []* kubetest.ExtK8s{}

	for i := 0; i < clustersCount; i++ {
		kubeconfig := os.Getenv(fmt.Sprintf("KUBECONFIG_CLUSTER_%d", i + 1))
		Expect(len(kubeconfig)).ToNot(Equal(0))

		k8s, err := kubetest.NewK8sForConfig(true, kubeconfig)

		Expect(err).To(BeNil())

		nseNoHeal.Namespace = k8s.GetK8sNamespace()
		nseNoHeal.DataplaneVariables = kubetest.DefaultDataplaneVariables(k8s.GetForwardingPlane())

		nodesSetup, err := kubetest.SetupNodesConfig(k8s, 1, defaultTimeout, []*pods.NSMgrPodConfig{
			nseNoHeal,
			nseNoHeal,
		}, k8s.GetK8sNamespace())
		Expect(err).To(BeNil())

		k8ss = append(k8ss, &kubetest.ExtK8s{
			K8s:      k8s,
			NodesSetup: nodesSetup,
		})

		pnsmdName := fmt.Sprintf("pnsmgr-%s", nodesSetup[0].Node.Name)
		kubetest.DeployProxyNSMgr(k8s, nodesSetup[0].Node, pnsmdName, defaultTimeout)

		serviceCleanup := kubetest.RunProxyNSMgrService(k8s)
		defer serviceCleanup()

		defer k8ss[i].K8s.Cleanup()
	}

	// Run ICMP
	icmpPodNode := kubetest.DeployICMP(k8ss[clustersCount - 1].K8s, k8ss[clustersCount - 1].NodesSetup[0].Node, "icmp-responder-nse-1", defaultTimeout)

	nseExternalIP, err := kubetest.GetNodeExternalIP(k8ss[clustersCount - 1].NodesSetup[0].Node)
	if err != nil {
		nseExternalIP, err = kubetest.GetNodeInternalIP(k8ss[clustersCount - 1].NodesSetup[0].Node)
		Expect(err).To(BeNil())
	}

	nscPodNode := kubetest.DeployNSCWithEnv(k8ss[0].K8s, k8ss[0].NodesSetup[0].Node, "nsc-1", defaultTimeout, map[string]string{
		"OUTGOING_NSC_LABELS": "app=icmp",
		"OUTGOING_NSC_NAME":   fmt.Sprintf("icmp-responder@%s", nseExternalIP),
	})

	var nscInfo *kubetest.NSCCheckInfo

	failures := InterceptGomegaFailures(func() {
		nscInfo = kubetest.CheckNSC(k8ss[0].K8s, nscPodNode)

		ipResponse, errOut, err := k8ss[clustersCount - 1].K8s.Exec(icmpPodNode, icmpPodNode.Spec.Containers[0].Name, "ip", "addr")
		Expect(err).To(BeNil())
		Expect(errOut).To(Equal(""))
		Expect(strings.Contains(ipResponse, "nsm")).To(Equal(true))
	})
	// Do dumping of container state to dig into what is happened.
	if len(failures) > 0 {
		logrus.Errorf("Failures: %v", failures)
		for i := 0; i < clustersCount; i++ {
			kubetest.PrintLogs(k8ss[i].K8s, k8ss[i].NodesSetup)
		}
		nscInfo.PrintLogs()

		t.Fail()
	}

	var podToKill *v1.Pod
	var clusterToKill int
	var podToCheck *v1.Pod
	var clusterToCheck int
	if killSrc {
		podToKill = nscPodNode
		clusterToKill = 0
		podToCheck = icmpPodNode
		clusterToCheck = clustersCount - 1
	} else {
		podToKill = icmpPodNode
		clusterToKill = clustersCount - 1
		podToCheck = nscPodNode
		clusterToCheck = 0
	}

	k8ss[clusterToKill].K8s.DeletePods(podToKill)

	success := false
	for attempt := 0; attempt < 20; <-time.After(300 * time.Millisecond) {
		attempt++
		ipResponse, errOut, err := k8ss[clusterToCheck].K8s.Exec(podToCheck, podToCheck.Spec.Containers[0].Name, "ip", "addr")
		if err == nil && errOut == "" && !strings.Contains(ipResponse, "nsm") {
			success = true
			break
		}
	}

	Expect(success).To(Equal(true))
}

