FROM docker.io/library/golang:1.22.3-alpine3.18 as build

WORKDIR /marmitton

## git is needed for the go build process with the library github.com/earthboundkid/versioninfo
RUN apk update && apk add --no-cache \
    tini-static \
    git
COPY . .
RUN go build -ldflags="-s -w"

# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static:nonroot
USER nonroot:nonroot

COPY --from=build /marmitton/marmitton /app/marmitton
COPY --from=build --chown=nonroot:nonroot /sbin/tini-static /sbin/tini

ENTRYPOINT ["/sbin/tini", "--", "/app/marmitton"]
