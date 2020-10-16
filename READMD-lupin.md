
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
    3. haproxy_have_dns_server 如果有對應的 dns server 可以試著打開這個功能，否則會直接寫入 /etc/hosts 下，直接加入機器Ip 對應
```
* 除錯
  1. vip 的設定如果出問題，需要追查 keepalived 而不是 haproxy (設定檔已修正)
  2. haproxy 再次安裝仍然有出錯 (設定檔已修正)

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


# make all chain
```
cd /etc/kubernetes/pki/
cat ca.pem intermediate/IM-CA-bundle.crt.pem > ca-chain.pem
cat front-proxy-ca.pem intermediate/IM-CA-bundle.crt.pem > front-proxy-ca-chain.pem
cat ./etcd/etcd-ca.pem ./intermediate/IM-CA-bundle.crt.pem > ./etcd/etcd-ca-chain.pem

cat apiserver.pem ca-chain.pem > apiserver-chain.pem
cat admin.pem ca-chain.pem > admin-chain.pem
cat controller-manager.pem ca-chain.pem > controller-manager-chain.pem
cat kubelet.pem ca-chain.pem > kubelet-chain.pem
cat scheduler.pem ca-chain.pem > scheduler-chain.pem

cat front-proxy-client.pem front-proxy-ca-chain.pem > front-proxy-client-chain.pem

cat ./etcd/etcd.pem ./etcd/etcd-ca-chain.pem > ./etcd/etcd-chain.pem
```

# 取得token
```sh
簡易指令：
kubectl -n kubernetes-dashboard get sa admin-user -o json | grep \"secrets\" -A 5 | grep name | awk '{print $2} ' | xargs kubectl -n kubernetes-dashboard describe secrets
```
