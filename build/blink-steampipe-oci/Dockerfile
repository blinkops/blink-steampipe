FROM golang:1.18 as builder

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM alpine/git:2.36.3 as mods

WORKDIR /mods
RUN git clone https://github.com/turbot/steampipe-mod-oci-compliance --branch v0.7
RUN git clone https://github.com/turbot/steampipe-mod-oci-insights --branch v0.4
RUN git clone https://github.com/turbot/steampipe-mod-oci-thrifty --branch v0.8

FROM turbot/steampipe:0.18.2

# Install the oci and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY --from=mods /mods /home/steampipe
COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config/oci.spc /home/steampipe/.steampipe/config/oci.spc
COPY config/db.spc /home/steampipe/.steampipe/config/db.spc
COPY plugins/local/oci/oci.plugin /home/steampipe/.steampipe/plugins/local/oci/oci.plugin
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
