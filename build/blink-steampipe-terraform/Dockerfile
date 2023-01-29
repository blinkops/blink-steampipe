FROM golang:1.18 as builder

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM turbot/steampipe:0.18.2

# Install the terraform and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install terraform@0.2.0

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
