# (s)FAAS project example

This directory contains an example project developed with nubes. It contaings nobjects' definitions that later are translated into serverless functions using the [nubes generator](https://github.com/Astenna/Nubes/tree/main/generator).

## Requirements

To deploy the aws lambdas defined here, the following dependencies are needed:

- Makefile
- Docker
- serverless framework
- configured aws credentials
- for Windows users: bash/shell for windows (e.g. Git Bash) to run the Makefile

Note that, all of those dependencies are required for any project defined using nubes, except for the Makefile that is used only to temporarily copy the local package [lib](https://github.com/Astenna/Nubes/tree/main/lib).

## How to run

First, use generator to generate code defining the corresponding AWS lambdas:

```
generator handlers -t=<path to this repo>/types -r=<path to this repo>/faas/repositories" "-o=<path to this repo> -m=github.com/Astenna/Nubes/faas
```

To build the generator executable check readme [here](https://github.com/Astenna/Nubes/tree/main/generator). Replace the `<path to this repo>` with path to this repo root.

Once the generator is finished, you should see the `generated` directory with lambda handlers definitions added to this directory (and some other files too).

Then, to deploy generated AWS lambdas run:
```
make deploy
```