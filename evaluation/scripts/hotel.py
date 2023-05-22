#!/usr/bin/env python3
import boto3
import numpy as np
import os
import csv
import toml
from argparse import ArgumentParser
from datetime import datetime, timedelta

log_client = boto3.client("logs")
cloudwatch = boto3.resource("cloudwatch")


def get_metric(cloudwatch, metric, fname, duration, unit=None):
    end_time = datetime.utcnow()
    metric = cloudwatch.Metric("AWS/Lambda", metric)
    response = metric.get_statistics(
        Dimensions=[{"Name": "FunctionName", "Value": fname}] if fname is not None  else [],
        Statistics=["Sum", "Average", "Maximum"],
        ExtendedStatistics=["p01", "p10", "p50", "p90", "p99"],
        StartTime=end_time - timedelta(minutes=duration + 1),
        EndTime=end_time + timedelta(minutes=1),
        Period=120,
        Unit=unit,
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
                percs["p01"],
                percs["p10"],
                percs["p50"],
                percs["p90"],
                percs["p99"],
            ]
        )

    return res


def dump_metric(experiment, config, task, metric, data, path="result/"):
    if len(data) == 0:
        print(f"Skipping {metric}: no data.")
        return

    try:
        os.mkdir(path)
    except:
        pass

    with open(
        f"{path}/{experiment}-{task}-{config}-{metric.lower()}.csv", "w"
    ) as f:
        w = csv.writer(f)
        w.writerow(["experiment", "task", "config", "timestamp", "avg", "sum", "max", "p01", "p10", "p50", "p90", "p99"])
        for r in data:
            w.writerow([experiment, task, config, *r])

    print(f"{metric} median:\t{np.mean([x[-2] for x in data])}")


def main():
    parser = ArgumentParser()
    parser.add_argument("--experiment", required=True)
    parser.add_argument("--config", required=True)
    parser.add_argument("--duration", required=False, default = 60)
    args = parser.parse_args()


    with open(args.experiment, "r") as f:
        config = toml.load(f)

    name = config["name"]
    path = f"result/{name}"

    metrics = [
        ["Duration", "Milliseconds"],
        ["Invocations", "Count"],
        ["ConcurrentExecutions", "Count"],
        ["UrlRequestLatency", "Milliseconds"],
        ["UrlRequestCount", "Count"],
        ["Throttles", "Count"],
    ]

    if not os.path.exists(path):
        os.makedirs(path)

    for function in config["functions"]:
        fname = function["lambda"]
        task = function["task"]
        print(f"-- {name}-{task} ({fname})")
        for [metric, unit] in metrics:
            r = get_metric(cloudwatch, metric, fname, int(args.duration), unit)
            dump_metric(name, args.config, task, metric, r, path=path)

    

    print(f"-- {name}-total")
    for [metric, unit] in metrics:
        r = get_metric(cloudwatch, metric, None, int(args.duration), unit)
        dump_metric(name, args.config, "total", metric, r, path=path)


if __name__ == "__main__":
    main()
