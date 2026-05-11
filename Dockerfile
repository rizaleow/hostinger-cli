FROM gcr.io/distroless/static-debian12:nonroot
COPY hostinger-cli /usr/local/bin/hostinger-cli
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/hostinger-cli"]
