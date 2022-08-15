FROM turbot/steampipe

# Install the aws and steampipe plugins for Steampipe (as steampipe user).
USER steampipe:0
RUN steampipe plugin install steampipe aws github azure gcp kubernetes

COPY config/aws.spc /home/steampipe/.steampipe/config/aws.spc
COPY scripts/bin/ /home/steampipe/bin/
# COPY plugins/ /home/steampipe/.steampipe/plugins/local/

USER root:0
RUN chown -R steampipe /home/steampipe/bin/
RUN chmod -R +x /home/steampipe/bin/

# RUN chown -R steampipe /home/steampipe/.steampipe/plugins/local
# RUN chown -R steampipe /home/steampipe/.steampipe/config/

USER steampipe:0

ENTRYPOINT ["/home/steampipe/bin/entrypoint.sh"]