#!/bin/sh

# copy pre-commit to .git/hooks folder and give it execution permission

cp scripts/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit