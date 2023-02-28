# Nubes app example

This directory contains example of a project developed with Nubes. The `faas` directory comprises Nobjects' definitions that later are translated into serverless functions using the [nubes generator](https://github.com/Astenna/Nubes/tree/main/generator).

## Prerequisites

*the same as in the README in the root directory of the project*

Nubes requires **Golang version 1.18 or greater**.

To successfully run Nubes generator, [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) must be installed.

To deploy serverless functions:

- [AWS credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) must be configured
- [serverless framework](https://www.serverless.com/framework/docs/getting-started) must be installed

Additionally, as one of the commands required for serverless functions deployemnt uses shell script, Windows users need bash-emulating command-line (e.g. Git Bash usually installed `along with Git for Windows).

## How to run

First, run Nubes generator to generate deployment files as well as perform source-to-source code translation of Nobjects' definitions. The paths in the command are relative paths assuming the generator is run in the `example` directory.

```bash
generator handlers -t=./faas/types -o=./faas -m=github.com/Astenna/Nubes/example/faas -g=true -i=false
```

Then, run the deployment commands from within the `faas` directory.

```bash
./build_handlers.sh generated/
sls deploy --verbose
```

At this point, a set of serverless functions listed in the `faas/serverless.yml` should be deployed to the currently configured AWS account. As a next step, another run of the Nubes generator is done to generate the *client library* with Nobjects' types redefinitions to be used in projects being the clients of the defined types (= deployed serverless functions). The paths in the command are relative paths assuming the generator is run in the `example` directory.

```bash
generator client -t=./faas/types -o=./example
```