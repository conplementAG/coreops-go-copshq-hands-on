# CoreOps Go copshq Hands-on Workshop

### Prerequisites

Access to Conplement Network (https://wiki.conplement.de/x/vQQWAQ)
Direct Access enabled (Windows only?)
OR
VPN with force tunneling (Windws, Linux, Mac)

---

Docker (https://docs.docker.com/get-docker/)
or Podman (https://podman.io/docs/installation)

Workshop resources are primarily provided for Docker. You can customize them as needed.

| Tool      | Tested Version | Installation Link                                                                                             | Note                                                                                                                                 |
| --------- | -------------- | ------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| Go        | 1.24.2         | [Go Installation](https://go.dev/doc/install)                                                                 |
| kubectl   | 1.32.3         | [kubectl Installation](https://kubernetes.io/docs/tasks/tools/)                                               |
| kubelogin | 0.2.7          | [kubelogin Installation](https://azure.github.io/kubelogin/install.html)                                      |
| copsctl   | 0.14.0         | [copsctl Installation](https://github.com/conplementAG/copsctl#installation)                                  |
| Helm      | 3.17.3         | [Helm Installation](https://helm.sh/docs/intro/quickstart/#install-helm)                                      |
| Azure CLI | 2.71.0         | [Azure CLI Installation](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) |
| Terraform | 1.11.4         | [Terraform Installation](https://developer.hashicorp.com/terraform/install)                                   |
| SOPS      | 3.10.2         | [SOPS Installation](https://github.com/getsops/sops/releases)                                                 | system variable to safely use VS code: https://stackoverflow.com/a/79523950 `EDITOR="code --wait --new-window --disable-extensions"` |

## Workshop

### 00_hello_world

```
go run main.go
go build main.go
./main.exe
```

### 01_commands

```
go run . hello
go run . hello welcome -n Conplement
```

### 02_executor

```
go run . executor
```

### 03_helm

Build and push app docker image first

> **Important:** To ensure we have no naming collisions everyone should use a unique abbriviation for their resources. For example the first letters of your name. This will be used for the container repository, cluster namespace and hostname for the ingress.

cponeneucopsacr\.azurecr.io/azure-notes-**>your-namespace<**:latest

```
cd app
docker build . -t azure-notes
docker tag azure-notes cponeneucopsacr.azurecr.io/azure-notes-**>your-namespace<**:latest
```

Login to Azure Container Registry used by Cluster

```
az acr login -n cponeneucopsacr
```

Push to Azure Container Registry

```
docker push cponeneucopsacr.azurecr.io/azure-notes-**>your-namespace<**:latest
```

Deploy

```
go run . deploy -c '<connection-string>'
```

#### HELM UPGRADE FAILED: another operation (install/upgrade/rollback) is in progress

helm ls -a -n ws-**>your-namespace<**
helm history cp-notes -n ws-**>your-namespace<**
helm uninstall cp-notes -n ws-**>your-namespace<**
helm rollback cp-notes -n ws-**>your-namespace<**
