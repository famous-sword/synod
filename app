#!/usr/bin/env bash

# shellcheck disable=SC1009
cat << 'EOF'
   _____                       _
  / ____|                     | |
 | (___  _   _ _ __   ___   __| |
  \___ \| | | | '_ \ / _ \ / _` |
  ____) | |_| | | | | (_) | (_| |
 |_____/ \__, |_| |_|\___/ \__,_|
          __/ |
         |___/
EOF

# shellcheck disable=SC1073
if [ $# -lt 2 ]; then
  echo "usage: app run command"
  echo "- api"
  echo "- storage"
  exit 1
fi

go run main.go "$1" "$2"
