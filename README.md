# Latency meter

## Client
1. Build<br>`docker build -t lm-client . -f build/client/Dockerfile`
2. Run:<br>`docker run --network kind --rm -v /c/Users/loren/Desktop/kubeconfig.txt:/root/kubeconfig lm-client -kubeconfig=/root/kubeconfig -label=app=lm-server -lmport=30007`

## Client
1. Build<br>`docker build -t devrols/lm-server . -f build/server/Dockerfile`
2. Push:<br>`docker push devrols/lm-server`
3. Start (control plane):<br>`kubectl apply -f lm-server-deployment.yml` and `kubectl apply -f lm-server-service.nodeport.yml`