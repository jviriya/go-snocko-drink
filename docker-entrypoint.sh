#!/bin/sh
chmod ug+x ./csr
eval $(./csr)
echo $CMD
exec ./go-pentor-bank.bin $CMD
