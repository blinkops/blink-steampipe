connection "gcp" {
  plugin = "gcp@0.32.0"
  project = "{{PROJECT}}"
  credentials = "~/.config/gcloud/application_default_credentials.json"
}
