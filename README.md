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
3. `make deploy`
4. remove the above role binding `kubectl delete -f https://gist.githubusercontent.com/warroyo/9984a4e7ec1ee667153613153c8670ea/raw/58271b688583bd1f5c4feeecfeec014913d8277a/override-rbac.yml`

## Building

### Build docker image

1. `cp .netrc-sample .netrc` and update with creds
2. `export IMG=<your-image-name>`
3. `make docker-build`


### Build locally

1. `make`


