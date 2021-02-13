#!/bin/bash

pip3 install pre-commit==2.7.1

pre-commit install
pre-commit install --hook-type commit-msg

yarn global add @graphql-inspector/graphql-loader @graphql-inspector/git-loader @graphql-inspector/diff-command @graphql-inspector/cli graphql
