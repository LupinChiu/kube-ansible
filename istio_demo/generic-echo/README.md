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
* hosts : 一樣是針對 k8s 的 service，或由一個 ServiceEntry 定義的外部服務作服務；可以針對多個服務進行工作。 (白話其實是提供給 Client 呼叫的位置，與 Gateway 的 hosts 匹配)
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
這邊以使用來源版本標籤的方式來模擬金絲雀部署；v1 Client 的請求會導向 v1 Server，而其他版本的請求會導向 v2 Server 。

e.g.
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: generic-echo-vsvc
spec:
  hosts:
  - "generic-echo-server-svc.default.svc.cluster.local"
  gateways:
  - generic-echo-server-gateway
  http:
  - match:
    - sourceLabels:
        app: generic-echo-client
        version : v1
    route:
    - destination:
        host: generic-echo-server-svc.default.svc.cluster.local
        subset: v1
  - route:
    - destination:
        host: generic-echo-server-svc.default.svc.cluster.local
        subset: v2
```

## Gateway
即 mesh 的 edge proxy，負責 traffic 的入口與出口
![](https://i.imgur.com/RugVST8.png)

### 欄位說明
* selector :標籤選擇器，用於指定使用哪一個 Gateway Pod 來負責此 Gateway 物件執行；一般都是直接使用 istio 預設。
* hosts: 負責提供給 Client 呼叫的位置

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

## Timeout
### Before you begin
* v1 負責外部流量，且設定 respond timeout 10s，delay 2s sql 回傳時間
* v2 負責 mesh 流量
* 設定 Request Timeouts 1s
* match 與 route 是同層級物件，所以撰寫時須留意順序

e.g.
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: generic-echo-vsvc
spec:
  hosts: 
  - "generic-echo-server-svc.default.svc.cluster.local"
  - 10.0.1.179
  gateways:
  - generic-echo-server-gateway
  - mesh
  http:
  - match:
    - gateways:
      - generic-echo-server-gateway
    route:
      - destination:
          host: generic-echo-server-svc.default.svc.cluster.local
          subset: v1
    timeout: 1s
  - route:
    - destination:
        host: generic-echo-server-svc.default.svc.cluster.local
        subset: v2
```

## 故障植入測試
### Delay
#### Before you begin
* v1 負責外部流量，且設定 respond timeout 10s，delay 10s sql 回傳時間
* v2 負責 mesh 流量
* 所以會造成外部在獲取回應時會造成 503 

e.g.
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: generic-echo-vsvc
spec:
  hosts: 
  - "generic-echo-server-svc.default.svc.cluster.local"
  - 10.0.1.179
  gateways:
  - generic-echo-server-gateway
  - mesh
  http:
  - fault:
      delay:
        fixedDelay: 7s
        percentage:
          value: 100
    match:
    - gateways:
      - generic-echo-server-gateway
    route:
      - destination:
          host: generic-echo-server-svc.default.svc.cluster.local
          subset: v1
  - route:
    - destination:
        host: generic-echo-server-svc.default.svc.cluster.local
        subset: v2
```

### Fault
#### Before you begin
* v1 負責外部流量，且設定 respond timeout 10s，delay 10s sql 回傳時間
* v2 負責 mesh 流量
* 所以會造成外部得到 500 錯誤訊息

e.g.
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: generic-echo-vsvc
spec:
  hosts: 
  - "generic-echo-server-svc.default.svc.cluster.local"
  - 10.0.1.179
  gateways:
  - generic-echo-server-gateway
  - mesh
  http:
  - fault:
      abort:
        httpStatus: 500
        percentage:
          value: 100
    match:
    - gateways:
      - generic-echo-server-gateway
    route:
      - destination:
          host: generic-echo-server-svc.default.svc.cluster.local
          subset: v1
  - route:
    - destination:
        host: generic-echo-server-svc.default.svc.cluster.local
        subset: v2
```

## Mirroring
### Before you begin
* v1 負責所有流量，且設定 respond timeout 10s，delay 0s sql 回傳時間
* v2 負責從 v1 複製過來的流量
* 記得 DestinationRule subset 也要設定正確，不能遺漏 v2 

e.g.
#### DestinationRule
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: generic-echo-server-dr
spec:
  host: generic-echo-server-svc.default.svc.cluster.local
  trafficPolicy:
    loadBalancer:
      simple: LEAST_CONN
  subsets:
  - name: v1
    labels:
      version: v1
    trafficPolicy:
      loadBalancer:
        simple: ROUND_ROBIN
  - name: v2
    labels:
      version: v2
```

#### VirtualService
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: generic-echo-vsvc
spec:
  hosts: 
  - "generic-echo-server-svc.default.svc.cluster.local"
  - 10.0.1.179
  gateways:
  - generic-echo-server-gateway
  - mesh
  http:
  - route:
    - destination:
        host: generic-echo-server-svc.default.svc.cluster.local
        subset: v1
    mirror:
      host: generic-echo-server-svc.default.svc.cluster.local
      subset: v2
    mirror_percent: 100  
```

## 鎔斷
鎔斷機制是由 DestniationRule 所負責。
### outlierDetection
* consecutiveErrors: 允許出錯的次數
* interval : 每秒做幾次請求次數
* baseEjectionTime: 發生故障的 pod 最少在被移除多久後才能再次加入負載平衡池
* maxEjectionPercent: 可以從負載平衡池中移除多少 % 的 Pod

e.g.
``` yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: generic-echo-server-dr
spec:
  host: generic-echo-server-svc.default.svc.cluster.local
  subsets:
  - name: v1
    labels:
      version: v1
    trafficPolicy:
      loadBalancer:
        simple: ROUND_ROBIN
      connectionPool:
        tcp:
          maxConnections: 1
        http:
          http1MaxPendingRequests: 1
          maxRequestsPerConnection: 1
      outlierDetection:
        consecutiveErrors: 1
        interval: 1s
        baseEjectionTime: 3m
        maxEjectionPercent: 100
  - name: v2
    labels:
      version: v2
```

## Demo
此次新增了泛用的 [demo project](https://github.com/LupinChiu/kube-ansible/tree/master/istio_demo/generic-echo)，往後的範例將會透過修改 ConfigMap 的方式啟用不同內容。

### 說明
#### server configmap
* protocol:
    * name: 什麼 protocol
    * enable: 是否啟用該 protocol 內容
    * port: 哪一個 port 監聽

<span class=red>若全部啟用，會以優先讀到的順序為主。</span>

* db:
    * external: 是否啟用連接外部 db
    * host: db IP or FQDN
    * type: db 類型
    * db: db 名稱
    * hardcodedelay: 負責模擬延遲

``` yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: generic-echo-server-cm1
data:
  CONF: |
    tcp:
      name: tcp
      enable: false
      port: :31400
    http:
      name: http
      enable: true
      port: :8033
    https:
      name: https
      enable: false
      port: :8093
    grpc:
      name: grpc
      enable: false
      port: :31402
    db:
      external: true
      host: 10.0.1.101:3308
      user: root
      passwd: 1234
      type: mysql
      db: arpg
      hardcodedelay: 0
      conn:
        maxidle: 50
        maxopen: 120
```

#### client configmap
* protocol:
    * name: 什麼 protocol
    * enable: 是否啟用該 protocol 內容
    * host: server IP or FQDN
    * interval: 多久發一次請求 (單位 s)

<span class=red>若全部啟用，會以優先讀到的順序為主。</span>

``` yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: generic-echo-client
data:
  CONF: |
    tcp:
      name: tcp
      enable: false
      host: echo-server.default.svc.local:31400
      interval: 1 #s
    http:
      name: http
      enable: true
      host: http://generic-echo-server-svc.default.svc.cluster.local:8033/echo
      interval: 1 #s
    https:
      name: https
      enable: false
      host: echo-server.default.svc.local:8034
      interval: 1 #s
    grpc:
      name: grpc
      enable: false
      host: echo-server.default.svc.local:8022
      interval: 1 #s
```

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
