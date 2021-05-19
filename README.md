# NSXT Tag Controller

This controllers goal is to read labels from the `tkc` clusters and propogate them to the nsxt segment that is associated with the cluster as tags/scope so that they can be leveraged for grouping and firewalls.


## Architecture

 ![arch](/images/nsxt-tag-controller.jpg)


## Required Environment Variables

* `NSXT_HOST` - the hostname or ip for nsxt manager
* `NSXT_USERNAME` - NSXT user , must have priviliges to overwrite protected objects and update tags on segments
* `NSXT_PASSWORD` - Password for the nsxt user

## Deployment

**This will need to be deployed as a full cluster admin on the supervisor cluster.** 

1. temporarily give administrators cluster admin access to the supervisor
   1. ssh to vcenter and run `/usr/lib/vmware-wcp/decryptK8Pwd.py` 
   2. ssh to the ip that is output with the password that is output.
   3. `kubectl apply -f https://gist.githubusercontent.com/warroyo/9984a4e7ec1ee667153613153c8670ea/raw/58271b688583bd1f5c4feeecfeec014913d8277a/override-rbac.yml`

2. clone this repo to your desktop
3. cp `config/manager/env-sample.txt` `config/manager/env.txt`
4. update the values in the `env.txt` file
5. pull down the latest version of this image `warroyo90/nsxt-tag-controller:<version>` and move into your local repo or use a proxy cache
6. `export IMG=<path to the image>`
7. `make deploy`
8. validate its running
9. remove the above role binding `kubectl delete -f https://gist.githubusercontent.com/warroyo/9984a4e7ec1ee667153613153c8670ea/raw/58271b688583bd1f5c4feeecfeec014913d8277a/override-rbac.yml`


## Usage

after deploying the controller  it will watch for changes on `tkc` objects. it will only update tags in nsxt if the label has the prefix of `policytag/`

1. edit a tkc and add a new label with a prefix of `policytag/`

ex.
```yaml
labels:
  policytag/hello: world
```

2. you should see tags/scopes updated on the segment in nsxt

## Building

### Build docker image

1. `cp .netrc-sample .netrc` and update with creds
2. `export IMG=<your-image-name>`
3. `make docker-build`


### Build locally

1. `make`


