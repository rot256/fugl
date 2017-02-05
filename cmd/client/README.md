# Usage

You specify an operation using the --operation flag.
See details below for the different operations supported.

## Before you start

Before you start you should have a PGP key-pair,
the public key should be configured on the server and
the private key available on client (if you wish to create new proofs).
Furthermore you should create a directory for storing the proof chain (default is ./store)

## Creating

You can create new proofs using the create operation,
the `previous` field will be populated using the newest proof in the store (unless the store is empty).
Here an example of creating a proof, with an expiry date 100 days (2400 hours) into the future.

```
~> ./client --operation=create --private-key=./private.pgp --expire=2400h
Private key encrypted, please enter passphrase:
Wrote new proof to: temp.proof
It is recommended that you add this to the local store
```

## Adding

When you have created (or pulled) a new proof you may want to add this to the store,
which is the directory storing the state of the proof chain.
This is used both when following chains from other people and after creating a new proof.
An example of using the "add" operation:

```
~> ./client --operation=add --input=./temp.proof --public-key=./public.pgp
New deadline: 2017-05-16T22:40:21+02:00
```

Of course this also validates the signature and previous field.

## Pushing

Proofs are added to the server by pushing.
Pushing is really a fancy word for issuing a post request, after which the server will validate the proof (in completely the same way as the "add" operation).

Here a silly example using a server on localhost:

```
~> ./client --operation=push --address=http://127.0.0.1:8080/submit
Successfully pushed new proof to server
```

## Pulling

This is essentially just a HTTP get request and can also be carried out using wget, curl or your favourite browser.
The client "pulls" the newest proof from the server into a temporary file (output),
which you can then validate and add to the store using "add" (see next section).

```
~> ./client --operation=pull --address=http://127.0.0.1:8080/latest
Saved to: temp.proof
```
