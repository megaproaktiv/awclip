# awclip

Local caching of aws cli commands

- Store command results in `.awsclip` as md5hash of parameters
- Store metadata as hash-md.json

Use it as aws cli alias:

If awclip binary is in /usr/local/bin:

`alias aws='/usr/local/bin/awclip'`

## Limits

- No time to live implemented, delete `.awclip\*` for refresh
- You have to specify region and output
- Only one account, you have to clean `.awclip\*` for account switch
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

## run certain api calls in awclip

The following aws cli command are used by prowler and are beeing performed by awclip.

- `aws ec2 describe-instances --query "Reservations[*].Instances[*].[InstanceId]" --output text`
- `aws ec2 describe-regions --query "Regions[].RegionName" --output text`
- `aws sts get-caller-identity`

If a call is recognized as a supported call, the metadata says: `"Provider":"go"`.
The aws cli calls have `"Provider":"python"`

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

### current
  ## ADD
  - jmespath query. Because map ordering is not guarenteed, a key array is needed. This also fixes the ordering bug of aws cli text python

### v0.1.8
  ## ADD
  - add iam list-users
  - add prefetch lambda list function (prowler 762)
  - baxch script tests
  ## CHANGE
  - change id to struct parameter based, not cmd line based
### v0.1.6
  - implement iam list-user-policies with additional parameters
### v0.1.4  
- implement api calls with specific query in program
### v0.1.0
- reads command line
- calls aws with (python) aws cli
- writes metadata
- does not cache "generate-credential-report"

