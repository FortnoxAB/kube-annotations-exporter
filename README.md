# kube annotations exporter

Created for people using a lower version of kube prometheus. This feature exists in kube-state-metrics v2.2.1

## Download

To download the binary to your computer.

`go get github.com/fortnoxab/kube-annotations-exporter/cmd/kube-annotations-exporter`

Or use the providied docker image from quay.

`quay.io/fortnox/kube-annotations-exporter`

## How to use

This will export annotations that matched the given annotations from the CLI options. It will only export the kind, name and the annotation on the metrics given for the object. If it does not match anything it won't export the metric (useless data for prometheus).
```
kube_annotations_exporter{annotation="my.annotation",kind="Service",name="my-service"} 1
```

### Configuration flags
```
Usage of ./kube-annotations-exporter:
  -annotations
    	Change value of Annotations. (default [])
  -port
    	Change value of Port. (default 8080)

Generated environment variables:
   CONFIG_ANNOTATIONS
   CONFIG_PORT
```