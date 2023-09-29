# Pulumi Automation Api Apps

This repository contains a set of tools based on pulumi [automation api](https://www.pulumi.com/automation/).


The list of tools:
- [hetzner-snapshots-manager](hetzner-snapshots-manager/README.md) - tool for managing snapshots of hetzner cloud servers. It's a simple tool that creates snapshots before pulumi program deletes servers. It also delete stalled snapshots. It's a good example of how to use pulumi automation api for creating a simple tool.
