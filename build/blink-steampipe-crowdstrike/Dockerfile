FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.4-alpine3.18 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/blinkops/blink-steampipe

RUN apk add --no-cache make

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM --platform=${TARGETPLATFORM:-linux/amd64} turbot/steampipe:0.18.2

# Install steampipe and all steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install crowdstrike@0.2.0

# Validate steampipe is running
RUN steampipe plugin list

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config/db.spc /home/steampipe/.steampipe/config/db.spc
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
