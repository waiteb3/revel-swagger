#!/bin/bash

set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cat $DIR/swagger.yml | ruby -ryaml -rjson -e 'puts JSON.pretty_generate(YAML::load(ARGF.read))' > $DIR/swagger.json
