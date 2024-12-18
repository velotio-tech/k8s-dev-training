#!bin/bash
kubectl patch redis redis-sample --type=merge --subresource=status -p '{"status": {"reconciliationCount" : 10, "runningReplicas" : true}}'