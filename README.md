## SQS Queues Proxy [![Build Status](https://travis-ci.org/lalvarezguillen/sqs_proxy.svg?branch=master)](https://travis-ci.org/lalvarezguillen/sqs_proxy) [![codecov](https://codecov.io/gh/lalvarezguillen/sqs_proxy/branch/master/graph/badge.svg)](https://codecov.io/gh/lalvarezguillen/sqs_proxy)

This is a work in progress, and a learning project.

This project's purpose is, given an SQS queue A, take all the messages it receives and send them to queues B, C, etc; deleting them from A in the process.

It is configurable through a JSON file, being able to listen to several queues, and set several recipient queues per source.

Currently it only works within a single AWS account.


### Sample config.

Assuming we have 2 source queues (A and B), and want to proxy all messages that A receives to queues C and D, and all messages that B receives to D, E and F; the config file would look like:

```json
{
    "proxyOps": [
        {
            "src": "https://sqs.region.amazonaws.com/accountid/A",
            "dest": [
                "https://sqs.region.amazonaws.com/accountid/C",
                "https://sqs.region.amazonaws.com/accountid/D"
            ],
            "interval": 20
        },
        {
            "src": "https://sqs.region.amazonaws.com/accountid/B",
            "dest": [
                "https://sqs.region.amazonaws.com/accountid/D",
                "https://sqs.region.amazonaws.com/accountid/E",
                "https://sqs.region.amazonaws.com/accountid/F"
            ],
            "interval": 35
        }
    ]
}
```

The attribute `interval` specifies the number of seconds to wait between polls to that particular source queue.

### AWS Credentials

This project relies on AWS SDK fetching your AWS Credentials. Look at [their docs](https://github.com/aws/aws-sdk-go) for more info


### Running the proxy

The proxy is called with a `--config`  flag that  points to the JSON config file to use.

```bash
go build .
./sqs_proxy --config ./config.json
```

### TODO

* Extend unittests
* Make it possible to proxy messages to the recipient queues in a round robin fashion. Consider a random recipient method as well.
* Make the long-polling time configurable
* Set up CI
* Consider storing some metadata on the messages proxied