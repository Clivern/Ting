FROM golang:1.27.1 AS builder

ARG TING_VERSION=0.1.0
ARG TING_COMMIT=none
ARG TING_BUILD_DATE=unknown
ARG TING_BUILT_BY=docker

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${TING_VERSION} -X main.commit=${TING_COMMIT} -X main.date=${TING_BUILD_DATE} -X main.builtBy=${TING_BUILT_BY}" \
    -o /out/ting .

RUN mkdir -p /out/configs && \
    cp config.dist.yml /out/configs/config.dist.yml

FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder --chown=65532:65532 /out/ting /app/ting
COPY --from=builder --chown=65532:65532 /out/configs /app/configs

EXPOSE 8080

VOLUME ["/app/configs"]

USER 65532:65532

ENTRYPOINT ["/app/ting"]
CMD ["server", "-c", "/app/configs/config.dist.yml"]
