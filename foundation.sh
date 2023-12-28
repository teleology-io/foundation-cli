#!/bin/sh 
# foundation cli

while [[ $# -gt 0 ]]; do
  case $1 in
    -k|--api-key)
      API_KEY="$2"
      shift # past argument
      shift # past value
      ;;
     -n|--name)
      NAME="$2"
      shift 
      shift
      ;;
    env)
      ARG="env"
      shift
      ;;
    config)
      ARG="config"
      shift
      ;;
    variable)
      ARG="variable"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

function fetch_environment() {
  curl --request GET \
    --url http://localhost:3000/test/v1/environment \
    --header "X-Api-Key: $API_KEY"
}

function fetch_configuration() {
  curl --request GET \
    --url http://localhost:3000/test/v1/configuration \
    --header "X-Api-Key: $API_KEY"
}

function fetch_variable() {
  curl --request POST \
    --url http://localhost:3000/test/v1/variable \
    --header "X-Api-Key: $API_KEY" \
    --data '{
  "name": "$NAME",
  "fallback_value": "test"
}'
}

case $ARG in
  env)
    fetch_environment
    shift
    ;;
  config)
    fetch_configuration
    shift
    ;;
  variable)
    fetch_variable
    shift
    ;;
  *)
    shift
    ;;
esac


  