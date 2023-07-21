## TLDR
```
docker image build -t acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0 .
```

```
docker run -it \
-p 8080:8080 \
-e PROVISIONER_API_ENDPOINT='10.71.2.5' \
-e PROVISIONER_API_TOKEN='token' \
acraocpshsrvnonprod.azurecr.io/ec-cli:v1.0.0 get deploymentManifest
```

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