#!/bin/bash

# Fail if error
set -ex

# Go into ui
cd ui

# Build UI
yarn build

# Come back
cd ..

# Copy UI into backend static
cp -Rf ui/build/* backend/static/

# Go into backend
cd backend

# Build cross platform backend and release it
make release/all

# Remove new static
rm -Rf static/

# Reset static
git checkout static/

# Come back
cd ..
