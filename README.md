# manifest_splitter

Extracts Ops Files from your YAMLs.
Optimised for BOSH manifests.

## Installation

`go get github.com/alex-slynko/manifest_splitter`

## Sample usage

```
$ manifest_splitter my_cfcr_manifest.yml kubo-deployment/manifests/cfcr.yml 
 
- type: replace
  path: /instance_groups/name=master/jobs/name=kube-apiserver/provides?
  value:
    kube-apiserver:
      as: master
      shared: true
- type: replace
  path: /instance_groups/name=master/jobs/name=etcd/provides?
  value:
    etcd:
      as: etcd
      shared: true
```
