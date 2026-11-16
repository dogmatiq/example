FROM scratch

ARG TARGETPLATFORM
COPY artifacts/build/release/$TARGETPLATFORM/* /bin/

EXPOSE 8080

ENTRYPOINT ["/bin/bank"]
