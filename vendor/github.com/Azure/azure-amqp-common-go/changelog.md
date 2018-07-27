# Change Log

## `v0.6.0`
- **Breaking Change** change the parse connection signature and make it more strict
- fix errors imports

## `v0.5.0`
- **Breaking Change** lock dependency to AMQP

## `v0.4.0`
- **Breaking Change** remove namespace from SAS provider and return struct rather than interface 

## `v0.3.2`
- Return error on retry. Was returning nil if not retryable.

## `v0.3.1`
- Fix missing defer on spans

## `v0.3.0`
- add opentracing support
- upgrade amqp to pull in the changes where close accepts context (breaking change)

## `v0.2.4`
- connection string keys are case insensitive 

## `v0.2.3`
- handle remove trailing slash from host

## `v0.2.2`
- handle connection string values which contain `=`

## `v0.2.1`
- parse connection strings using key / values rather than regex

## `v0.2.0`
- add file checkpoint persister

## `v0.1.0`
- initial release