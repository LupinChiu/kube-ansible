# Istio 相關

## 目錄
* [Prerequisite](#Prerequisite)
* [Istio-Ingressgateway](Istio-Ingressgateway)
    * [Add Private Docker Registry After Cluster Was Created](#Private-Docker-Registry)
* [Namespace](#Namespace)
* [Command](#Command)
* [Gateway](#Gateway)
* [Demo](#Demo)
* [Note](#Note)

## Prerequisite
* 添加 enable_istio 變數至 all.yml 
* deployment istio yaml files
* 添加 istio role 至 addons 
* 添加 insecure_registries 設定 (option)
    * <span class=red>Note: [創建集群後才添加](#Private-Docker-Registry)</span>

e.g.
path: <span class=orange>"roles\cluster-default\defaults\main.yml"</span>
``` yaml
insecure_registries:
- "grd-dev.urad.com.tw"

```



## Istio-Ingressgateway

* 啟用 autoInject 
* 啟用 enableNamespacesByDefault
建議直接在 istio-iop.yml.j2 搜尋 <span class=red>sidecarInjectorWebhook</span>

e.g.
``` yaml
sidecarInjectorWebhook:
      enableNamespacesByDefault: true
      rewriteAppHTTPProbe: true
      injectLabel: istio-injection
      objectSelector:
        enabled: false
        autoInject: true
```
* exteranl ip and expose tcp port

e.g.
``` yaml
ingressGateways:
    - name: istio-ingressgateway
      enabled: true
      k8s:
        env:
          - name: ISTIO_META_ROUTER_MODE
            value: "sni-dnat"
        service:
          externalIPs:
            - 10.0.1.179
          ports:
            - port: 15021
              targetPort: 15021
              name: status-port
            - port: 80
              targetPort: 8080
              name: http2
            - port: 443
              targetPort: 8443
              name: https
            - port: 15443
              targetPort: 15443
              name: tls
            - port: 31400
              targetPort: 31400
              name: tcp
```

## Namespace
* 啟用 istio-injection

e.g.
``` yaml
apiVersion: v1
kind: Namespace
metadata:
  name: istio-system
  labels:
    istio-injection: enabled
```


## Private-Docker-Registry
> **<span class=red>強烈建議</span> : 不要再創建 cluster 後才作添加私有庫**
1. 修改<span class=red>各個 master node</span> 下的 docker daemon.json

e.g.
path: <span class=orange>/etc/docker/daemon.json</span>
``` json
{
    "insecure-registries": ["grd-dev.urad.com.tw"],
    "debug": false
}
```

2. reload daemon
``` 
systemctl daemon-reload
```
3. 重啟 docker.sock
``` 
systemctl start docker.sock
```
4. 重啟 docker.service
``` 
systemctl start docker.service
```
5. 留意 kube-controller-manager 是否有發生問題
> 由於 containerd 為 kube-controller-manager 父 pid ，而 containerd 又與 docker.scok 有關，故重啟 docker.sock 可能會導致 kube-controller-manager 有不預期的錯誤。


## Command
* 查看哪些 Namespace 有啟用 istio-injection
``` cmd
kubectl get namespace -L istio-injection
```
* 啟用某個 namespace 的 istio-injection
``` cmd
kubectl label namespace {NamespaceName} istio-injection=enabled
```

## Gateway
* 設定 Service
``` yaml
apiVersion: v1
kind: Service
metadata:
  name: tcp-echo-server
  labels:
    app: tcp-echo-server
spec:
  ports:
  - name: tcp
    port: 8020
  # Port 9002 is omitted intentionally for testing the pass through filter chain.
  selector:
    app: echo-server
```
* 設定 Gateway
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: tcp-echo-gateway-server
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 31400
      name: tcp
      protocol: TCP
    hosts:
    - "*"
```
* 設定 VirtualService
``` yaml 
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: tcp-echo
spec:
  hosts:
  - "*"
  gateways:
  - tcp-echo-gateway-server
  tcp:
  - match:
    - port: 31400
    route:
    - destination:
        host: tcp-echo-server
        port:
          number: 8020
```


## Demo
[echo sevice demo](https://github.com/LupinChiu/kube-ansible/tree/master/istio_demo/echo_service_for_k8s)

```sequence
client -> gateway: send request
gateway -> virtual service:轉導到相應 VS
virtual service -> service:轉到 match service
service -> server(Pod): 分配給可執行的 pod
server(Pod) -> client: reply result
```

成功發布 demo 後，可在 kiali 看到此圖
![](https://i.imgur.com/TVfch1q.png)

## Note
1. golang dial tcp 目前有問題 (待釐清)
2. 留意其他台 Master Node 是否有設置私有 repository ，否則 replicas 可能造成 pull image 失敗
