package gke

import (
	"fmt"

	"github.com/aquasecurity/tfsec/pkg/result"
	"github.com/aquasecurity/tfsec/pkg/severity"

	"github.com/aquasecurity/tfsec/pkg/provider"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/hclcontext"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"

	"github.com/aquasecurity/tfsec/pkg/rule"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		LegacyID:  "GCP009",
		Service:   "gke",
		ShortCode: "enforce-pod-security-policy",
		Documentation: rule.RuleDocumentation{
			Summary:    "Pod security policy enforcement not defined.",
			Impact:     "Pods could be operating with more permissions than required to be effective",
			Resolution: "Use security policies for pods to restrict permissions to those needed to be effective",
			Explanation: `
By default, Pods in Kubernetes can operate with capabilities beyond what they require. You should constrain the Pod's capabilities to only those required for that workload.

Kubernetes offers controls for restricting your Pods to execute with only explicitly granted capabilities. 

Pod Security Policy allows you to set smart defaults for your Pods, and enforce controls you want to enable across your fleet. 

The policies you define should be specific to the needs of your application
`,
			BadExample: `
resource "google_container_cluster" "bad_example" {
	pod_security_policy_config {
        enabled = "false"
	}
}`,
			GoodExample: `
resource "google_container_cluster" "good_example" {
	pod_security_policy_config {
        enabled = "true"
	}
}`,
			Links: []string{
				"https://cloud.google.com/kubernetes-engine/docs/how-to/hardening-your-cluster#admission_controllers",
				"https://www.terraform.io/docs/providers/google/r/container_cluster.html#pod_security_policy_config",
			},
		},
		Provider:        provider.GoogleProvider,
		RequiredTypes:   []string{"resource"},
		RequiredLabels:  []string{"google_container_cluster"},
		DefaultSeverity: severity.High,
		CheckFunc: func(set result.Set, resourceBlock block.Block, _ *hclcontext.Context) {

			pspBlock := resourceBlock.GetBlock("pod_security_policy_config")
			if pspBlock == nil {
				set.Add(
					result.New(resourceBlock).
						WithDescription(fmt.Sprintf("Resource '%s' defines a cluster with no Pod Security Policy config defined. It is recommended to define a PSP for your pods and enable PSP enforcement.", resourceBlock.FullName())).
						WithRange(resourceBlock.Range()),
				)
				return
			}

			enforcePSP := pspBlock.GetAttribute("enabled")
			if enforcePSP != nil && enforcePSP.IsFalse() {
				set.Add(
					result.New(resourceBlock).
						WithDescription(fmt.Sprintf("Resource '%s' defines a cluster with Pod Security Policy enforcement disabled. It is recommended to define a PSP for your pods and enable PSP enforcement.", resourceBlock.FullName())).
						WithRange(enforcePSP.Range()),
				)
			}

		},
	})
}
