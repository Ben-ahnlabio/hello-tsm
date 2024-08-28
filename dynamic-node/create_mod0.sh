#!/bin/bash
openssl ecparam -name P-256 -genkey -param_enc named_curve -outform DER -out private_mod0.key
openssl base64 -A -in private_mod0.key; echo
openssl ec -inform DER -in private_mod0.key -pubout -outform DER -out public_mod0.key
openssl base64 -A -in public_mod0.key; echo
