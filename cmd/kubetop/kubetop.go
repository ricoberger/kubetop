package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ricoberger/kubetop/pkg/api"
	"github.com/ricoberger/kubetop/pkg/term"
	"github.com/ricoberger/kubetop/pkg/term/widgets"
	"github.com/ricoberger/kubetop/pkg/version"

	"github.com/spf13/cobra"
)

var (
	kubeconfig string
	namespace  string
)

var rootCmd = &cobra.Command{
	Use:   "kubetop",
	Short: "Display Resource (CPU/Memory/Storage) usage of nodes and pods",
	Long:  "Display Resource (CPU/Memory/Storage) usage of nodes and pods",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Kubernetes API client.
		client, err := api.NewClient(kubeconfig)
		if err != nil {
			log.Fatalf("Failed to initialize API client: %#v", err)
		}

		// Initialize and run the terminal user interface for kubetop.
		t := term.Term{
			APIClient: client,
			ViewType:  widgets.ViewTypePods,
		}

		err = t.Run(api.Filter{Namespace: namespace, Node: "", Status: 10})
		if err != nil {
			log.Fatalf("Failed to initialize ui: %#v", err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information for kubetop",
	Long:  "Print version information for kubetop",
	Run: func(cmd *cobra.Command, args []string) {
		// Print the version information for kubetop.
		v, err := version.Print("kubetop")
		if err != nil {
			log.Fatalf("Failed to print version information: %#v", err)
		}

		fmt.Fprintln(os.Stdout, v)
		os.Exit(0)
	},
}

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Display Resource (CPU/Memory/Storage) usage of nodes",
	Long:  "Display Resource (CPU/Memory/Storage) usage of nodes",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Kubernetes API client.
		client, err := api.NewClient(kubeconfig)
		if err != nil {
			log.Fatalf("Failed to initialize API client: %#v", err)
		}

		// Initialize and run the terminal user interface for kubetop.
		t := term.Term{
			APIClient: client,
			ViewType:  widgets.ViewTypeNodes,
		}

		err = t.Run(api.Filter{Namespace: namespace, Node: "", Status: 10})
		if err != nil {
			log.Fatalf("Failed to initialize ui: %#v", err)
		}
	},
}

var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Display Resource (CPU/Memory/Storage) usage of pods",
	Long:  "Display Resource (CPU/Memory/Storage) usage of pods",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Kubernetes API client.
		client, err := api.NewClient(kubeconfig)
		if err != nil {
			log.Fatalf("Failed to initialize API client: %#v", err)
		}

		// Initialize and run the terminal user interface for kubetop.
		t := term.Term{
			APIClient: client,
			ViewType:  widgets.ViewTypePods,
		}

		err = t.Run(api.Filter{Namespace: namespace, Node: "", Status: 10})
		if err != nil {
			log.Fatalf("Failed to initialize ui: %#v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(nodesCmd)
	rootCmd.AddCommand(podsCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "If present, the namespace scope for this CLI request")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to initialize kubetop: %#v", err)
	}
}
