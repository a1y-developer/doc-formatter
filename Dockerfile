from golang:1.25.4 as build
COPY . /src
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg --mount=type=cache,target=/root/.cache/go-build make build-local-linux


FROM --platform=linux/amd64 ubuntu:22.04 AS base
# Install DOC_FORMATTER Dependencies
RUN apt-get update -y && apt-get install python3 python3-pip git -y
# DOC_FORMATTER PATH
ENV PATH="/root/go/bin:${PATH}"
ENV DOC_FORMATTER_HOME="$HOME/.doc-formatter"
ENV LANG=en_US.utf8

FROM base AS goreleaser
COPY df /usr/local/bin/df
RUN /usr/local/bin/df

FROM base
COPY --from=build /src/_build/bundles/doc-formatter-linux/bin/df /usr/local/bin/df
RUN /usr/local/bin/df
