FROM golang:1.25.4 AS build
COPY . /src
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build make build-local-linux


FROM ubuntu:22.04 AS base
# Install DOC_FORMATTER Dependencies
RUN apt-get update -y && apt-get install python3 python3-pip git -y
# DOC_FORMATTER PATH
ENV PATH="/root/go/bin:${PATH}"
ENV LANG=en_US.utf8

FROM base AS goreleaser
COPY doc-formatter /usr/local/bin/df
RUN /usr/local/bin/df

FROM base
COPY --from=build /src/_build/bundles/doc-formatter-linux/bin/df /usr/local/bin/df
RUN /usr/local/bin/df
