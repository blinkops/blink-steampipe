FROM turbot/steampipe

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install steampipe aws github azure gcp

COPY config/ /home/steampipe/.steampipe/config/
COPY plugins/ /home/steampipe/.steampipe/plugins/local/

USER root:0
RUN chown -R steampipe /home/steampipe/.steampipe/plugins/local
RUN chown -R steampipe /home/steampipe/.steampipe/config/

USER steampipe:0
