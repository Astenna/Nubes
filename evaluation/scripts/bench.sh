#!/bin/bash
set -uxo pipefail

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
RATE_FORTH=$(( $RATE / 4 ))
CONFIG=large-r${RATE}-${DURATION}s
WRK_ARGS="-c 100 -t 8"

mkdir -p result

# Align to even minutes - 15s
align() {
    sleep $(( ($(date +%s) / ${DURATION} + 1) * ${DURATION} - $(date +%s) - ${EXTRA_TIME} ))
}

# Retry function
retry_command() {
    local command="$1"
    local max_attempts="$2"
    local out_file="$3"
    local attempt=1

    # Loop until successful execution or reaching the maximum number of attempts
    while [ $attempt -le $max_attempts ]; do
        # Execute the command in a subshell to isolate it from the main script
        $command 2>&1 | tee $out_file

        # Check the exit status of the command
        if [ $? -eq 0 ]; then
            return 0  # Success
        else
            echo "Command failed. Retrying..."
            sleep 1
            ((attempt++))
        fi
    done

    echo "Maximum number of attempts reached. Exiting."
    return 1  # Failure
}


### NUBES

align
retry_command "wrk2 ${WRK_ARGS} --latency -R ${RATE} -d ${WRK_DURATION}s -s hotel-full.lua ${GATEWAY_NUBES}" 15 "result/wrk-nubes-${CONFIG}.log"
sleep 60
./hotel.py --experiment "nubes-full.toml" --config "${CONFIG}" --duration "${DURATION}"

align
retry_command "wrk2 ${WRK_ARGS} --latency -R "${RATE_FORTH}" -d ${WRK_DURATION}s -s hotel-s1.lua ${GATEWAY_NUBES}" 15 "result/wrk-nubes-s1-${CONFIG}.log"
sleep 60
./hotel.py --experiment "nubes-s1.toml" --config "${CONFIG}" --duration "${DURATION}"

align
retry_command "wrk2 ${WRK_ARGS} --latency -R "${RATE_FORTH}" -d ${WRK_DURATION}s -s hotel-s2.lua ${GATEWAY_NUBES}" 15 "result/wrk-nubes-s2-${CONFIG}.log"
sleep 60
./hotel.py --experiment "nubes-s2.toml" --config "${CONFIG}" --duration "${DURATION}"

### BASELINE

align
retry_command "wrk2 ${WRK_ARGS} --latency -R ${RATE} -d ${WRK_DURATION}s -s hotel_baseline-full.lua ${GATEWAY_BASELINE}" 15 "result/wrk-baseline-${CONFIG}.log"
sleep 60
./hotel.py --experiment "baseline-full.toml" --config "${CONFIG}" --duration "${DURATION}"

### BASELINE SIMPLE

align
retry_command "wrk2 ${WRK_ARGS} --latency -R ${RATE} -d ${WRK_DURATION}s -s hotel_baseline_simple-full.lua ${GATEWAY_SIMPLE}" 15 "result/wrk-mixed-${CONFIG}.log"
sleep 60
./hotel.py --experiment "baseline_simple-full.toml" --config "${CONFIG}" --duration "${DURATION}"
