#!/usr/bin/env bash

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

ulimit -n 102400

if [[ "$#" = 0 ]]; then
  echo "usage: start.sh command"
  exit 1
fi

if [[ "$1" = "api" ]]; then
    startAPI()
else
    startStorage()
fi

startAPIServer() {
    export SYNOD_APP_ID=0
    export SYNOD_API_ADDR=":5555"
    go run main.go run api
}

startDataServer() {
    if [[ ! "$1" =~ ^-?[0-9]+$ ]]; then
      echo "you must input a number"
      exit 1
    fi

    if [ ! -d var/disk/"$1" ]; then
        mkdir "var/disk/$1"
    fi

    if [ ! -d var/temp/"$1" ]; then
        mkdir "var/temp/$1"
    fi

    export SYNOD_APP_ID="$1"
    export SYNOD_API_ADDR=":5555"
    export SYNOD_DATA_ADDR=":556$1"
    export SYNOD_DATA_DIR="var/disk/$1"
    export SYNOD_TEMP_DIR="var/temp/$1"


    go run main.go run storage
}
