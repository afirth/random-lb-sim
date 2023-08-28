# random lb sim

Simulates a fleet of workers receiving work in a random fashion, demonstrating the
inefficient resource usage when approaching low error rates.

An error is defined as a request to a worker which is already at capacity.
This is representative for workloads where the initial connection has high latency and the
response is time sensitive.

## Usage

See makefile or `go run . -h`

## Contributions

This code is a giant hack that I smashed together with no particular plan. If you find it useful feel free to fork

## Caveats

The code, particularly around stats, is quite a mess.

The autoscaler uses the overall average so on long runs it can spool up massive pools. It's not even proportional, so just make a good guess when you start, or fix it.

There are no tests - who knows if the code is correct. It answered my questions
