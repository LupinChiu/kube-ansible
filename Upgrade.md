# Upgrade k8s

在K8s中與主要版本相關的元件有 
- kube-apiserver 
- kube-controller-manager 
- kube-scheduler 
- kube-proxy 
- kubectl 
- kubelet 

以及
- etcd 
  
其中只有etcd有自己的版本號，其餘都是跟隨著k8s主要版本號更新 

## 更動項目
### kube-apiserver
更新`kube-apiserver`時在`roles/k8s-setup/templates/manifests/kube-apiserver.yml.j2`路徑下的設定有變動

需要加入以下設定
```
      - --service-account-key-file={{ sa_public_key }}
      - --service-account-signing-key-file={{ sa_private_key }}
      - --service-account-issuer=https://kubernetes.default.svc.cluster.local
```

### kube-proxy
在1.20版本以後不再支援此設定
```yaml
    featureGates:
      SupportIPVSProxyMode: true
```
(修改後bug激增，未解) 如需要使用IPVS將設為 `mode: ipvs` 即可 (修改後bug激增，未解)

### etcd
來源包有更動，有做修改
```yaml
    file: "etcd-v{{ base.etcd.version }}-linux-amd64.tar.gz"
```
改為
```yaml
    file: "etcd-v{{ base.etcd.version }}-linux-amd64.tar"
```

### kubectl & kubelet
kubelet版本不能比你所用的kubectl版本差超過一個小版本號。

例如，如果你的集群為 v1.18，那麼你可以使用 v1.17、v1.18、v1.19 的 kubectl，所有其他的组合都是不支援的。
 
 