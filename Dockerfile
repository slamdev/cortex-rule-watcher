FROM gcr.io/distroless/static:nonroot
WORKDIR /
ADD controller controller
USER 65532:65532

ENTRYPOINT ["/controller"]
