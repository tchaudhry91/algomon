#!/usr/bin/env python

import argparse
import requests
import json

def sendToTeams(inputs, params):
    '''sends alert to teams channel'''
    headers={"Content-Type": "application/json"}
    username = params.get("username")
    title = username + " " + params.get("name")
    message = {
        "check_output": inputs
    }
    payload = {
        "summary": title,
        "themeColor": "#FF0000",
        "sections": [{
            "activityTitle": title,
            "activitySubtitle": formatMessage(message),
        }],
    }
    webhook_url = params.get("webhook_url")
    res = requests.post(url=webhook_url, json=payload, headers=headers)
    if res.status_code != 200:
        print("Not able to send to teams")
        print(res.text)
    return res.text

def formatMessage(message):
    formatted = """```
    {}
    """.format(json.dumps(message, indent=4))
    return formatted

def applyAction(inputs, params):
    sendToTeams(inputs, params)

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
