# Terraform Provider for Google Play Console

> _This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._

This repository is a [Terraform](https://www.terraform.io) provider for managing resources in the [Google Play Console](https://developers.google.com/android-publisher/api-ref/rest).

This project is [available on the Terraform Registry](https://registry.terraform.io/providers/Oliver-Binns/googleplay/latest).

## Usage

### Provider

Start by declaring a Google Play provider.

This should contain:

- a base64 encoded string representation of your Google Service Account json file which is used for authenticating against the Google Play Developer API. Refer to the [Google Developer documentation](https://developers.google.com/android-publisher/getting_started/?hl=en) for more details.
- Your Google Play developer ID, this is a 19-digit number than can be found in the Google Play Console.

```tf
provider "googleplay" {
  service_account_json_base64 = filebase64("~/service-account.json")
  developer_id = "5166846112789481453"
}
```

### Managing users

You can manage Google Play Console users as a Terraform resource (`googleplay_user`).

Each user requires an email address that they will use to authenticate with the Google Play Console.

A set of [Developer Account permissions](https://developers.google.com/android-publisher/api-ref/rest/v3/users#DeveloperLevelPermission) is also required. This list can be empty.

```tf
resource "googleplay_user" "oliver" {
  email = "example@oliverbinns.co.uk"
  global_permissions = ["CAN_SEE_ALL_APPS", "CAN_MANAGE_DRAFT_APPS_GLOBAL"]
  app_permissions = []
}
```

### App specific permissions

Users can be granted specific permissions to a particular app using the `googleplay_app_iam` resource.

Each App IAM resource declares:
- the email of the user who should be granted access 
- the 19-digit app ID 
- the set of permissions to be granted

```tf
resource "googleplay_app_iam" "test_app" {
  app_id  = "0000000000000000000"
  user_id = googleplay_user.oliver.email
  permissions = [
    "CAN_REPLY_TO_REVIEWS"
  ]
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22

### Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Using the provider

Fill this in for each provider

### Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Commits

[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) are required for each pull request to ensure that release versioning can be managed automatically.
Please ensure that you have enabled the Git hooks, so that you don't get caught out!:
```
git config core.hooksPath hooks
```