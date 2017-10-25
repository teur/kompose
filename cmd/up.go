/*
Copyright 2017 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	log "github.com/Sirupsen/logrus"

	"github.com/kubernetes/kompose/pkg/app"
	"github.com/kubernetes/kompose/pkg/kobject"
	"github.com/spf13/cobra"
)

// TODO: comment
var (
	UpReplicas     int
	UpEmptyVols    bool
	UpVolumes      string
	UpInsecureRepo bool
	UpNamespace    string
	UpOpt          kobject.ConvertOptions
	UpBuild        string
	UpBuildBranch  string
	UpBuildRepo    string
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Deploy your Dockerized application to a container orchestrator.",
	Long:  `Deploy your Dockerized application to a container orchestrator. (default "kubernetes")`,
	PreRun: func(cmd *cobra.Command, args []string) {

		// Check that build-config wasn't passed in with --provider=kubernetes
		if GlobalProvider == "kubernetes" && UpBuild == "build-config" {
			log.Fatalf("build-config is not a valid --build parameter with provider Kubernetes")
		}

		// Create the Convert options.
		UpOpt = kobject.ConvertOptions{
			Build:              UpBuild,
			Replicas:           UpReplicas,
			InputFiles:         GlobalFiles,
			Provider:           GlobalProvider,
			EmptyVols:          UpEmptyVols,
			Volumes:            UpVolumes,
			Namespace:          UpNamespace,
			InsecureRepository: UpInsecureRepo,
			BuildBranch:        UpBuildBranch,
			BuildRepo:          UpBuildRepo,
			IsNamespaceFlag:    cmd.Flags().Lookup("namespace").Changed,
		}

		// Validate before doing anything else.
		app.ValidateComposeFile(&UpOpt)
	},
	Run: func(cmd *cobra.Command, args []string) {
		app.Up(UpOpt)
	},
}

func init() {
	upCmd.Flags().IntVar(&UpReplicas, "replicas", 1, "Specify the number of replicas generated")
	upCmd.Flags().StringVar(&UpVolumes, "volumes", "persistentVolumeClaim", `Volumes to be generated ("persistentVolumeClaim"|"emptyDir")`)
	upCmd.Flags().BoolVar(&UpInsecureRepo, "insecure-repository", false, "Use an insecure Docker repository for OpenShift ImageStream")
	upCmd.Flags().StringVar(&UpNamespace, "namespace", "default", "Specify Namespace to deploy your application")
	upCmd.Flags().StringVar(&UpBuild, "build", "local", `Set the type of build ("local"|"build-config" (OpenShift only)|"none")`)
	upCmd.Flags().StringVar(&UpBuildRepo, "build-repo", "", "Specify source repository for buildconfig (default remote origin)")
	upCmd.Flags().StringVar(&UpBuildBranch, "build-branch", "", "Specify repository branch to use for buildconfig (default master)")
	upCmd.Flags().MarkHidden("insecure-repository")
	upCmd.Flags().MarkHidden("build-repo")
	upCmd.Flags().MarkHidden("build-branch")

	// Deprecated
	upCmd.Flags().BoolVar(&UpEmptyVols, "emptyvols", false, "Use empty volumes. Do not generate PersistentVolumeClaim")
	upCmd.Flags().MarkDeprecated("emptyvols", "emptyvols has been marked as deprecated. Use --volumes empty")

	// In order to 'separate' both OpenShift and Kubernetes only flags. A custom help page is created
	customHelp := `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}

OpenShift Flags:
      --build-branch             Specify repository branch to use for buildconfig (default is current branch name)
      --build-repo               Specify source repository for buildconfig (default is current branch's remote url
      --insecure-repository      Specify to use insecure docker repository while generating Openshift image stream object

Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
	// Set the help template + add the command to root
	upCmd.SetUsageTemplate(customHelp)

	RootCmd.AddCommand(upCmd)
}
