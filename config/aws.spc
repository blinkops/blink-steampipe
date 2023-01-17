connection "aws" {
    plugin = "aws@0.90.0"
    regions = ["*"]
    {{ACCESS_KEY}}
    {{SECRET_KEY}}
    {{SESSION_TOKEN}}
}
