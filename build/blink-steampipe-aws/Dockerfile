FROM golang:1.18 as builder

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM alpine/git:2.36.3 as mods

WORKDIR /mods
RUN git clone https://github.com/turbot/steampipe-mod-aws-compliance --branch v0.55
RUN git clone https://github.com/turbot/steampipe-mod-aws-perimeter --branch v0.2
RUN git clone https://github.com/turbot/steampipe-mod-aws-tags --branch v0.7
RUN git clone https://github.com/turbot/steampipe-mod-aws-thrifty --branch v0.18

FROM turbot/steampipe:0.18.2

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install aws@0.105.0

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY --from=mods /mods /home/steampipe
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
