package actions

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/k8s/minikube"
	"github.com/guumaster/hostctl/pkg/profile"
	"github.com/guumaster/hostctl/pkg/types"
)

func newSyncMinikubeCmd(removeCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "minikube [profile] [flags]",
		Short: "Sync a minikube profile with a hostctl profile.",
		Long: `
Reads from Minikube the list of ingresses and add names and IPs to a profile in your hosts file.
`,
		Args: commonCheckArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			ns, _ := cmd.Flags().GetString("namespace")
			allNs, _ := cmd.Flags().GetBool("all-namespaces")

			if ns == "" && !allNs {
				return types.ErrKubernetesNamespace
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")
			ns, _ := cmd.Flags().GetString("namespace")
			allNs, _ := cmd.Flags().GetBool("all-namespace")

			if allNs {
				ns = ""
			}

			profileName := profiles[0]

			mini, err := minikube.GetProfile(profileName)
			if err != nil {
				return err
			}

			p, err := profile.NewProfileFromMinikube(mini, ns)
			if err != nil {
				return err
			}

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			p.Name = mini.Name
			p.Status = types.Enabled

			err = h.ReplaceProfile(p)
			if err != nil {
				return err
			}

			return h.Flush()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return postActionCmd(cmd, args, removeCmd, true)
		},
	}
}
