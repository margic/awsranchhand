# AWS Ranch Hand is a set of tools for maintaining Rancher.

[![CircleCI](https://circleci.com/gh/margic/awsranchhand.svg?style=shield)](https://circleci.com/gh/margic/awsranchhand)

## Usage
run awsranchand [command] [options]

## Commands
- finish           Complete a rancher service upgrade
- labelec2instance Label Rancher host with EC2 Instance ID
- labelhost        Apply a label to a rancher host
- listhosts        A brief description of your command
- rollback         Rollback a rancher service upgrade
- upgrade          Upgrade a service
- waitforstate     Waits for a service to be in a specified state.

## Configuration
Configuration can be done via configuration file either in the same folder as
the binary or in the user home directory.

Command params can be passed as environment variables or as flags.

Environment variables:
CATTLE_ACCESS_KEY - Rancher cattle key
CATTLE_SECRET_KEY - Rancher secret
CATTLE_URL - Rancher api url ie https://rancherhost.com/v1
