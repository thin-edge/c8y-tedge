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
c8y extension install thin-edge/c8y-tedge
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

### Bootstrap device via ssh

Bootstrap a thin-edge.io enable device using SSH.

```sh
c8y tedge bootstrap root@raspberrypi3-64.local
```

The bootstrapping processes does:

* Create the device certificate (if required)
* Fetch public device certificate and upload it to Cumulocity IoT (private key does not leave the device)
* Open the device in the Cumulocity IoT Device Management application
