# CRI-O
(Container Runtime Interface) using OCI (Open Container Initiative)

## 啟用模組
```shell
sudo modprobe overlay
sudo modprobe br_netfilter
```

## Source code
https://storage.googleapis.com/k8s-conform-cri-o/artifacts/crio-v1.19.1.tar.gz


https://github.com/cri-o/cri-o/releases/download/v1.19.1/crio-v1.19.1.tar.gz

## Install
```shell
sudo yum install make -y
sudo make install
```
### 注意事項
不知道為什麼`install`不會建立新資料夾，可能是因為CentOS 7(?) 

所以我這邊有手動建立需要的資料夾
```shell
mkdir -p /opt/cni/bin
mkdir -p /usr/local/share/oci-umount/oci-umount.d
mkdir -p /usr/local/lib/systemd/system
mkdir -p /etc/crio
```


## 加入群組
需要在 `/etc/crio/crio.conf` 或是 `/etc/crio/crio.conf.d/02-cgroup-manager.conf` 
中設定crio使用群組
```shell
[crio.runtime]
conmon_cgroup = "pod"
cgroup_manager = "cgroupfs"
```

## 配置Registries庫
放置設定黨`/etc/containers/registries.conf`
```shell
[registries.search]
registries = ["docker.io","grd.urad.com.tw"]

[registries.insecure]
registries = []

[registries.block]
registries = []
```

## kubelet 設定
需要在`kubelet.service`加入
```shell
Wants=docker.socket crio.service
```
You need to add following parameters to `KUBELET_ARGS`:

- `--container-runtime=remote` - Use remote runtime with provided socket.
- `--container-runtime-endpoint=unix:///var/run/crio/crio.sock` - Socket for remote runtime (default crio socket localization).
- `--runtime-request-timeout=10m` - Optional but useful. Some requests, especially pulling huge images, may take longer than default (2 minutes) and will cause an error.

Sample:
```shell
Environment="CRI_RUNTIME_ARGS=--container-runtime=remote --container-runtime-endpoint=unix://{{ runtime_sockets['' + container_runtime] }}"
```

## CNI (Container Network Interface)
預設會抓取`/etc/cni/net.d`中第一個檔案


## Start CRI-O:
```shell
sudo systemctl daemon-reload
sudo systemctl enable crio --now
sudo systemctl start --now crio
```

## 查看狀態
```shell
sudo systemctl status crio
```

## 移除
先停止後刪除container
```shell
crictl stop $(crictl ps -a -q)
crictl rm $(crictl ps -a -q)
```

```shell
sudo make uninstall
```
