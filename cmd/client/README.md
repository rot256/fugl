# Usage

You specify an operation using the --operation flag.
See details below for the different operations supported.

## Before you start

Before you start you should have a PGP key-pair,
the public key should be configured on the server and
the private key available on client (if you wish to create new proofs).
Furthermore you should create a directory for storing the proof chain (default is ./store)

## Creating

You can create new proofs using the create operation.
Before doing so you must create a canary "manifest", which is a template for creating new canaries.
An example of such a manifest can be found in the current directory (see manifest.toml).

After editing the manifest, a new proof can be created like this:

```
~> ./client --operation=create --private-key=./private.pgp
Private key encrypted, please enter passphrase:
Wrote new proof to: temp.proof
```

## Verifying

For verifying the validity of proofs, you must use the "verify" operation.
It verifies the PGP signature as well as fields of the canary in the metadata section.

An example follows (using the newly created proof):

```
~> ./client --operation=verify --public-key=./public.pgp
Author: Test author
Expires: 2017-02-21T11:27:54+01:00
Description:
# Test canary

You can:

* Explain the purpose of the canary
* Maintain a human readable canary
...
```

## Pushing

Proofs are added to the server by pushing.
Pushing is really a fancy word for issuing a POST request, after which the server will validate the proof
(in the same way as the "verify" operation on the client).

A silly example using a server on localhost follows (using the default value for the "proof" parameter):

```
~> ./client --operation=push --address=http://127.0.0.1:8080/submit
Successfully pushed new proof to server
```

## Pulling

The client "pulls" the newest proof from the server into a temporary file using a GET request.
You then use the verify operation to verify its validity.

Here an example of using the "pull" operation:

```
~> ./client --operation=pull --address=http://127.0.0.1:8080/latest
Saved to: temp.proof
```
