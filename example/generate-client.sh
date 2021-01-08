#!/bin/bash

for l in go javascript php; do
  docker run --rm -v "$(pwd):/go-work" swaggerapi/swagger-codegen-cli generate \
    -i "/go-work/rest.swagger.json" \
    -l "$l" \
    -o "/go-work/clients/$l"
done

docker run --rm -v "$(pwd):/go-work" swaggerapi/swagger-codegen-cli langs

echo "See https://github.com/swagger-api/swagger-codegen for more info"