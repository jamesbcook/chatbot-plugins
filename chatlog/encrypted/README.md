# Encrypted

* Log entires are encrypted with AES-GCM. The block size is determined by key assigned to the CHATBOT_LOG_KEY environmental variable.
* A random 12 byte nonce is used for each message, and put in the front of each entry, this is needed to decrypt the messages.
  * ```export CHATBOT_LOG_KEY=abcdefghijklmnopqrstuvwxyz123456```