FROM turbot/steampipe

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install steampipe aws github

COPY config/ /home/steampipe/.steampipe/config/
