# Overview

This directory contains standalone Diarkis examples that are able to be built and executed locally.

## Build Instructions

In order to build the following examples you must have a valid Diarkis **Project ID** and its
corresponding **Build Token**. 

Next invoke the following command from the root of the
`diarkis-server-template` repository:

```sh
make examples project_id=PROJECT_ID builder_token=BUILD_TOKEN output=OUTPUT
```

Once the example installation is finished, navigate to your `OUTPUT_DIRECTORY` and explore.

## Examples

### Dive

- [User Online Status](./dive/user-online-status/README.md)

### HTTP

- [JSON Endpoint](./http/json-endpoint/README.md)

### Matchmaker

- [Simple Candidate Pooling](./matching/simple/README.md)
- [Matchmaking via Custom Criteria](./matching/custom-criteria/README.md)
- [Team-Based Matchmaking](./matching/team/README.md)
- [Secure Matchmaking](./matching/secure-matching/README.md)



---

_First created on 2025-07-03. Updated on 2025-07-04._

