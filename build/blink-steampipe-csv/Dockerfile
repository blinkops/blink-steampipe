FROM golang:1.18 as builder

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts
COPY internal ./internal

RUN make install

FROM turbot/steampipe:0.16.1

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install csv@0.4.1

COPY --from=builder /bin/generate /home/steampipe/bin/
COPY docker-entrypoint.sh /home/steampipe/bin

USER root:0
COPY config/csv.spc /home/steampipe/.steampipe/config/csv.spc
COPY config/db.spc /home/steampipe/.steampipe/config/db.spc
RUN chown -R steampipe /home/steampipe/bin/ /home/steampipe/.steampipe/config
RUN chmod -R +x /home/steampipe/bin/
RUN chown steampipe /home/steampipe/bin/docker-entrypoint.sh

USER steampipe:0

ENTRYPOINT ["/home/steampipe/bin/docker-entrypoint.sh"]
