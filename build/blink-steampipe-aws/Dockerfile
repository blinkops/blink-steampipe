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
RUN git clone https://github.com/turbot/steampipe-mod-aws-compliance --branch v0.66
RUN git clone https://github.com/turbot/steampipe-mod-aws-perimeter --branch v0.2
RUN git clone https://github.com/turbot/steampipe-mod-aws-tags --branch v0.7
RUN git clone https://github.com/turbot/steampipe-mod-aws-thrifty --branch v0.18

# Setup for the aws-well-architected mod and its dependencies
# This mod is unique as it acts as a mod wrapper of aws-compliance.
#
# As a pre-condition to running this mod, it is required to
# run steampipe mod install. However, as it's not possible to lock the dependencies (aws plugin/compliance mod)
# at this time, future steampipe updates may break existing automations,
# if the inner mod installation is done dynamically (either in dockerfile or code).
#
# To overcome this, below is a "snapshot" of the mod installation flow.
RUN git clone https://github.com/turbot/steampipe-mod-aws-well-architected --branch v0.8
RUN mkdir -p /mods/steampipe-mod-aws-well-architected/.steampipe/mods/github.com/turbot/steampipe-mod-aws-compliance@v0.66.0
RUN cp -R /mods/steampipe-mod-aws-compliance/. /mods/steampipe-mod-aws-well-architected/.steampipe/mods/github.com/turbot/steampipe-mod-aws-compliance@v0.66.0
RUN ls -a /mods/steampipe-mod-aws-well-architected/.steampipe/mods/github.com/turbot/steampipe-mod-aws-compliance@v0.66.0
COPY build/blink-steampipe-aws/aws-well-architected-mod-cache.json /mods/steampipe-mod-aws-well-architected/.mod.cache.json

FROM --platform=${TARGETPLATFORM:-linux/amd64} turbot/steampipe:0.20.5

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install aws@latest aws@0.106.0 aws@0.92.1

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY --from=mods /mods /home/steampipe

USER root:0
RUN chown -R steampipe /home/steampipe/steampipe-mod-aws-well-architected

USER steampipe:0
RUN steampipe mod list --mod-location /home/steampipe/steampipe-mod-aws-well-architected

COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config/aws.spc /home/steampipe/.steampipe/config/aws.spc
COPY config/db.spc /home/steampipe/.steampipe/config/db.spc
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh
RUN apt-get update && apt-get install -y git

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
