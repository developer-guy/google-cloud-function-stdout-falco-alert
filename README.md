# google-cloud-function-stdout-falco-alert

### Prerequisites

* gcloud 342.0.0

### Tutorial

To test workloadidentity first create a GKE cluster with workloadidentity enabled

```bash
$ GOOGLE_PROJECT_ID=$(gcloud config get-value project)
$ CLUSTER_NAME=demo
$ gcloud container clusters create $CLUSTER_NAME \
                   --workload-pool ${GOOGLE_PROJECT_ID}.svc.id.goog
```

Let's deploy the Google Cloud Functions first, because in the later steps, we'll need the name of the function.

```bash
$ git clone https://github.com/developer-guy/google-cloud-function-stdout-falco-alert
$ cd google-cloud-function-stdout-falco-alert
$ gcloud functions deploy HelloWorld --runtime go113 --trigger-http
Allow unauthenticated invocations of new function [HelloWorld]? (y/N)? N
...
```

Get the name of the function
```bash
$ CLOUD_FUNCTION_NAME=$(gcloud functions describe --format=json HelloWorld | jq -r '.name')
```

Once it's created, lets install `Falco`, and `Falcosidekick` with enabled `Google Cloud Functions` output type. In order to do that, 
we should clone the `developer-guy/charts1`, and deploy the `Falco` and `Falcosidekick` through this chart.
 Because the upstream Chart repository of the `Falcosecurity` does not involve the latest upgrades yet.

> Don't forget to change [Chart.yaml](https://github.com/developer-guy/charts-1/blob/master/falco/Chart.yaml#L24) with the location of your working directory.

```bash
$ git clone https://github.com/developer-guy/charts-1
$ cd charts-1
$ helm dependency update falco
$ helm upgrade --install falco falco \
--namespace falco --create-namespace \
--set ebpf.enabled=true \
--set falcosidekick.enabled=true \
--set falcosidekick.config.gcp.cloudfunctions.name=${CLOUD_FUNCTION_NAME} \
--set falcosidekick.webui.enabled=true \
--set falcosidekick.image.repository=falcosecurity/falcosidekick \
--set falcosidekick.image.tag=latest
```

Finally set up the your SA and Rolebindings
```bash
$ SA_ACCOUNT=falco-falcosidekick-sa
$ gcloud iam service-accounts create $SA_ACCOUNT

$ gcloud projects add-iam-policy-binding ${GOOGLE_PROJECT_ID} \
--member="serviceAccount:${SA_ACCOUNT}@${GOOGLE_PROJECT_ID}.iam.gserviceaccount.com" \
--role="roles/cloudfunctions.developer"

$ gcloud projects add-iam-policy-binding ${GOOGLE_PROJECT_ID} \
--member="serviceAccount:${SA_ACCOUNT}@${GOOGLE_PROJECT_ID}.iam.gserviceaccount.com" \
--role="roles/cloudfunctions.invoker"

$ gcloud iam service-accounts add-iam-policy-binding \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:${GOOGLE_PROJECT_ID}.svc.id.goog[falco/falco-falcosidekick]" \
  ${SA_ACCOUNT}@${GOOGLE_PROJECT_ID}.iam.gserviceaccount.com
```

Finally set up the Falcosidekick SA to impersonate a GCP SA
```bash
$ kubectl annotate serviceaccount \
  --namespace falco \
  falco-falcosidekick \
  iam.gke.io/gcp-service-account=${SA_ACCOUNT}@${GOOGLE_PROJECT_ID}.iam.gserviceaccount.com
```

### Test

Create an alpine pod first, then try to exec into it.

```bash
$ kubectl run alpine  --image=alpine --restart='Never' -- sh -c "sleep 600"
```

Exec into it.
```bash
$ kubectl exec -i --tty alpine -- sh -c "uptime" 
```
