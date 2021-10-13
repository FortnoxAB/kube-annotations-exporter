FROM scratch

COPY kube-annotations-exporter /kube-annotations-exporter

ENTRYPOINT ["/kube-annotations-exporter"]