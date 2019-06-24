<div align="center">
  <img src="./assets/logo.png" width="20%" />
  <br><br>

  Another terminal based activity monitor for Kubernetes.

  <img src="./assets/screenshot1.png" width="100%" />
  <br><br>
  <img src="./assets/screenshot2.png" width="100%" />
</div>

## Installation

See [https://github.com/ricoberger/kubetop/releases](https://github.com/ricoberger/kubetop/releases) for the latest release.

```sh
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
wget https://github.com/ricoberger/kubetop/releases/download/v1.0.0/kubetop-$GOOS-$GOARCH
sudo install -m 755 kubetop-$GOOS-$GOARCH /usr/local/bin/kubetop
```

## Usage

kubetop has two entrypoints. The first one is the `pods` view, which shows the ressources of all running pods in the cluster. The second one is the `nodes` view which shows the ressources of all running nodes in the cluster. By selecting a node in the nodes view you get an overview of all running pods on this node. When you select a pod you get some details about this pod, like events and logs.

```
Display Resource (CPU/Memory/Storage) usage of pods

Usage:
  kubetop [flags]
  kubetop [command]

Available Commands:
  help        Help about any command
  nodes       Display Resource (CPU/Memory/Storage) usage of nodes
  version     Print version information for kubetop

Flags:
  -h, --help                help for kubetop
      --kubeconfig string   Path to the kubeconfig file to use for CLI requests
  -n, --namespace string    If present, the namespace scope for this CLI request

Use "kubetop [command] --help" for more information about a command.
```

The following keys can be used for the navigation in kubetop.

| Key | Action |
| --- | ------ |
| `q`, `<C-c>` | Quit kubetop |
| `k`, `<Up>`, `<MouseWheelUp>` | Scroll up (pods, nodes, logs) |
| `j`, `<Down>`, `<MouseWheelDown>` | Scroll down (pods, nodes, logs) |
| `<Tab>` | Select next container |
| `p` | Pause updating data |
| `<Enter>` | Select (pod, node, sortorder, filter) |
| `<Escape>` | Cancle (pod details, sortorder, filter) |
|  `<F1>` | Select sortorder |
|  `<F2>` | Select namespace filter |
|  `<F3>` | Select node filter |
|  `<F4>` | Select status filter |

## Dependencies

- [gotop](https://github.com/cjbassi/gotop): A terminal based graphical activity monitor inspired by gtop and vtop
- [cobra](https://github.com/spf13/cobra): A Commander for modern Go CLI interactions
- [termui](https://github.com/gizak/termui): A cross-platform and fully-customizable terminal dashboard and widget library
- [kubernetes](https://github.com/kubernetes/kubernetes): Production-Grade Container Scheduling and Management
