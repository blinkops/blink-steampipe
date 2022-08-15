FROM turbot/steampipe

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install steampipe aws github azure gcp kubernetes

COPY config/ /home/steampipe/.steampipe/config/
COPY scripts/bin/ /home/steampipe/bin/

USER root:0
RUN chown -R steampipe /home/steampipe/bin/
RUN chmod -R +x /home/steampipe/bin/


USER steampipe:0

ENTRYPOINT ["/home/steampipe/bin/entrypoint.sh"]
