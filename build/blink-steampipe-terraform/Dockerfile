FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.4-alpine3.18 as builder

ENV GOPRIVATE=github.com/blinkops

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/blinkops/blink-steampipe

RUN apk add --no-cache make git openssh-client

RUN mkdir -p ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN git config --global url."git@github.com:blinkops/blink-sdk.git".insteadOf "https://github.com/blinkops/blink-sdk"

COPY go.mod go.sum Makefile ./
RUN --mount=type=ssh go mod download

COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM --platform=${TARGETPLATFORM:-linux/amd64} turbot/steampipe:0.18.2

# Install the terraform and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install terraform@0.3.0

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config/terraform.spc /home/steampipe/.steampipe/config/terraform.spc
COPY config/db.spc /home/steampipe/.steampipe/config/db.spc
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
