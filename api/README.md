# API Plugin

This is plugin setups up a listener on port 55449 by default. This can be changed by setting the environment variable `CHATBOT_LISTEN_PORT`.

## Details

* Allow for incoming connections that will send data to a channel.
  * Incoming messages need to use the chatbot-external-api protobuf package.
  * There are example clients in the chatbot-external-api package, but if you want to use a different language the following steps are needed to be taken.
    * Generate a ED25519 key pair
    * Generate a curve25519 key pair
    * Exchange keys with the KeyExchange protobuf
    * Send Encrypted message to the server
    * Receive Encrypted message
    * Send finish packet
* The API only allows for public keys that have been added via the keybase chat to fully communicate.
* Packets follow the following layout
  * 4 bytes for the length of the message (Little Endian format)
  * 64 bytes for the signed output of the hashed encrypted message and nonce
  * 12 bytes for the nonce
  * The encrypted protobuf message

## Add Example

```
/api add 5364affd9d4d8596fd8662c02f30ced0ec486562ce97a4354c1137c22b9cec88
---------------------------------
Adding 5364affd9d4d8596fd8662c02f30ced0ec486562ce97a4354c1137c22b9cec88
```

## Info Example

```
/api info
---------------------------------
Server Info
Public 206.189.225.117:55449
Private 206.189.225.117:55449
Imported Keys
Key0: 5364affd9d4d8596fd8662c02f30ced0ec486562ce97a4354c1137c22b9cec88
```

## Delete Example

```
/api delete 5364affd9d4d8596fd8662c02f30ced0ec486562ce97a4354c1137c22b9cec88
---------------------------------
Deleting 5364affd9d4d8596fd8662c02f30ced0ec486562ce97a4354c1137c22b9cec88
```