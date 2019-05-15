## Script

Demos,
1. Install a source; target a function.
1. Source targets the broker, trigger a function.
1. First Party events.

### Install a Source Demo

1. Setup
    - Empty cluster with Serving and Eventing installed.
    - Namespace labeled with inject broker.
    
1. Demo
    - Install GCS source.
    - Point to a function to display event.

1. Highlights
    - Function does not have GCP creds.

### Source to Broker Demo

1. Setup
    - Source is installed.

1. Demo
    - Repoint the source to the broker.
    - Create a trigger from broker to function.
    - Show function sending events.

1. Highlights
    - Show function sending events.


### First Party Events

1. Setup
    - Source pointed to Broker.

1. Demo
    - Install a function that splits the file and makes an event.
    - The event is sent to the broker directly.
    - Create other functions that react to it.
    - Show function sending events.

1. Highlights
    - Show function sending events.
    


---


To make the cluster:

```shell
 export CLUSTER_NAME=<cluster-name>
 export CLUSTER_ZONE=europe-west1-b
 export PROJECT=<project>
 export PROJECT_ID=$PROJECT
 gcloud config set core/project $PROJECT
 gcloud services enable \
    cloudapis.googleapis.com \
    container.googleapis.com \
    containerregistry.googleapis.com
    
 gcloud beta container clusters create $CLUSTER_NAME \
   --addons=HorizontalPodAutoscaling,HttpLoadBalancing,Istio \
   --machine-type=n1-standard-4 \
   --cluster-version=latest --zone=$CLUSTER_ZONE \
   --enable-stackdriver-kubernetes --enable-ip-alias \
   --enable-autoscaling --min-nodes=1 --max-nodes=10 \
   --enable-autorepair \
   --scopes cloud-platform

 kubectl create clusterrolebinding cluster-admin-binding \
   --clusterrole=cluster-admin \
   --user=$(gcloud config get-value core/account)
```

Install Knative v0.6.0:

```shell
 kubectl apply --selector knative.dev/crd-install=true \
 --filename https://github.com/knative/serving/releases/download/v0.6.0/serving.yaml \
 --filename https://github.com/knative/eventing/releases/download/v0.6.0/release.yaml \
 --filename https://github.com/knative/eventing-sources/releases/download/v0.6.0/eventing-sources.yaml

 kubectl apply \
 --filename https://github.com/knative/serving/releases/download/v0.6.0/serving.yaml \
 --filename https://github.com/knative/eventing/releases/download/v0.6.0/release.yaml \
 --filename https://github.com/knative/eventing-sources/releases/download/v0.6.0/eventing-sources.yaml
```

Install GCS source:

```shell
 gcloud services enable pubsub.googleapis.com
 kubectl apply -f https://github.com/knative/eventing-sources/releases/download/v0.6.0/gcppubsub.yaml
 
 gcloud projects add-iam-policy-binding $PROJECT_ID \
   --member=serviceAccount:knative-source@$PROJECT_ID.iam.gserviceaccount.com \
   --role roles/pubsub.editor
 
 gcloud iam service-accounts keys create knative-source.json \
   --iam-account=knative-source@$PROJECT_ID.iam.gserviceaccount.com
   
 # Note that the first secret may already have been created when installing
 # Knative Eventing. The following command will overwrite it. If you don't
 # want to overwrite it, then skip this command.
 kubectl -n knative-sources create secret generic gcppubsub-source-key --from-file=key.json=knative-source.json --dry-run -o yaml | kubectl apply --filename -
 
 # The second secret should not already exist, so just try to create it.
 kubectl -n default create secret generic google-cloud-key --from-file=key.json=knative-source.json

 kubectl apply -f https://raw.githubusercontent.com/vaikas-google/gcs/master/release.yaml
```

Enable Eventing:

```shell
 # edit the default namespace and add the label knative-eventing-injection=enabled
 kubectl get namespaces -l knative-eventing-injection=enabled
```