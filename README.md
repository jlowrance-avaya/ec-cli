## TLDR
```
docker image build -t ec .
```

```
docker run -it \
-p 8080:8080 \
-e PROVISIONER_API_ENDPOINT='localhost' \
-e PROVISIONER_API_TOKEN='token' \
ec get deploymentManifest
```