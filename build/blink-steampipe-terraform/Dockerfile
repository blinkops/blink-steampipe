FROM golang:1.18 as builder

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY ../../go.mod go.sum Makefile ./
COPY ../../scripts ./scripts
COPY ../../internal ./internal

RUN make install

FROM turbot/steampipe:0.16.1

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install terraform@0.2.0

# Validate steampipe is running
RUN steampipe plugin list

# Sanity check
RUN steampipe query "select 1" --output csv

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY ../../docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY ../../config/terraform.spc /home/steampipe/.steampipe/config/terraform.spc
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh

USER steampipe:0
ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
