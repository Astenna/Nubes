#!/usr/bin/env python3
import boto3
import numpy as np
import os
import csv
from pprint import pprint
from argparse import ArgumentParser
from datetime import datetime, timedelta


log_client = boto3.client("logs")
cloudwatch = boto3.resource("cloudwatch")

def get_metric(cloudwatch, metric, fname, duration):
    end_time = datetime.utcnow()
    metric = cloudwatch.Metric("AWS/Lambda", metric)
    req = {"Name": "FunctionName", "Value": fname}
    response = metric.get_statistics(
        Dimensions=[req],
        Statistics=["Sum", "Average", "Maximum"],
        ExtendedStatistics=["p50", "p99"],
        StartTime=end_time - timedelta(minutes=duration + 1),
        EndTime=end_time + timedelta(minutes=1),
        Period=60,
        Unit="Milliseconds",
    )
    points = response["Datapoints"]
    points.sort(key=lambda x: x["Timestamp"])
    res = []
    for point in points:
        percs = point["ExtendedStatistics"]
        res.append(
            [
                point["Timestamp"],
                point["Average"],
                point["Sum"],
                point["Maximum"],
                percs["p50"],
                percs["p99"],
            ]
        )

    return res


def dump_metric(config, metric, data):
    if len(data) == 0:
        print(f"Skipping {metric}: no data.")
        return

    try:
        os.mkdir("result/")
    except:
        pass

    with open("result/{}-{}.csv".format(config, metric.lower()), "w") as f:
        w = csv.writer(f)
        w.writerow(["timestamp", "avg", "sum", "max", "p50", "p99"])
        w.writerows(data)

    print("=========================================")
    print("Median: {}".format(np.mean([x[-2] for x in data])))
    print("99 Percentile: {}".format(np.mean([x[-1] for x in data])))
    print("=========================================")


def main():
    parser = ArgumentParser()
    parser.add_argument("--name", required=True)
    parser.add_argument("--config", required=True)
    parser.add_argument("--duration", required=False, default = 60)
    args = parser.parse_args()

    metrics = ["Duration", "Invocations", "ConcurrentExecutions"]

    for metric in metrics:
        r = get_metric(cloudwatch, metric, args.name, int(args.duration))
        dump_metric(args.config, metric, r)


if __name__ == "__main__":
    main()
