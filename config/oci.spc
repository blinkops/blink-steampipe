connection "oci" {
    plugin           = "local/oci"
    tenancy_ocid     = "{{TENANCY_OCID}}"
    user_ocid = "{{USER_OCID}}"
    fingerprint = "{{FINGERPRINT}}"
    private_key_path = "~/.ssh/oci_private.pem"
    regions = ["{{REGION}}"]
}
