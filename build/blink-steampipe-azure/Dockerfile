FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.4-alpine as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine/git:2.36.3 as mods

WORKDIR /mods
RUN git clone https://github.com/turbot/steampipe-mod-azure-compliance --branch v0.29
RUN git clone https://github.com/turbot/steampipe-mod-azure-tags --branch v0.6
RUN git clone https://github.com/turbot/steampipe-mod-azure-thrifty --branch v0.9
RUN git clone https://github.com/turbot/steampipe-mod-azure-insights --branch v0.10

FROM --platform=${TARGETPLATFORM:-linux/amd64} turbot/steampipe:0.18.2

# Install the azure and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install azure@0.45.0 azuread@0.11.0

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY --from=mods /mods /home/steampipe
COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config/db.spc /home/steampipe/.steampipe/config/db.spc
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh
RUN apt-get update && apt-get install -y git

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
