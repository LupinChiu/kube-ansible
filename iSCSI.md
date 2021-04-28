# iSCSI (Internet Small Computer System Interface)

## iSCSI Server設置
每個裝置不一樣，NAS有自己的設定，這裡示範安裝在CentOS


https://www.server-world.info/en/note?os=Fedora_21&p=iscsi


利用`targetcli`可以將設定完成
## iSCSI Client端
https://www.server-world.info/en/note?os=Fedora_21&p=iscsi&f=2


可以用這個方法來做SERVER驗證

### K8s中的設置
https://github.com/kubernetes/examples/tree/master/volumes/iscsi


https://blog.csdn.net/hongxiaolu/article/details/113711593


利用這個範例來實現

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: chap-secret
type: "kubernetes.io/iscsi-chap"
data:
  discovery.sendtargets.auth.username: dXNlcm5hbWU=
  discovery.sendtargets.auth.password: cGFzc3dvcmQ=
  discovery.sendtargets.auth.username_in: dXNlcm5hbWU=
  discovery.sendtargets.auth.password_in: cGFzc3dvcmQ=
  node.session.auth.username: dXNlcm5hbWU=
  node.session.auth.password: cGFzc3dvcmQ=
  node.session.auth.username_in: dXNlcm5hbWU=
  node.session.auth.password_in: cGFzc3dvcmQ=
```
(因為Secret所以要求base64編碼)


所有的帳號密碼在server端都需要有相對應的設置(`targetcli`中做設定)
```shell
cd /iscsi
set discovery_auth enable=1
set discovery_auth userid=username
set discovery_auth password=password
set discovery_auth mutual_userid=username
set discovery_auth mutual_password=password

cd /iscsi/iqn.2021-04.k8s.example:storage/tpg1/acls/iqn.2021-04.k8s.example:node
set auth userid=username
set auth password=password
set auth mutual_userid=username
set auth mutual_password=password
```

Node 需要安裝 
```
yum install iscsi-initiator-utils -y
```

最後範例yaml
我沒有去設定`mutual`跟`discovery`認證
```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: iscsipd
spec:
  containers:
    - name: iscsipd-ro
      image: nginx
      volumeMounts:
        - mountPath: "/mnt"
          name: iscsivol
  volumes:
    - name: iscsivol
      iscsi:
        targetPortal: 10.0.1.89
        iqn: iqn.2021-04.k8s.srv:storage.target00
        lun: 0
        fsType: ext4
        readOnly: false
        initiatorName: iqn.2021-04.k8s.srv:node
        chapAuthSession: true
        chapAuthDiscovery: false
        secretRef:
          name: chap-secret
```

## Troubleshooting
- 遇到登入驗證問題可能是帳號密碼不對
- K8s的Node只會使用到`iscsi-initiator-utils`工具，不會吃到設定黨
- 重開機




## 詳細介紹
http://linux.vbird.org/linux_server/0460iscsi.php