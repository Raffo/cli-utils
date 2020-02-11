// Copyright 2020 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package destroy

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"sigs.k8s.io/cli-utils/pkg/apply"
)

// NewCmdDestroy creates the `destroy` command
func NewCmdDestroy(f util.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	destroyer := apply.NewDestroyer(f, ioStreams)
	printer := &apply.BasicPrinter{
		IOStreams: ioStreams,
	}

	cmd := &cobra.Command{
		Use:                   "destroy (-f FILENAME | -k DIRECTORY)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Destroy all the resources related to configuration managed by kpt"),
		Args:                  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				destroyer.ApplyOptions.DeleteFlags.FileNameFlags.Kustomize = &args[0]
			}

			cmdutil.CheckErr(destroyer.Initialize(cmd))

			// Run the destroyer. It will return a channel where we can receive updates
			// to keep track of progress and any issues.
			ch := destroyer.Run()

			// The printer will print updates from the channel. It will block
			// until the channel is closed.
			printer.Print(ch)
		},
	}

	destroyer.SetFlags(cmd)

	cmdutil.AddValidateFlags(cmd)
	cmd.Flags().BoolVar(&destroyer.DryRun, "dry-run", destroyer.DryRun,
		"If true, only print the object that would be sent and action which would be performed, without performing it.")
	cmdutil.AddServerSideApplyFlags(cmd)
	return cmd
}