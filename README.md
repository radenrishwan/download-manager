# Parrarel Download Manager
This project is a demo how to download files from internet with parrarel method.

## How this is work?
HTTP has a header called `Range` that can be used to download a file from a specific byte to another byte. So we can download a file in parrarel by splitting it to some parts and download each part in a thread. after all parts downloaded, we can merge them to a single file.

## How this is make download faster?
Sometimes, the server that we download from it, has a limit for download speed. So we can download a file in parrarel to bypass this limit.

## How to use?
I'm already write a simple code on main.go file