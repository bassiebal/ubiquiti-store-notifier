#!/bin/zsh
source ~/.zshrc
PATH="/usr/local/bin:$PATH"
env-cmd go run main.go
