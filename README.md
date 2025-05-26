# c8y-tedge

[go-c8y-cli](https://goc8ycli.netlify.app/) extension to add some useful command to help manage thin-edge.io with Cumulocity IoT.

## What is included?

**Note**

Use ✅ or 🔲 indicates if the extension includes the given functionality or not.


|Type|Included|Notes|
|----|:-:|-----|
|Commands|✅|Commands to manage thin-edge.io devices|
|Views|✅|thin-edge.io specific views|

## Install

The extension can be installed using the following command.

```sh
c8y extensions install thin-edge/c8y-tedge
```

Or if you have it already installed, then update to the latest version using:

```sh
c8y extensions update tedge
```

## Examples

### Start a demo container

Start a demo container using a randomly generated device name:

```sh
c8y tedge demo start
```

Or you can specify the device name to be used:

```
c8y tedge demo start tedge0001
```

Or you can specify the auth-type you want to use (e.g. certificate (default) or basic):

```sh
c8y tedge demo start tedge0001 --auth-type basic

c8y tedge demo start tedge0001 --auth-type certificate
```

### Start a tedge-container-bundle container

The [tedge-container-bundle](https://github.com/thin-edge/tedge-container-bundle) is a lighter weight containerized version of thin-edge.io where it can be used to deploy on devices which are already running a container engine and you don't want to install thin-edge.io on the host.

The `container-bundle` subcommands provide an easy way to start/stop/list instances of the tedge-container-bundle on your machine so that you can explore the functionality provided by it.

Start a [tedge-container-bundle](https://github.com/thin-edge/tedge-container-bundle) container that can be used to managed other containers. See the project for more details.

```sh
# Start with an auto generated name
c8y tedge container-bundle start

# Start and use a given name
c8y tedge container-bundle start tedge12345

# Start and use basic auth credentials
c8y tedge container-bundle start tedge12345 --auth-type basic

# Start but don't publish any of the ports on the host
c8y tedge container-bundle start tedge12345 --no-ports

# Start and publish ports to randomly assigned ports on the host
c8y tedge container-bundle start tedge12345 --publish-all
```

Then you can list and then stop the instances using (note: only the instances which were started with the c8y-tedge extension will be shown)

```sh
# list existing instances
c8y tedge container-bundle list

# stop / remove an instances
c8y tedge container-bundle stop tedge12345
```

### Bootstrap device via ssh

Bootstrap a thin-edge.io enable device using SSH.

```sh
c8y tedge bootstrap root@raspberrypi3-64.local
```

The bootstrapping processes does:

* Create the device certificate (if required)
* Fetch public device certificate and upload it to Cumulocity IoT (private key does not leave the device)
* Open the device in the Cumulocity IoT Device Management application
