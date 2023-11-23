FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.4-alpine as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/blinkops/blink-steampipe

RUN apk add --no-cache make

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine/git:2.36.3 as mods

WORKDIR /mods
RUN git clone https://github.com/turbot/steampipe-mod-gcp-compliance --branch v0.13
RUN git clone https://github.com/turbot/steampipe-mod-gcp-labels --branch v0.6
RUN git clone https://github.com/turbot/steampipe-mod-gcp-thrifty --branch v0.10

FROM --platform=${TARGETPLATFORM:-linux/amd64} turbot/steampipe:0.18.2

# Install the gcp and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install gcp@0.32.0

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY --from=mods /mods /home/steampipe
COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config/gcp.spc /home/steampipe/.steampipe/config/gcp.spc
COPY config/db.spc /home/steampipe/.steampipe/config/db.spc
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh
RUN apt-get update && apt-get install -y git

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
