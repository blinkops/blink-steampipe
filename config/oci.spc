connection "oci" {
    plugin = "oci@0.17.2"
    tenancy_ocid     = "{{TENANCY_OCID}}"
    user_ocid = "{{USER_OCID}}"
    fingerprint = "{{FINGERPRINT}}"
    private_key_path = "/home/steampipe/.ssh/oci_private.pem"
    regions = ["{{REGION}}"]
}
