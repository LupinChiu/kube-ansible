# Istio Traffic Management

主要以流量的入口與出口為主，即 Gateway 的設置；而 Istio 官方建議部署業務應用時，至少
配置
1. 使用 app 標籤表明應用身分
2. 使用 version 標籤標明應用版本
3. 建立 DestniationRule 
4. 建立 Virtual Service

## 目錄
* [Prerequisite](#Prerequisite)
* [DestniationRule](#DestniationRule)
* [VirtualService](#VirtualService)
	* [Canary Deployment( 金絲雀部署 )](#金絲雀部署)
* [Gateway](#Gateway)
	* [Ingress](#Ingress)
	* [Egress](#Egress)
* [Timeout](#Timeout)
* [Fault Injection](#故障植入測試)
	* [Delay](#Delay)
	* [Fault](#Fault)
* [Mirroring](#Mirroring)
* [Circuit Breaking](#鎔斷)
* [Demo](#Demo)
* [Note](#Note)

## Prerequisite
* 已部署 Istio
* 已開啟自動注入 Istio Proxy(sidecar)

## DestniationRule
### 與 k8s 的不同
在 k8s 中，client 需要使用不同服務入口才可以存取多個不同服務；而 <span class=red>Istio 可以只使用一個服務入口</span>，Istio 透過流量的特徵來完成對後端服務的選擇。
### 欄位說明
* host (required): 代表 k8s 的 service，或由一個 ServiceEntry 定義的外部服務；建議使用 FQDN。
* trafficPolicy (optional): 流量策略；DestniationRule Level 和 Subset Level 皆可定義，<span class=red>Subset Level 會 override DestniationRule Level</span>。
* subsets (optional) : 使用標籤選擇器來定義不同子集；即可以用版本標籤來區別流量策略。

## VirtualService
在沒有定義 VirtualService 的情況下，DestniationRule 的 subset 是沒有作用的；會依造 [kube-proxy](https://hackmd.io/@daemonbuu/S1lNp8Tcv) 的預設隨機行為進行存取。

<span class=red>VirtualService 負責對流量進行判別和轉發。</span>

### 欄位說明
* hosts : 一樣是針對 k8s 的 service，或由一個 ServiceEntry 定義的外部服務作服務；可以針對多個服務進行工作。
* 支援多種協定
	* http
	* tcp
	* tls

### 流量拆分和移轉
正式和測試
1. 部署兩個應用 v1、v2
``` yaml
apiVersion: v1
kind: Service
metadata:
  name: simple-login-svc
  labels:
    app: simple-login
spec:
  ports:
  - name: http
    port: 8033
  selector:
    app: simple-login
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple-login-v1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: simple-login
      version: v1
  template:
    metadata:
      labels:
        app: simple-login
        version: v1
    spec:
      containers:
      - name: simple-login
        image: docker.io/whitewalker0506/simple_login:1.0.0
        imagePullPolicy: IfNotPresent
        args: [ "10.0.1.101", "3308" ]
        ports:
        - containerPort: 8033
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple-login-v2
spec:
  replicas: 1
  selector:
    matchLabels:
      app: simple-login
      version: v2
  template:
    metadata:
      labels:
        app: simple-login
        version: v2
    spec:
      containers:
      - name: simple-login
        image: docker.io/whitewalker0506/simple_login:1.0.0
        imagePullPolicy: IfNotPresent
        args: [ "10.0.1.101", "3308" ]
        ports:
        - containerPort: 8033
---
```
2. 部署 DestniationRule
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: simple-login-dr
spec:
  host: simple-login-svc.default.svc.cluster.local
  trafficPolicy:
    loadBalancer:
      simple: LEAST_CONN
  subsets:
  - name: prodversion
    labels:
      version: v1
    trafficPolicy:
      loadBalancer:
        simple: ROUND_ROBIN
  - name: testversion
    labels:
      version: v2
    trafficPolicy:
      loadBalancer:
        simple: RANDOM
```

3. 部署 VirtualService
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: simple-login-vsvc
spec:
  hosts:
  - simple-login-svc.default.svc.cluster.local
  gateways:
  - simple-login-gateway
  http:
  - name: "simple-login-v1-routes"
    match:
    - uri: 
        exact: /index
    - uri: 
        exact: /internal
    - uri: 
        exact: /v1/login
    - uri: 
        exact: /v1/logout
    route:
    - destination:
        host: simple-login-svc.default.svc.cluster.local
        port:
          number: 8033
        subset: v1
  - name: "simple-login-v2-routes"
    match:
    - uri: 
        exact: /index
    - uri: 
        exact: /internal
    - uri: 
        exact: /v1/login
    - uri: 
        exact: /v1/logout
    route:
    - destination:
        host: simple-login-svc.default.svc.cluster.local
        port:
          number: 8033
        subset: v2
```
4. 部署 Ingress
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: simple-login-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - simple-login-svc.default.svc.cluster.local
```
### 金絲雀部署
可以透過不同使用者來拆分正式和測試

## Gateway
即 mesh 的 edge proxy，負責 traffic 的入口與出口
![](https://i.imgur.com/RugVST8.png)

### Ingress
入口，負責外部可以存取 mesh 內的服務。
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: simple-login-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
```

### Egress
出口，負責使 mesh 內的服務可以存取外部服務 (不在 mesh 內)；建議使用 ServiceEntry 的方式訪問外部服務，優點是<span class=red>不會丟失流量監控和控制特性</span>。

* 固定 IP 沒有提供 FQDN 的方式
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: mysql
spec:
  hosts:
  - sample.db
  ports:
  - number: 3308
    name: tcp
    protocol: TCP
  resolution: STATIC
  location: MESH_EXTERNAL
  endpoints:
  - address: 10.0.1.101
```

* 有提供 FQND
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: cnn
spec:
  hosts:
  - edition.cnn.com
  ports:
  - number: 80
    name: http-port
    protocol: HTTP
  - number: 443
    name: https
    protocol: HTTPS
  resolution: DNS
```

## Demo
[demo project](https://github.com/LupinChiu/kube-ansible/tree/master/istio_demo/egress)
1. 成功發布 demo 後，可在 browser 輸入 10.0.1.179/index
![](https://i.imgur.com/LFhVBwy.png)
2. User name :1234；Password:1234
![](https://i.imgur.com/5zPtv6l.png)
3. traffic graph
![](https://i.imgur.com/GIoKzuQ.png)


## Note
1. 目前在從零到有部署 istio 時，需要先將 istio-injection 設為 disabled (原因為何 ? 待釐清，有可能和憑證有關)
e.g.
``` yaml
apiVersion: v1
kind: Namespace
metadata:
  name: istio-system
  labels:
    istio-injection: enabled
```
