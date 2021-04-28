# Rook Ceph

![rook](https://blog.kasten.io/hubfs/Blog%20Images/Backup%20and%20DR%20using%20Rook%201.4%20and%20Ceph%20CSI%203.0/rook_ceph_logos-social-01.png "rook")

## Source code
https://github.com/rook/rook.git

## 環境需求與安裝
- 完整 kubernetes cluster ，至少有三個 worker node
- 至少在一個 worker node 裡，有一個未使用的 disk ，提供為 ceph 使用

### Rook 所能拿來提供為 storage 使用的來源，只有下列三種
- Raw devices (no partitions or formatted filesystems)
- Raw partitions (no formatted filesystem)
- PVs available from a storage class in block mode


## Ansible安裝
1. 開啟`enable_ceph: true`
2. `ansible-playbook -i inventory/ceph_hosts.ini add-nodes.yml`
3. 等待安裝完成(如果卡死，有可能需要手動重啟VM)
4. `kubectl apply -f filesystem.yaml storageclass.yaml`
5. 完成


## 快速手動安裝部屬流程
```shell
git clone --single-branch --branch release-1.6 https://github.com/rook/rook.git
```
```shell
cd rook/cluster/examples/kubernetes/ceph/
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
kubectl create -f cluster.yaml
// 等待cluster完成 約10-20分鐘
kubectl create -f filesystem.yaml
kubectl create -f dashboard-ingress-https.yaml

cd rook/cluster/examples/kubernetes/ceph/csi/cephfs/
kubectl create -f storageclass.yaml
```

## 測試驗證
```shell
cd rook/cluster/examples/kubernetes/ceph/csi/cephfs/
kubectl create -f pvc.yaml
kubectl create -f pod.yaml
```

## 注意事項與bug
- 要確認安裝好`operator`才可以安裝`cluster`
- 安裝`cluster`會耗費20-30分鐘，如遇上api server掛掉、其中一台極為忙碌，可以嘗試重啟host
- 大部分相關設定可以在`cluster.yaml`中設定
- dashboard可以用ingress來引導出來，預設網址為 rook-ceph.example.com
- dashboard
    - 帳號: admin
    - 密碼:

```shell
kubectl -n rook-ceph get secret rook-ceph-dashboard-password -o jsonpath='{.data.password}'  |  base64 --decode
```

- 有重新安裝可能會抓不到硬碟，要記得把舊資料清除
    - 清除設定`rm -rf /var/lib/rook`，這在`cluster.yaml`中可修改
    - 刷硬碟

```
#!/usr/bin/env bash
DISK="/dev/vda"

# Zap the disk to a fresh, usable state (zap-all is important, b/c MBR has to be clean)

# You will have to run this step for all disks.
sgdisk --zap-all $DISK

# Clean hdds with dd
dd if=/dev/zero of="$DISK" bs=1M count=100 oflag=direct,dsync

# Clean disks such as ssd with blkdiscard instead of dd
blkdiscard $DISK

# These steps only have to be run once on each node
# If rook sets up osds using ceph-volume, teardown leaves some devices mapped that lock the disks.
ls /dev/mapper/ceph-* | xargs -I% -- dmsetup remove %

# ceph-volume setup can leave ceph-<UUID> directories in /dev and /dev/mapper (unnecessary clutter)
rm -rf /dev/ceph-*
rm -rf /dev/mapper/ceph--*
```




