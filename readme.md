# awclip

Local caching of aws cli commands

- Store command results in `.awsclip` as md5hash of parameters
- Store metadata as hash-md.json

User it as aws cli alias:

If awclip binary is in /usr/local/bin:

`alias aws='/usr/local/bin/awclip'`


## Status

Base functionality is working, no guarantee that it works on all services.

## Example

Time with aws cli:

```bash
time aws iam list-roles >/dev/null
aws iam list-roles > /dev/null  
0,36s user 0,17s system 39% cpu 1,324 total
```

Time varies from 1..2 seconds

1st time with awclip:

```bash
time aws iam list-roles >/dev/null
/Users/silberkopf/letsbuild/awclip/dist/awclip iam list-roles > /dev/null  
0,36s user 0,18s system 32% cpu 1,671 total
```

2nd time with awclip:

```bash
 time aws iam list-roles >/dev/null
/Users/silberkopf/letsbuild/awclip/dist/awclip iam list-roles > /dev/null 
 0,00s user 0,00s system 40% cpu 0,014 total
```

## Todo

- create `.awcli` automatically
- implement ttl (time to live), currently you have to clean `.awclip` yourself
- create storable configuration
- speed up region prefetch

## Version 



### v0.1.0
- reads command line
- calls aws with (python) aws cli
- writes metadata

