# Overview

[![License: Apache](https://img.shields.io/badge/License-Apache-yellow.svg)](https://opensource.org/licenses/MIT)

Server: ![kubernetes](https://img.shields.io/badge/-kubernetes-black.svg?logo=kubernetes&style=flat) ![docker](https://img.shields.io/badge/-docker-black.svg?logo=docker&style=flat) ![golang](https://img.shields.io/badge/-go-black.svg?logo=go&style=flat)

Client: ![golang](https://img.shields.io/badge/-go-black.svg?logo=go&style=flat) ![c#](https://img.shields.io/badge/-csharp-black.svg?logo=csharp&style=flat) ![c++](https://img.shields.io/badge/-c++-black?logo=c%2B%2B&style=flat)

This repository is a template for **Diarkis** engine's server-side code for a project.

One thing to note about this project is that **Diarkis** itself is proprietary software and will not pass the build process as is, but must be built using diarkis-cli.

A skeleton source code with function definitions and variable definitions is also provided for use with the IDE's completion functions, etc.

Diarkis server cluster is made up with HTTP, TCP, and UDP servers.

Each protocol servers run independently within the cluster, but you do not have to have all protocols.

Only HTTP server is required in the cluster and the rest of the servers should be chosen according to your application's requirements.

# How To Use The Template

The repository itself is under the src directory.
When you actually start using this repository, it is assumed that you will start your project using the src directory as a template.
To generate it, use the following command.

`make init project_id={project ID} builder_token={build token} output={absolute path to install}`

The repository itself is under the src directory.
When you actually start using this repository, it is assumed that you will start your project using the src directory as a template.
To generate it, use the following command.

To build, you will need the build token to build diarkis, which can be obtained by contacting us at https://diarkis.io .

The structure of the project after it is generated and how to use it can be found in the Readme that is included in the generated directory, which should help.
