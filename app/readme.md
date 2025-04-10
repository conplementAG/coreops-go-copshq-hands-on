docker build . -t cp-notes:latest
docker run --rm -p 8080:8080 cp-notes:latest

docker tag cp-notes cponeneucopsacr.azurecr.io/cp-notes-**>your-namespace<**:latest
docker push cponeneucopsacr.azurecr.io/cp-notes-**>your-namespace<**:latest