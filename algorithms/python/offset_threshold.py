#!/usr/bin/env python

import argparse
import json
import sys
import os


def applyAlgorithm(inputs, params):
    """
        This algorithm requires the following:
        - inputs["current"]: The current set of time series values. This is a dictionary with the time series labels as the key.
        - inputs["previous"]: The previous set of time series values. This is a dictionary with the time series labels as the key.
        - params["threshold"]: The permissible threshold between old and new entries

        The algorithm will the compare the non-empty labels as follows:
        (abs(current - previous)/current)*100 < threshold
    """
    output = {
        "violations": [],
        "environment": os.environ['ENVIRONMENT'],
        "error": None
    }
    current = inputs.get("current")
    previous = inputs.get("previous")
    threshold = params.get("threshold")
    if current == None or previous == None or threshold == None:
        output["error"] = "Inputs missing. Please ensure both current and previous keys exist and the threshold is defined in the params."
        printOutput(output)
        sys.exit(1)
    for series, value in current.items():
        if series in previous:
            # Check threshold
            # Everything is a string so convert values to float first
            delta = abs(float(value) - float(previous[series]))
            if ((delta/float(value)) * 100 > float(threshold)):
                output["violations"].append(series)
    printOutput(output)
    if len(output["violations"]) > 0:
        sys.exit(2)

def printOutput(output):
    print(json.dumps(output))

def main(args):
    with open(args.inputs, "r") as inputsF:
        inputs = json.load(inputsF)
    with open(args.params, "r") as paramsF:
        params = json.load(paramsF)
    applyAlgorithm(inputs, params)

if __name__=="__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--inputs", required=True)
    parser.add_argument("--params", required=True)
    args = parser.parse_args()
    main(args)
