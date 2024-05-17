package info

import (
    "fmt"
)

// DisplayHelp displays the help message for the k0pt command-line tool
func DisplayHelp() {
    fmt.Println(`
  oooo          .oooo.                  .   
  ` + "`" + `888         d8P'` + "`" + `Y8b               .o8   
   888  oooo  888    888 oo.ooooo.  .o888oo 
   888 .8P'   888    888  888' ` + "`" + `88b   888   
   888888.    888    888  888   888   888   
   888 ` + "`" + `88b.  ` + "`" + `88b  d88'  888   888   888 . 
  o888o o888o  ` + "`" + `Y8bd8P'   888bod8P'   "888" 
                          888               
                         o888o              
								 

k0pt - A command-line tool for managing Kubernetes clusters and namespaces with ease.
Author - Your Name

Usage:
  k0pt [command] [flags]

Supported Flags: namespace, clustername

Available Commands:
  use-namespace           Use a specific namespace
  delete-namespace        Delete a specific namespace
  start-cluster           Start a specific cluster
  stop-cluster            Stop a specific cluster
  get-cluster-state       Get the status of a specific cluster
  get-admin-credentials   Get admin credentials for a specific cluster
  analyze-resource-usage  Analyze resource usage across the cluster
  calculate-cost-savings  Calculate potential cost savings
  optimize-resources      Optimize resource allocation across the cluster
  version                 Get the version of k0pt

NOTE - Use "&" when using stop and start commands to run them in the background if needed.

Flags:
  -h, --help   help for k0pt

Use "k0pt [command] --help" for more information about a command.`)
}
