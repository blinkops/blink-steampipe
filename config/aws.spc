connection "aws" {
    plugin = "aws@0.92.1"
    regions = ["{{REGION}}"]
    {{ACCESS_KEY}}
    {{SECRET_KEY}}
    {{SESSION_TOKEN}}
}
