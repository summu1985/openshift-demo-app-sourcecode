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

We have created a trivial web socket server using nodejs that emits the metric which counts the active count of web socket connections.

You can review the code [here](https://github.com/summu1985/openshift-demo-app-sourcecode/blob/main/node/web-socket-server.js)

We will deploy the application to our Openshift cluster using S2I

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
```

Let us verify the route and check if the metric is reported correctly via the exposed "/metrics" endpoint.

URL=$(oc get route demo-ws-server -o jsonpath='{.spec.host}')
curl http://$URL/metrics

Initially it should show:
websocket_connection_count 0
Connect via WebSocket:
npx wscat -c ws://$URL
Re-check:
curl http://$URL/metrics
# Output: websocket_connection_count 1

📊 Step 5: Enable User Workload Monitoring
Create cluster-monitoring-config.yaml:
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-monitoring-config
  namespace: openshift-monitoring
data:
  config.yaml: |
    enableUserWorkload: true
Apply:
oc create -f cluster-monitoring-config.yaml
Confirm pods:
oc -n openshift-user-workload-monitoring get pod

🔍 Step 6: Create ServiceMonitor
Label discovery:
oc get deploy/demo-ws-server -o json | jq '.metadata.labels'
Use label app: demo-ws-server. Create ServiceMonitor:
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
Apply it:
oc apply -f service-monitor.yaml

🚀 Step 7: Load Testing the Metric
Install Node dependency:
npm install ws
Download and run the generator:
wget https://raw.githubusercontent.com/summu1985/openshift-demo-app-sourcecode/main/node/keda/ws-load-generator.js
node ws-load-generator.js ws://$URL 200
Check metric:
curl http://$URL/metrics
# Output: websocket_connection_count 200
Stop the script:
Ctrl + C
Recheck:
curl http://$URL/metrics
# Output: websocket_connection_count 0

🔒 Step 8: Prometheus Auth via Service Account
oc create serviceaccount thanos
Create token secret:
apiVersion: v1
kind: Secret
metadata:
  name: thanos-token
  annotations:
    kubernetes.io/service-account.name: thanos
  type: kubernetes.io/service-account-token
Apply:
oc apply -f token-secret.yaml
Create TriggerAuthentication:
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

⚖️ Step 9: RBAC Permissions for Prometheus
Create Role:
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: thanos-metrics-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get"]
- apiGroups: ["metrics.k8s.io"]
  resources: ["pods", "nodes"]
  verbs: ["get", "list", "watch"]
Create RoleBinding:
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: thanos-metrics-reader
  namespace: custom-metric-autoscaler-demo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: thanos-metrics-reader
subjects:
- kind: ServiceAccount
  name: thanos
  namespace: custom-metric-autoscaler-demo

📈 Step 10: Create KEDA ScaledObject
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: ws-server-scaled-object
  namespace: custom-metric-autoscaler-demo
spec:
  scaleTargetRef:
    name: demo-ws-server
  minReplicaCount: 1
  maxReplicaCount: 5
  pollingInterval: 15
  cooldownPeriod: 45
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

🎉 Final Test: Observe Autoscaling
Run the load generator again:


node ws-load-generator.js ws://$URL 200

Watch pods scale up in OCP Console.


Stop the script and observe pods scale down after ~300s.



📂 Resources
GitHub Source Code


Red Hat Documentation



🙌 Wrapping Up
Using KEDA and custom metrics, you can autoscale apps based on real-time application metrics, not just resource usage. This approach works great for WebSockets, queue depths, API hits, or custom business logic.
If this helped you, please like, share, and comment. Would love to hear what other metrics you've scaled on!
#OpenShift #Kubernetes #Autoscaling #KEDA #DevOps #WebSockets #CustomMetricsAutoscaler

