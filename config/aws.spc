connection "aws" {
    plugin = "aws@0.92.1"
    regions = ["*"]
    {{ACCESS_KEY}}
    {{SECRET_KEY}}
    {{SESSION_TOKEN}}
}
