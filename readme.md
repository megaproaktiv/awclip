# awclip

Local caching of aws cli commands

- Store command results in `.awsclip` as md5hash of parameters
- Store metadata as hash-md.json

Use it as aws cli alias:

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

- copy awclip executable to your local filesystem
    - see https://github.com/megaproaktiv/awclip/tags
    - e.g. to `/usr/local/bin/awclip`

- create an alias
    ```bash
    alias aws=`/usr/local/bin/awclip`
    ```
    

- create a local `.awclip` directory

    ```bash
    mkdir .awclip
    ```
    

- or clean existing directory

    ```bash
    rm .awclip/*
    ```

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

## Working with prowler

Change `include/awscli_detector`

Change

```bash
if [ ! -z $(which aws) ]; then
  AWSCLI=$(which aws)
elif [ ! -z $(type -p aws) ]; then
  AWSCLI=$(type -p aws)
else
  echo -e "\n$RED ERROR!$NORMAL AWS-CLI (aws command) not found. Make sure it is installed correctly and in your \$PATH\n"
  EXITCODE=1
  exit $EXITCODE
fi
```

to

```bash
AWSCLI=./awclip
```

if you copy awclip into the same directory.

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

