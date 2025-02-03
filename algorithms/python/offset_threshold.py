#!/usr/bin/env python
import argparse
import json

def main(args):
    inputs = json.load(args.inputs)
    params = json.load(args.params)
    print(inputs)
    print(params)

if __name__=="__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("inputs", required=True)
    parser.add_argument("params", required=True)
    args = parser.parse_args()
    main(args)
