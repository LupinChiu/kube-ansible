[![Build Status](https://travis-ci.org/kairen/kube-ansible.svg?branch=master)](https://travis-ci.org/kairen/kube-ansible)
# Kubernetes Ansible
A collection of playbooks for deploying/managing/upgrading a Kubernetes cluster onto machines, they are fully automated command to bring up a Kubernetes cluster on bare-metal or VMs.

[![asciicast](https://asciinema.org/a/fDjMx3fTZX9SZktqEdTtWwZwi.png)](https://asciinema.org/a/fDjMx3fTZX9SZktqEdTtWwZwi?speed=2)

Feature list:
- [x] Support Kubernetes v1.10.0+.
- [x] Highly available Kubernetes cluster.
- [x] Full of the binaries installation.
- [x] Kubernetes addons:
  - [x] Promethues Monitoring.
  - [x] EFK Logging.
  - [x] Metrics Server.
  - [x] NGINX Ingress Controller.
  - [x] Kubernetes Dashboard.
- [x] Support container network:
  - [x] Calico.
  - [x] Flannel.
- [x] Support container runtime:
  - [x] Docker.
  - [x] NVIDIA-Docker.(Require NVIDIA driver and CUDA 9.0+)
  - [x] Containerd.
  - [ ] CRI-O.

## Quick Start
In this section you will deploy a cluster via vagrant.

Prerequisites:
* Ansible version: *v2.5 (or newer)*.
* [Vagrant](https://www.vagrantup.com/downloads.html): >= 2.0.0.
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads): >= 5.0.0.
* Mac OS X need to install `sshpass` tool.

```sh
$ brew install http://git.io/sshpass.rb
```

The getting started guide will use Vagrant with VirtualBox to deploy a Kubernetes cluster onto virtual machines. You can deploy the cluster with a single command:
```sh
$ ./hack/setup-vms
Cluster Size: 1 master, 2 worker.
  VM Size: 1 vCPU, 2048 MB
  VM Info: ubuntu16, virtualbox
  CNI binding iface: eth1
Start to deploy?(y):
```
> * You also can use `sudo ./hack/setup-vms -p libvirt -i eth1` command to deploy the cluster onto KVM.

If you want to access API you need to create RBAC object define the permission of role. For example using `cluster-admin` role:
```sh
$ kubectl create clusterrolebinding open-api --clusterrole=cluster-admin --user=system:anonymous
```

Login the addon's dashboard:
- Dashboard: [https://API_SERVER:8443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/](https://API_SERVER:8443/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/)
- Logging: [https://API_SERVER:8443/api/v1/namespaces/kube-system/services/kibana-logging/proxy/](https://API_SERVER:8443/api/v1/namespaces/kube-system/services/kibana-logging/proxy/)

As of release 1.7 Dashboard no longer has full admin privileges granted by default, so you need to create a token to access the resources:
```sh
$ kubectl -n kube-system create sa dashboard
$ kubectl create clusterrolebinding dashboard --clusterrole cluster-admin --serviceaccount=kube-system:dashboard
$ kubectl -n kube-system get sa dashboard -o yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: 2017-11-27T17:06:41Z
  name: dashboard
  namespace: kube-system
  resourceVersion: "69076"
  selfLink: /api/v1/namespaces/kube-system/serviceaccounts/dashboard
  uid: 56b880bf-d395-11e7-9528-448a5ba4bd34
secrets:
- name: dashboard-token-vg52j

$ kubectl -n kube-system describe secrets dashboard-token-vg52j
...
token:      eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkYXNoYm9hcmQtdG9rZW4tdmc1MmoiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGFzaGJvYXJkIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiNTZiODgwYmYtZDM5NS0xMWU3LTk1MjgtNDQ4YTViYTRiZDM0Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Omt1YmUtc3lzdGVtOmRhc2hib2FyZCJ9.bVRECfNS4NDmWAFWxGbAi1n9SfQ-TMNafPtF70pbp9Kun9RbC3BNR5NjTEuKjwt8nqZ6k3r09UKJ4dpo2lHtr2RTNAfEsoEGtoMlW8X9lg70ccPB0M1KJiz3c7-gpDUaQRIMNwz42db7Q1dN7HLieD6I4lFsHgk9NPUIVKqJ0p6PNTp99pBwvpvnKX72NIiIvgRwC2cnFr3R6WdUEsuVfuWGdF-jXyc6lS7_kOiXp2yh6Ym_YYIr3SsjYK7XUIPHrBqWjF-KXO_AL3J8J_UebtWSGomYvuXXbbAUefbOK4qopqQ6FzRXQs00KrKa8sfqrKMm_x71Kyqq6RbFECsHPA
```
> Copy and paste the `token` to dashboard.

## Manual deployment
In this section you will manually deploy a cluster on your machines.

Prerequisites:
* Ansible version: *v2.5 (or newer)*.
* *Linux distributions*: Ubuntu 16+/Debian/CentOS 7.x.
* All Master/Node should have password-less access from `deploy` node.

For machine example:

| IP Address      |   Role           |   CPU    |   Memory   |
|-----------------|------------------|----------|------------|
| 172.16.35.9     | vip              |    -     |     -      |
| 172.16.35.10    | k8s-m1           |    4     |     8G     |
| 172.16.35.11    | k8s-n1           |    4     |     8G     |
| 172.16.35.12    | k8s-n2           |    4     |     8G     |
| 172.16.35.13    | k8s-n3           |    4     |     8G     |

Add the machine info gathered above into a file called `inventory/hosts.ini`. For inventory example:
```
[etcds]
k8s-m1
k8s-n[1:2]

[masters]
k8s-m1
k8s-n1

[nodes]
k8s-n[1:3]

[kube-cluster:children]
masters
nodes
```
* 失敗及待辨事項：
  1. etcd 使用的 ca.pem 經初步測試，使用自定的 中繼憑證 會失敗，原因仍然在排查中，暫時使用 舊有方式 自簽 rootCA 的方式讓 etcd 能通過
  初步判定可能是憑證的 link list 無法辨別，錯誤訊息如下：
```
  May 11 19:07:06 k8snode43.udream.local etcd[13795]: health check for peer a9d0c318e1272ae1 could not connect: dial tcp 172.29.19.48:2380: connect: connection refused (prober "ROUND_TRIPPER_RAFT_MESSAGE")
May 11 19:07:06 k8snode43.udream.local etcd[13795]: health check for peer a9d0c318e1272ae1 could not connect: dial tcp 172.29.19.48:2380: connect: connection refused (prober "ROUND_TRIPPER_SNAPSHOT")
May 11 19:07:06 k8snode43.udream.local etcd[13795]: health check for peer f9b6ca26f2b010b7 could not connect: x509: certificate signed by unknown authority (possibly because of "x509: invalid signature: parent certifi
May 11 19:07:06 k8snode43.udream.local etcd[13795]: health check for peer f9b6ca26f2b010b7 could not connect: x509: certificate signed by unknown authority (possibly because of "x509: invalid signature: parent certifi
May 11 19:07:06 k8snode43.udream.local etcd[13795]: rejected connection from "172.29.19.49:56828" (error "remote error: tls: bad certificate", ServerName "")
May 11 19:07:06 k8snode43.udream.local etcd[13795]: rejected connection from "172.29.19.49:56830" (error "remote error: tls: bad certificate", ServerName "")
May 11 19:07:06 k8snode43.udream.local etcd[13795]: rejected connection from "172.29.19.48:46207" (error "remote error: tls: bad certificate", ServerName "")
May 11 19:07:06 k8snode43.udream.local etcd[13795]: rejected connection from "172.29.19.48:46206" (error "remote error: tls: bad certificate", ServerName "")
May 11 19:07:06 k8snode43.udream.local etcd[13795]: rejected connection from "172.29.19.49:56838" (error "remote error: tls: bad certificate", ServerName "")


這裡是上面被中斷的錯誤訊息
May 11 19:02:49 k8snode49 etcd: health check for peer a9d0c318e1272ae1 could not connect: x509: certificate signed by unknown authority (possibly because of "x509: invalid signature: parent certificate cannot sign this kind of certificate" while trying to verify candidate authority certificate "etcd") (prober "ROUND_TRIPPER_RAFT_MESSAGE")
May 11 19:02:49 k8snode49 etcd: health check for peer a9d0c318e1272ae1 could not connect: x509: certificate signed by unknown authority (possibly because of "x509: invalid signature: parent certificate cannot sign this kind of certificate" while trying to verify candidate authority certificate "etcd") (prober "ROUND_TRIPPER_SNAPSHOT")
May 11 19:02:49 k8snode49 etcd: health check for peer d970fd5472d32b4e could not connect: x509: certificate signed by unknown authority (possibly because of "x509: invalid signature: parent certificate cannot sign this kind of certificate" while trying to verify candidate authority certificate "etcd") (prober "ROUND_TRIPPER_RAFT_MESSAGE")
May 11 19:02:49 k8snode49 etcd: health check for peer d970fd5472d32b4e could not connect: x509: certificate signed by unknown authority (possibly because of "x509: invalid signature: parent certificate cannot sign this kind of certificate" while trying to verify candidate authority certificate "etcd") (prober "ROUND_TRIPPER_SNAPSHOT")
```
  1. K8S-IM 也失敗，目前只要是使用自已的 中繼憑證 通通都宣告失敗，由於沒時間除錯，只能先全部復原，使用原本的自建root CA方式
  1. haproxy 需要新增一個 dns server 以利辨別 haproxy 裡面設定的 server name 解析
     1. 目前是使用此方式：開始前要先把機器名及ip寫入所有機器的 /etc/hosts 裡
```
  echo -e "172.29.19.43 k8s-m1 k8s-n1\n172.29.19.48 k8s-m2 k8s-n2\n172.29.19.49 k8s-m3 k8s-n3" >> /etc/hosts

例：
172.29.19.43 k8s-m1 k8s-n1
172.29.19.48 k8s-m2 k8s-n2
172.29.19.49 k8s-m3 k8s-n3
```

* 要記得改以下的設定檔：
```
vip_address -> Keepalived virtual ip (在這裡主要是用在 ha kubernetes api server)
vip_interface -> virtual ip 網卡設定 example: eth0/eth1.....
encryption_token -> 已修改

  inventory/group_vars/all.yml
  inventory/hosts.ini.lupin-example 修改完後請存成 hosts.ini
  roles/downloads/package/defaults/main.yml
  roles/cluster-defaults/defaults/main.yml
  roles/k8s-setup/defaults/main.yml
    1. keepalived_password
    2. keepalived_priority  將 master 權限由 100 => 200 
    (查了一下keepalived裡 master 的設定應該是要比較高的值才對，原本是第一台設為 100，代表一開始的 api server 希望不是由第一台任職)
    1. haproxy_stats_user
    2. haproxy_stats_password
```

* 憑證相關需在執行前處理完成
  1. 需提前瞭解及使用 cfssl 簽出 中繼憑證 並進行相關測試
  1. 使用提前簽出的 中繼憑證 替換Ansible內部的 ca 憑證 (不需要，已修改對應 yml 檔)
  1. 改寫對應 roles/cert 內容 (已進行相關修改)

* 其他注意事項
  1. Kubernetes Audit 主要目的為記錄 kubernetes 的各種運行記錄，正式上線前需要詳細調整
  1. 此 ansible 設定可能使用 fluentd 來將 Audit 的日志記錄 存在磁碟上 -- 後續需追查
  1. --audit-log-path 設定記錄檔，看到這個設定，判斷他是存在 api server 的 pod 上
  1. 此Ansible 在 kubernetes Api server內的設定檔有將 /etc/ssl/certs 映射至本機一樣的位置
  1. featureGates - PodPriority 參數經追查 PodPriority 已在 1.14 版本穩定化，不需要額外設定了
  1. featureGates - DevicePlugins 參數從 1.10 版本開始進入 Beta ，仍然未 穩定化
  1. 10-kubelet.conf.j2 內部參數追查有如下限制：
     1. 在 1.18 仍然為 alpha 的參數如下：
        1. --network-plugin    # 限制 cri 為 docker 時使用
        2. --cni-conf-dir      # 限制 cri 為 docker 時使用
        3. --cni-bin-dir       # 限制 cri 為 docker 時使用
        4. --node-labels       # 在 1.15 版之後 此參數規則有變，無法支援 node-role.kubernetes.io/master 該命名規則
           1. 修正如下：--node-labels=node.kubernetes.io/lupin=''
     2. 已完全消失參數
     3. --allow-privileged
  1. 由於此 Ansible dns 解析方式採用 kubedns 需考慮改用 coredns 以利後續佈署至雲端正式集群 (未)
     1. 原因： coredns 別人的測試報告指出 coredns解析 外部域名 的時間約比 kubedns 快上三倍(12ms vs 41ms)，且 coredns 比 kubedns 更加靈活
     2. 雖不見得未來會需要使用外部域名的服務，保險起見還是先當成會使用
     3. https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/dns/coredns
  1. 需將 kubernetes 設定改為優先取用私人 docker repo
     1. 需新增下列 secret 以存取私人 docker registry
```
-- Kubernetes 集群使用 docker-registry 类型的 Secret 来通过容器仓库的身份验证，进而提取私有映像。

创建 Secret，命名为 regcred：
kubectl create secret docker-registry regcred --docker-server=${YOUR_DOCKER_REGISTRY} --docker-username=${USER} --docker-password=${PASSWORD} --docker-email=you@domain.com
```
  1. 在 role\cert\main.yml 最下面有一段是刪除所有 .csr 及憑證設定檔 .json 若要長久使用同一憑證，這裡就不太適合刪除

Set the variables in `group_vars/all.yml` to reflect you need options. For example:
```yml
# overide kubernetes version(default: 1.10.6)
kube_version: 1.11.2

# container runtime, supported: docker, nvidia-docker, containerd.
container_runtime: docker

# container network, supported: calico, flannel.
cni_enable: true
container_network: calico
cni_iface: ''

# highly available variables
vip_interface: ''
vip_address: 172.16.35.9

# etcd variables
etcd_iface: ''

# kubernetes extra addons variables
enable_dashboard: true
enable_logging: false
enable_monitoring: false
enable_ingress: false
enable_metric_server: true

# monitoring grafana user/password
monitoring_grafana_user: "admin"
monitoring_grafana_password: "p@ssw0rd"
```

### Deploy a Kubernetes cluster
If everything is ready, just run `cluster.yml` playbook to deploy the cluster:
```sh
$ ansible-playbook -i inventory/hosts.ini cluster.yml
```

And then run `addons.yml` to create addons:
```sh
$ ansible-playbook -i inventory/hosts.ini addons.yml
```

## Verify cluster
Verify that you have deployed the cluster, check the cluster as following commands:
```sh
$ kubectl -n kube-system get po,svc

NAME                                 READY     STATUS    RESTARTS   AGE       IP             NODE
po/haproxy-master1                   1/1       Running   0          2h        172.16.35.10   k8s-m1
...
```

### Reset cluster
Finally, if you want to clean the cluster and redeploy, you can reset the cluster by `reset-cluster.yml` playbook.:
```sh
$ ansible-playbook -i inventory/hosts.ini reset-cluster.yml
```

## Contributing
Pull requests are always welcome!!! I am always thrilled to receive pull requests.
