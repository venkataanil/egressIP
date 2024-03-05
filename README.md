# EgressIP scale testing workload

This repo consists of workload configuration for scale testing EgressIPs.

## image
image folder consists of customized nginx server which adds source IP address to the return packet. Corresponding change in image/nginx/nginx.conf.template

http {
   server {
       location / {
          return 200 $remote_addr;
        }
  }
}

image/nginx/Containerfile provides LISTEN_PORT for user to specify as input in Pod so that nginx server can listen on this port. OCP worker on AWS platform opens only ports from 9001, so we can use these ports for serving requets

quay.io/vkommadi/nginx-custom:latest already built with this configuration.

## kube_burner_eip_simulate.go
This simulates how kube-burner-ocp and kube-burner interacts to generate EIP object template while running the job iterations.
kube_burner_ocp() will generate ip addresses from CIDR by excluding nodeips and then set them as environment variables.
kube_burner() will read this environment variable and generate IP addresses for each EIP object template
 
Pass --iterations and --addresses-per-iteration as cli arguments to kube_burner_eip_simulate.go.

```console
[vkommadi@fedora temp]$ go run kube_burner_eip_simulate.go --iterations 2 --addresses-per-iteration 2 --exclude-addresses "192.168.1.0 192.168.1.1"
EgressIPs:
 - 192.168.1.4
 - 192.168.1.5

EgressIPs:
 - 192.168.1.6
 - 192.168.1.7
```

## go_eips.go
This is a test file to simulate how we can generate EIP addresses in the workload object template.
Config params -
eip_cidr -  User provides CIDR to allocate addresses for EIP (kube-burner-ocp can read this from the node configuration).
nodeips -  User also provides IPs to exclude from this CIDR (this will be node IP addrsses which kube-burner can get from OCP cluster).
numIterations - total job iterations 
addrPerIteration - ip addresses per iteration. We will use 1 or 2 EIP ip addresses per EIP object

From the provided CIDR, test will generate IP addresses incrementaly for the objects which are part of iterations (excluding the node IPs) and add these addresses to object template file.

Sample output
```console
[vkommadi@fedora temp]$ go run go_eips.go 
EgressIPs:
 - 192.168.1.5
 - 192.168.1.6

EgressIPs:
 - 192.168.1.7
 - 192.168.1.8
```
This code will be added to kube-burner-ocp egressip workload. kube-burner-ocp will only generate ip addresses (exluding node ips) and then pass them to kube-burner which will add them to the template.

Note: There is alternate way of passing all the params (eip_cidr, nodeips) to kube-burner by kube-burner-ocp and kube-burner generate the ip addresses. But need to think a way for kube-burner to provide alternate IPs for node IP.

