
🔄 **Auto-Scaling WebSocket Applications on OpenShift Using Custom Metrics with KEDA**

In the age of real-time applications and microservices, traditional CPU/memory-based autoscaling often falls short. 
What if you could scale your applications based on business-specific metrics like active WebSocket connections? 
With OpenShift, Prometheus, and the Custom Metrics Autoscaler (KEDA), it’s possible to autoscale pods based on custom metrics that align with your application's real-time behavior.
This article walks you through the entire process of implementing WebSocket-based autoscaling using custom metrics.

📄 Prerequisites
- Access to an OpenShift 4.18+ cluster with cluster-admin permissions
- Basic knowledge of pods, deployments, routes, and metrics
- Node.js installed on your local machine for load testing

✅ **Step 1: Login and Project Setup**

```
oc login -u <user> -p <password> --server=https://<openshift_api_endpoint>
oc new-project custom-metric-autoscaler-demo
```

✅ **Step 2: Install Custom Metrics Autoscaler Operator**

Go to Operators > OperatorHub in OpenShift Web Console.

Search and install Custom Metrics Autoscaler.

Choose:

Installation Mode: All namespaces

Installed Namespace: openshift-keda

Verify the installation:

```
oc get all -n openshift-keda
```

You should see the operator and pod running.

```
 NAME                                                      READY   STATUS    RESTARTS   AGE
pod/custom-metrics-autoscaler-operator-67c5988854-sgc9r   1/1     Running   0          100s

NAME                                                 READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/custom-metrics-autoscaler-operator   1/1     1            1           101s

NAME                                                            DESIRED   CURRENT   READY   AGE
replicaset.apps/custom-metrics-autoscaler-operator-67c5988854   1         1         1       101s
```

✅ Step 3: Deploy KedaController
Create keda-controller.yaml:
```
#cat keda-controller.yaml
apiVersion: keda.sh/v1alpha1
kind: KedaController
metadata:
  name: keda
  namespace: openshift-keda
spec:
  watchNamespace: 'custom-metric-autoscaler-demo'
```
Apply:

```
oc apply -f keda-controller.yaml
oc get kedacontrollers -n openshift-keda
```

You should see similar output :

```
NAME   AGE
keda   22h
```

We have configured the controller to only scale our created namespace 
```custom-metric-autoscaler-demo```

✅ **Step 4: Deploy WebSocket Server with Custom Metrics**

For this demo, we will use a trivial web socket server using nodejs that emits the metric which counts the active count of web socket connections.

You can review the code [here](https://github.com/summu1985/openshift-demo-app-sourcecode/blob/main/node/web-socket-server.js)

Deploy the application to our Openshift cluster using S2I

```
oc new-app --name demo-ws-server --context-dir=node --strategy=docker https://github.com/summu1985/openshift-demo-app-sourcecode.git#main
oc expose service/demo-ws-server
oc get svc,deploy,route,pods
NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/demo-ws-server   ClusterIP   172.30.91.245   <none>        8080/TCP   22h

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/demo-ws-server   1/1     1            1           22h

NAME                                      HOST/PORT                                                                                       PATH   SERVICES         PORT       TERMINATION   WILDCARD
route.route.openshift.io/demo-ws-server   demo-ws-server-custom-metric-autoscaler-demo.apps.cluster-5ng44.5ng44.sandbox1582.opentlc.com          demo-ws-server   8080-tcp                 None

NAME                                  READY   STATUS    RESTARTS   AGE
pod/demo-ws-server-5c9cffdf74-ss8qj   1/1     Running   1          21h
oc status
W0414 15:58:33.418831   16351 warnings.go:70] apps.openshift.io/v1 DeploymentConfig is deprecated in v4.14+, unavailable in v4.10000+
In project custom-metric-autoscaler-demo on server https://api.cluster-5ng44.5ng44.sandbox1582.opentlc.com:6443

http://demo-ws-server-custom-metric-autoscaler-demo.apps.cluster-5ng44.5ng44.sandbox1582.opentlc.com to pod port 8080-tcp (svc/demo-ws-server)
  deployment/demo-ws-server deploys istag/demo-ws-server:latest <-
    bc/demo-ws-server docker builds https://github.com/summu1985/openshift-demo-app-sourcecode.git#main on istag/nodejs:18-minimal-ubi9 
    deployment #2 running for 27 seconds - 0/1 pods (warning: 2 restarts)
    deployment #1 deployed 52 seconds ago - 0/1 pods growing to 1
oc get bc
NAME             TYPE     FROM       LATEST
demo-ws-server   Docker   Git@main   1
oc logs -f bc/demo-ws-server

Cloning "https://github.com/summu1985/openshift-demo-app-sourcecode.git" ...
  Commit: bbc461437c671610debf228b6f5de8d39ec9cdb0 (Rename Containerfile to Dockerfile)
  Author: Sumit Mukherjee <69989028+summu1985@users.noreply.github.com>
  Date: Mon Apr 14 15:57:08 2025 +0530
Replaced Dockerfile FROM image image-registry.openshift-image-registry.svc:5000/openshift/nodejs:18-minimal-ubi9
....

Copying config sha256:7d8f45bdeb3fac9ed2d0edd46a528691cb5a019f857ec0aab3be74c92e9c98a5
Writing manifest to image destination
Successfully pushed image-registry.openshift-image-registry.svc:5000/custom-metric-autoscaler-demo/demo-ws-server@sha256:495586f011c7011e8037b27ea7c9804b2792001cf4ea3b3d961db88a956ceef9
Push successful
```

Confirm that the application is successfully running as a pod

```
oc logs -f deploy/demo-ws-server
WebSocket server running on port 8080
```

Verify the route and check if the metric is reported correctly via the exposed "/metrics" endpoint.

```
URL=$(oc get route demo-ws-server -o jsonpath='{.spec.host}')
curl http://$URL/metrics
```

Initially it should show:

```
websocket_connection_count 0
```

Connect via WebSocket:

```
npx wscat -c ws://$URL
```
Let us re-check the reported metric

```
curl http://$URL/metrics
websocket_connection_count 1
```

✅ **Step 5: Enable User Workload Monitoring**

To enable user workload monitoing we need to create cluster-monitoring-config.yaml with the following content:

```
#cluster-monitoring-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-monitoring-config
  namespace: openshift-monitoring
data:
  config.yaml: |
    enableUserWorkload: true
```

Create the config and confirm that the pods are running :

```oc create -f cluster-monitoring-config.yaml
oc -n openshift-user-workload-monitoring get pod
NAME                                   READY   STATUS    RESTARTS   AGE
prometheus-operator-75687bf59b-smpnn   2/2     Running   2          22h
prometheus-user-workload-0             6/6     Running   6          22h
thanos-ruler-user-workload-0           4/4     Running   4          22h
```

✅ **Step 6: Create ServiceMonitor**

Create a ServiceMonitor object to scrape the metrics from our deployment's endpoint which exposes the metrics (i.e. /metrics endpoint)

```
#cat service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: demo-ws-server-monitor
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app: demo-ws-server
  endpoints:
  - port: 8080-tcp
    path: /metrics
    interval: 15s

oc apply -f service-monitor.yaml
oc get servicemonitor/demo-ws-server-monitor
NAME                     AGE
demo-ws-server-monitor   22h
```

✅ **Step 7: Load Testing the Metric**

Test that the metrics collection is working fine by opening some web socket connection to our application.

```
LOAD_GENERATOR=ws-load-generator.js
LOAD_GENERATOR_URL=https://raw.githubusercontent.com/summu1985/openshift-demo-app-sourcecode/refs/heads/main/node/keda/$LOAD_GENERATOR
wget $LOAD_GENERATOR_URL
--2025-04-14 21:29:14--  https://raw.githubusercontent.com/summu1985/openshift-demo-app-sourcecode/refs/heads/main/node/keda/ws-load-generator.js
Resolving raw.githubusercontent.com (raw.githubusercontent.com)... 185.199.108.133, 185.199.109.133, 185.199.111.133, ...
Connecting to raw.githubusercontent.com (raw.githubusercontent.com)|185.199.108.133|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 1079 (1.1K) [text/plain]
Saving to: ‘ws-load-generator.js
URL=$(oc get route demo-ws-server -o jsonpath='{.spec.host}')
```
You can use node Cli from local desktop to execute a script written in node.js to test this. Ensure that you have node/npm installed locally.

Install some dependencies first.

```
npm install ws
```
And now execute the script

```
node $LOAD_GENERATOR ws://$URL 200
```

URL is the actual openshift route of our application which we extracted earlier.
200 means open 200 connections

```
node ws-load-generator.js ws://$URL 200
Opening 200 WebSocket connections to ws://demo-ws-server-custom-metric-autoscaler-demo.apps.cluster-5ng44.5ng44.sandbox1582.opentlc.com...
Connection 51 opened
Connection 47 opened
...
Connection 199 opened
Connection 198 opened
Connection 196 opened
Connection 197 opened
```

Validate that the metric is showing up correctly

```
curl http://$URL/metrics
websocket_connection_count 200
```
Validate the metric from Openshift web console integrated monitoring.

Navigate to Developer view of Openshift web console -> Observe -> Metrics

![Screenshot 2025-04-15 at 15 28 30](https://github.com/user-attachments/assets/19d3367d-5d2c-47c0-94fb-5f167843dbff)

![Screenshot 2025-04-14 at 20 01 11](https://github.com/user-attachments/assets/8def329f-550f-4ee2-baf5-e0dbf5aec3af)

✅ **Step 8: Configuring the custom metrics autoscaler to use OpenShift Container Platform monitoring**

Perform the following tasks in order to trigger custom metrics autoscaler using Openshift platform monitoring i.e. via the user workload monitoring that we just configured.

**Step 8a : Create a service account.**

```
oc project custom-metric-autoscaler-demo
oc create serviceaccount thanos
```

**Step 8b: Create a secret that generates a token for the service account.**

```
cat << EOF | oc apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: thanos-token
  annotations:
    kubernetes.io/service-account.name: thanos 
type: kubernetes.io/service-account-token
EOF
secret/thanos-token created
oc describe serviceaccount thanos
Name:                thanos
Namespace:           custom-metric-autoscaler-demo
Labels:              <none>
Annotations:         openshift.io/internal-registry-pull-secret-ref: thanos-dockercfg-lg7jl
Image pull secrets:  thanos-dockercfg-lg7jl
Mountable secrets:   thanos-dockercfg-lg7jl
Tokens:              thanos-token
Events:              <none>
```

**Step 8c: Create the trigger authentication.**

```
cat << EOF | oc apply -f -
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication 
metadata:
  name: keda-trigger-auth-prometheus
spec:
  secretTargetRef: 
  - parameter: bearerToken 
    name: thanos-token 
    key: token 
  - parameter: ca
    name: thanos-token
    key: ca.crt
EOF
triggerauthentication.keda.sh/keda-trigger-auth-prometheus created
```

**Step 8d: Create a role.**
```
cat << EOF | oc apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: thanos-metrics-reader
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
- apiGroups:
  - metrics.k8s.io
  resources:
  - pods
  - nodes
  verbs:
  - get
  - list
  - watch
EOF

role.rbac.authorization.k8s.io/thanos-metrics-reader created
```

**Step 8e: Add that role to the service account.**
```
cat << EOF | oc apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding 
metadata:
  name: thanos-metrics-reader 
  namespace:  custom-metric-autoscaler-demo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: thanos-metrics-reader
subjects:
- kind: ServiceAccount
  name: thanos 
  namespace: custom-metric-autoscaler-demo
EOF

rolebinding.rbac.authorization.k8s.io/thanos-metrics-reader created
```

**Step 8f: Reference the token in the trigger authentication object used by Prometheus.**

Create a ScaledObject resource and associate it with the trigger authentication. 

This object will have the scaling policies based on the authenticated triggers received from Prometheus by scraping the application metrics.
```
cat << EOF | oc apply -f -
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: ws-server-scaled-object
  namespace: custom-metric-autoscaler-demo
spec:
  scaleTargetRef:
    #apiVersion: apps.openshift.io/v1
    #kind: Deployment
    name: demo-ws-server
  minReplicaCount: 1
  maxReplicaCount: 5
  pollingInterval: 15
  cooldownPeriod: 45
  #advanced:
    #restoreToOriginalReplicaCount: false
    horizontalPodAutoscalerConfig:
      name: keda-hpa-scale-down
      behavior:
        scaleDown:
          stabilizationWindowSeconds: 60
          policies:
          - type: Percent
            value: 20
            periodSeconds: 60
  triggers:
  - type: prometheus
    metadata:
      serverAddress: https://thanos-querier.openshift-monitoring.svc.cluster.local:9092
      namespace: custom-metric-autoscaler-demo
      metricName: websocket_connection_count
      threshold: '40' 
      query: sum(websocket_connection_count{namespace="custom-metric-autoscaler-demo"})
      authModes: "bearer"
    authenticationRef: 
      name: keda-trigger-auth-prometheus
      kind: TriggerAuthentication

EOF

scaledobject.keda.sh/ws-server-scaled-object created
```

This creates a scaling policy where:

- The deployment will be minimum 1 instance and maximum 5 instances
- The deployment will be scaled up for every increase of 40 in the metric
- The deployment will be scaled down by maximum 20% of currently deployed instances for every period of 60 seconds, until it reaches the minimum deployment

🎉 Final Test: Observe Autoscaling

Run the load generator again:
```
node ws-load-generator.js ws://$URL 200
Opening 200 WebSocket connections to ws://demo-ws-server-custom-metric-autoscaler-demo.apps.cluster-5ng44.5ng44.sandbox1582.opentlc.com...
Connection 53 opened
Connection 63 opened
Connection 61 opened
…
Connection 194 opened
Connection 196 opened
Connection 199 opened
```
Watch pods scale up in OCP Console.

![Screenshot 2025-04-14 at 21 40 02](https://github.com/user-attachments/assets/9102d4ba-bd91-424f-b48a-c238eac21bbf)


![Screenshot 2025-04-14 at 21 40 23](https://github.com/user-attachments/assets/9f3c5649-c5a0-47a6-bcd9-ff9d6265d0c6)

Stop the script and observe pods scale down one-by-one after every 60s.
```
^C
Closing all WebSocket connections...
curl http://$URL/metrics
websocket_connection_count 200
```

![Screenshot 2025-04-15 at 13 33 36](https://github.com/user-attachments/assets/d4bb3aff-d49c-4022-8c36-6c412954d6b8)


![Screenshot 2025-04-14 at 21 40 02](https://github.com/user-attachments/assets/7ecd1636-b986-45cb-809a-be2dc6188965)


🙌 Wrapping Up
Using KEDA and custom metrics, you can autoscale apps based on real-time application metrics, not just resource usage. 

This approach works great for WebSockets, queue depths, API hits, or custom business logic.

