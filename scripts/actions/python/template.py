#!/usr/bin/env python

import argparse
import json

def applyAction(inputs, params):
    print(inputs)
    print(params)

def main(args):
    with open(args.inputs, "r") as inputsF:
        inputs = json.load(inputsF)
    with open(args.params, "r") as paramsF:
        params = json.load(paramsF)
    applyAction(inputs, params)

if __name__=="__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--inputs", required=True)
    parser.add_argument("--params", required=True)
    args = parser.parse_args()
    main(args)
