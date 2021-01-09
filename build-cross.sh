#!/bin/bash

# Fail if error
set -ex

# Go into ui
cd ui

# Build UI
yarn build

# Come back
cd ..

# Create static directory just in case
mkdir -p backend/static/

# Copy UI into backend static
cp -Rf ui/build/* backend/static/

# Go into backend
cd backend

# Build cross platform backend
make code/build-cross

# Come back
cd ..
