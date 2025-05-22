# Drax-the-OS



## colima setup for mac 
https://www.opencredo.com/blogs/building-the-best-kubernetes-test-cluster-on-macos

### prerequisit
- Homebrew
- Kubectl CLI client

### steps
- Docker CLI client
```
brew install docker
```
- Colima container runtime
```
brew install colima
colima start --network-address
colima list
```
- Kind
```
brew install kind
```
Save the following to a file called `kind-config.yaml:`
```
# 1 control plane node and 2 worker nodes
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name:  kind-multi-node
nodes:
- role: control-plane
- role: worker
- role: worker
```
Create a Kind cluster with config
```
kind create cluster --config=kind-config.yaml
```
- Configure Networking
```
export colima_host_ip=$(ifconfig bridge100 | grep "inet " | cut -d' ' -f2)
echo $colima_host_ip
export colima_vm_ip=$(colima list | grep docker | awk '{print $8}')
echo $colima_vm_ip
export colima_kind_cidr=$(docker network inspect -f '{{.IPAM.Config}}' kind | cut -d'{' -f2 | cut -d' ' -f1)
echo $colima_kind_cidr
export colima_kind_cidr_short=$(docker network inspect -f '{{.IPAM.Config}}' kind | cut -d'{' -f2 | cut -d' ' -f1| cut -d '.' -f1-2)
echo $colima_kind_cidr_short
export colima_vm_iface=$(colima ssh -- ip -br address show to $colima_vm_ip | cut -d' ' -f1)
echo $colima_vm_iface
export colima_kind_iface=$(colima ssh -- ip -br address show to $colima_kind_cidr | cut -d' ' -f1)
echo $colima_kind_iface
sudo route -nv add -net $colima_kind_cidr_short $colima_vm_ip
```

```
ssh_cmd="sudo iptables -A FORWARD -s $colima_host_ip -d $colima_kind_cidr -i $colima_vm_iface -o $colima_kind_iface -p tcp -j ACCEPT"
echo $ssh_cmd
colima ssh -- $ssh_cmd
exit
```
- MetalLB
```
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.9/config/manifests/metallb-native.yaml
kubectl wait --namespace metallb-system \
             --for=condition=ready pod \
             --selector=app=metallb \
             --timeout=90s
```
Save the following to a file called `metallb-conf.yaml`, replacing the addresses to suit your kind_cidr:
```
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: example
  namespace: metallb-system
spec:
  addresses:
  - 172.18.255.200-172.18.255.250
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: empty
  namespace: metallb-system
```
‍‍Apply the manifest:
```
kubectl apply -f metallb-conf.yaml
```


# install minimal prometheus and alertmanger
```
helm install minimal-prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheusOperator.enabled=true \
  --set grafana.enabled=false \
  --set kubeStateMetrics.enabled=true \
  --set kubeStateMetrics.resources.requests.memory=32Mi \
  --set kubeStateMetrics.resources.requests.cpu=50m \
  --set nodeExporter.enabled=false \
  --set alertmanager.enabled=true \
  --set prometheus.prometheusSpec.resources.requests.memory=256Mi \
  --set prometheus.prometheusSpec.resources.requests.cpu=100m
```

```
kubectl port-forward svc/minimal-prometheus-kube-pr-prometheus -n monitoring 9090:9090

kubectl port-forward svc/minimal-prometheus-kube-pr-alertmanager -n monitoring 9093:9093 
```
Delete unnecessary alert rules:
```
kubectl delete prometheusrules -n monitoring --all  
```
To uninstall everything:
```
helm uninstall minimal-prometheus -n monitoring
kubectl delete ns monitoring
```