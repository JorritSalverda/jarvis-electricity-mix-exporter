## Installation

To install this application using Helm run the following commands: 

```bash
helm repo add jorritsalverda https://helm.jorritsalverda.com
kubectl create namespace jarvis-electricity-mix-exporter

helm upgrade \
  jarvis-electricity-mix-exporter \
  jorritsalverda/jarvis-electricity-mix-exporter \
  --install \
  --namespace jarvis \
  --set secret.gcpServiceAccountKeyfile='{abc: blabla}' \
  --set secret.entsoeToken='token' \
  --wait
```

# Misc

Losely based on https://github.com/tmrowco/electricitymap-contrib/blob/master/parsers/NL.py