FROM golang:1.18 as builder

WORKDIR /go/src/github.com/blinkops/blink-steampipe

COPY go.mod go.sum Makefile ./
COPY scripts ./scripts

RUN make install

FROM turbot/steampipe:0.16.1

USER steampipe:0
COPY config /home/steampipe/blink/config
COPY --from=builder /bin/generate /home/steampipe/blink/bin/
COPY docker-entrypoint.sh /home/steampipe/blink/bin

RUN steampipe plugin install aws github@0.19.0 azure@0.31.0 gcp@0.26.0 kubernetes@0.10.0

RUN steampipe plugin list

USER root:0
RUN chown -R steampipe /home/steampipe/blink/bin/
RUN chmod -R +x /home/steampipe/blink/bin/
RUN chown steampipe /home/steampipe/blink/bin/docker-entrypoint.sh

USER steampipe:0
ENTRYPOINT ["/home/steampipe/blink/bin/docker-entrypoint.sh"]
CMD ["steampipe"]
