FROM golang:1.18 as builder

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM turbot/steampipe:0.18.2

# Install steampipe and all steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install aws@0.76.0 github@0.19.0 azure@0.31.0 gcp@0.26.0 kubernetes@0.10.0 terraform@0.2.0 csv@0.4.1 oci@0.17.2

# Validate steampipe is running
RUN steampipe plugin list

# Sanity check
RUN steampipe query "select 1" --output csv

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config /home/steampipe/.steampipe/config/
COPY plugins/local/oci/oci.plugin /home/steampipe/.steampipe/plugins/local/oci/oci.plugin
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]