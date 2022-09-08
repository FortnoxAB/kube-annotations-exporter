FROM gcr.io/distroless/static-debian11:nonroot
COPY kube-annotations-exporter /kube-annotations-exporter
USER nonroot
ENTRYPOINT ["/kube-annotations-exporter"]
