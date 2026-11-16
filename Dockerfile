FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG TARGETPLATFORM
COPY artifacts/build/release/$TARGETPLATFORM/bank /bin/bank

EXPOSE 8080

ENTRYPOINT ["/bin/bank"]
