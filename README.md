## TLDR
```
podman image build -t acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0 .
```

```
podman run -it \
-p 8080:8080 \
acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0 get deploymentManifest
```

## Podman instructions

### build image with podman
```
podman image build -t acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0 .
```

### run this command and copy the access token
```
az acr login --name acraocpshsrvnonprod --expose-token
```

### push the image using the `--creds` flag and paste the access token you had copied
```
podman push --creds=00000000-0000-0000-0000-000000000000:"<access token>" acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0
```

## Docker instructions

### WARNING: we are not supposed to be using docker.

### push to ACR with docker 
```
az acr login --name acraocpshsrvnonprod
docker image build -t acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0 .
docker push acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0
```

### pull image from server
```
az acr login --name acraocpshsrvnonprod
docker image pull acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0
```

## troubleshooting

### x509: certificate signed by unknown authority
```
 => ERROR [builder 4/6] RUN go mod download                                                                                                 0.5s
------
 > [builder 4/6] RUN go mod download:
0.503 go: github.com/alecthomas/kingpin@v2.2.6+incompatible: Get "https://proxy.golang.org/github.com/alecthomas/kingpin/@v/v2.2.6+incompatible.mod": x509: certificate signed by unknown authority
```

This is happening because zscaler doesn't trust `https://proxy.golang.org/github.com/alecthomas/kingpin/@v/v2.2.6+incompatible.mod`. To work around this issue, temporarily disable zscaler network protection and run the `docker build` command again.