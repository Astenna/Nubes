#!/bin/bash
set -euxo pipefail

## medium

# const UserCount = 2000
# const CitiesCount = 25
# const HotelsPerCity = 10
# const RoomsPerHotel = 10
# const ReservationsPerRoom = 20

GATEWAY_NUBES=https://qzod5szgz6hgwox3suztsgc4ze0mkefm.lambda-url.us-east-1.on.aws/
GATEWAY_BASELINE=https://6dy6oeoalo3apzaq2vqglb42rq0zcmlp.lambda-url.us-east-1.on.aws/
GATEWAY_SIMPLE=https://vek53sm3o6qtnn2gtdat6wvuza0qnyyv.lambda-url.us-east-1.on.aws/
DURATION=120
EXTRA_TIME=15
WRK_DURATION=$((${DURATION} + 2 * ${EXTRA_TIME}))
RATE=1000
RATE_SIMPLE=$(( $RATE / 4 ))
CONFIG=small-r${RATE}-${DURATION}s
WRK_ARGS="-c 100 -t 8"

mkdir -p result

# Align to even minutes - 15s
align() {
    sleep $(( ($(date +%s) / ${DURATION} + 1) * ${DURATION} - $(date +%s) - ${EXTRA_TIME} ))
}

align
wrk2 ${WRK_ARGS} --latency -R "${RATE}" -d "${WRK_DURATION}s" -s hotel.lua ${GATEWAY_NUBES} | tee "result/wrk-nubes-${CONFIG}.log"
sleep 60
./hotel.py --experiment "nubes.toml" --config "${CONFIG}" --duration "${DURATION}"

align
wrk2 ${WRK_ARGS} --latency -R "${RATE}" -d "${WRK_DURATION}s" -s hotel_baseline.lua ${GATEWAY_BASELINE} | tee "result/wrk-baseline-${CONFIG}.log"
sleep 60
./hotel.py --experiment "baseline.toml" --config "${CONFIG}" --duration "${DURATION}"

align
wrk2 ${WRK_ARGS} --latency -R "${RATE_SIMPLE}" -d "${WRK_DURATION}s" -s hotel_baseline_simple.lua ${GATEWAY_SIMPLE} | tee "result/wrk-simple-${CONFIG}.log"
sleep 60
./hotel.py --experiment "baseline_simple.toml" --config "${CONFIG}" --duration "${DURATION}"
