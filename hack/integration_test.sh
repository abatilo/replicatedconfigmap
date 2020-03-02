#!/bin/bash
set -ex

set +x
echo "##########################################"
echo "# Verify creation of replicatedconfigmap #"
echo "##########################################"
set -x
kubectl get replicatedconfigmap replicatedconfigmap-sample -o json | jq '.spec.data'

set +x
echo "############################################################################"
echo "# Verify creation in syncable and syncable2, but no creation in unsyncable #"
echo "############################################################################"
set -x
kubectl -n syncable get configmap rcm-test -o json | jq '.data'
kubectl -n syncable2 get configmap rcm-test -o json | jq '.data'
kubectl -n unsyncable get configmap

set +x
echo "###########################################"
echo "# Patch unsyncable namespace, verify data #"
echo "###########################################"
set -x
kubectl patch ns unsyncable -p '{"metadata":{"annotations":{"rcm-sync":"true"}}}'
kubectl -n syncable2 get configmap rcm-test -o json | jq '.data'

set +x
echo "################################################################"
echo "# Patch syncable namespace configmap, data will be overwritten #"
echo "################################################################"
set -x
kubectl -n syncable patch configmap rcm-test -p '{"data":{"bar":"bazinga"}}'
kubectl -n syncable get configmap rcm-test -o json | jq '.data'

set +x
echo "#####################################################################"
echo "# Patch ReplicatedConfigMap and see that all ConfigMaps get updated #"
echo "#####################################################################"
set -x
kubectl patch replicatedconfigmap replicatedconfigmap-sample --type=merge -p '{"spec":{"data":{"beep":"boop"}}}'
kubectl -n syncable get configmap rcm-test -o json | jq '.data'
kubectl -n syncable2 get configmap rcm-test -o json | jq '.data'
kubectl -n unsyncable get configmap rcm-test -o json | jq '.data'

set +x
echo "####################################################"
echo "# Delete all configmaps but they will be recreated #"
echo "####################################################"
set -x
kubectl -n syncable delete configmap rcm-test
kubectl -n syncable2 delete configmap rcm-test
kubectl -n unsyncable delete configmap rcm-test
kubectl -n syncable get configmap rcm-test -o json | jq '.data'
kubectl -n syncable2 get configmap rcm-test -o json | jq '.data'
kubectl -n unsyncable get configmap rcm-test -o json | jq '.data'

set +x
echo "##############################"
echo "# Delete ReplicatedConfigMap #"
echo "##############################"
set -x
kubectl delete replicatedconfigmap replicatedconfigmap-sample
kubectl -n syncable get configmap
kubectl -n syncable2 get configmap
kubectl -n unsyncable get configmap
