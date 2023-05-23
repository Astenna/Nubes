#!/bin/bash
set -euxo pipefail

## medium

# const UserCount = 2000
# const CitiesCount = 25
# const HotelsPerCity = 10
# const RoomsPerHotel = 10
# const ReservationsPerRoom = 20

GATEWAY_NUBES=AAAA
GATEWAY_BASELINE=BBBB
DURATION=330
RATE=10
CONFIG=medium-r${RATE}-${DURATION}s

mkdir -p result

# Align to even minutes - 15s
align() {
    sleep $(( ($(date +%s) / 300 + 1) * 300 - $(date +%s) - 15 ))
}

align
wrk2 --latency -R "${RATE}" -d "${DURATION}s" -s hotel.lua ${GATEWAY_NUBES} | tee "result/wrk-nubes-${CONFIG}.log"
sleep 60
./hotel.py --experiment "nubes.toml" --config "${CONFIG}" --duration "${DURATION}"

align
wrk2 --latency -R "${RATE}" -d "${DURATION}s" -s hotel_baseline.lua ${GATEWAY_BASELINE} | tee "result/wrk-baseline-${CONFIG}.log"
sleep 60
./hotel.py --experiment "baseline.toml" --config "${CONFIG}" --duration "${DURATION}"
