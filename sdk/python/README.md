# Google Cloud Platform Scaffolding

This repo is a [Pulumi Package](https://www.pulumi.com/docs/guides/pulumi-packages/) used to deploy create a Google Cloud Project,
connect it to a billing account and be able to enable any APIs for the project

It's written in Go, but thanks to Pulumi's multi language SDK generating capability, it create usable SDKs for all of Pulumi's [supported languages](https://www.pulumi.com/docs/intro/languages/)

> :warning: **This package is a work in progress**: Please do not use this in a production environment!

# Building and Installing

## Building from source

But if you need to build it yourself, just download this repository, [install](https://taskfile.dev/#/installation) [Task](https://taskfile.dev/):

```sh
go get github.com/go-task/task/v3/cmd/task
```

And run the following command to build and install the plugin in the correct folder (resolved automatically based on the current Operating System):

```sh
task install
```

## Install Plugin Binary

Before you begin, you'll need to install the latest version of the Pulumi Plugin using `pulumi plugin install`:

```
pulumi plugin install resource globalgcpcloudrun v0.0.1 --server https://cobraz.jfrog.io/artifactory/pulumi-packages/pulumi-gcp-scaffold
```

This installs the plugin into `~/.pulumi/plugins`.

## Install your chosen SDK

Next, you need to install your desired language SDK using your languages package manager.

### Python

```
pip3 install cobraz-pulumi-gcp-scaffold
```

### NodeJS

```
npm install @cobraz/pulumi-gcp-scaffold
```

### DotNet

```
Coming Soon
```

### Go

```
go get -t github.com/cobraz/pulumi-gcp-scaffold/sdk/go/gcp
```

# Usage

Once you've installed all the dependencies, you can use the library like any other Pulumi SDK. See the [examples](examples/) directory for examples of how you might use it.

To create a Google Cloud Project and enable some APIs then you would use the following example:

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as scaffold from "@cobraz/pulumi-gcp-scaffold";

const proj = new scaffold.Project("my-project", {
    projectID: "my-project-name",
    orgID: "<my gcp org id>",
    billingAccountID: "<my billing account id>",
    activatedApis: [
        "compute.googleapis.com",
        "container.googleapis.com",
        "cloudbilling.googleapis.com"
    ]
})

export const projectName = proj.projectName;
export const projectID = proj.projectID;
```

# Limitations

This package currently requires the user to be aware of the `orgID` and `billingAccountID`. There are plans to make this easier.
