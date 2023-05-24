#!/bin/sh

export PATH="$HOME/.serverless/bin:$PATH"

# BASELINE
cd ./../hotel_baseline/

echo "============  Build baseline handlers ============ ";
./build_handlers.sh handlers/

echo "============  Deploy baseline ============ ";
sls deploy 

# SIMPLE
cd ./../hotel_baseline_simple/

echo "============  Build simple handlers ============ ";
./build_handlers.sh handlers/

echo "============  Deploy simple ============ ";
sls deploy 


# NUBES
cd ./../hotel/

echo "============  Build nubes handlers ============ ";
./build_handlers.sh generated/

echo "============  Deploy nubes ============ ";
sls deploy