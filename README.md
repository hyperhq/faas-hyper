FaaS-hyper
===========

This is a plugin to enable Hyper.sh as an [OpenFaaS](https://github.com/alexellis/faas) backend.

## Requirements

- A working Hyper.sh account with credential
- A Hyper.sh CLI with a working `hyper`

## Usage

### Get Started with Hyper.sh

Follow the [instruction](https://docs.hyper.sh/GettingStarted/index.html) to get `Access Key` and `Secret Key` of Hyper.sh account.

Go to [FIP](https://console.hyper.sh/fips) page to allocate a new FIP, get the FIP address `x.x.x.x`.

### Configure FaaS-hyper

Copy `docker-compose-tpl.yml` file to `docker-compose.yml`.

Edit the `docker-compose.yml` and input environment variables using the values above:

- `<HYPER_ACCESS_KEY>`: insert `Access Key`
- `<HYPER_SECRET_KEY>`: insert `Secret Key`
- `<FIP>`: insert `x.x.x.x`

### Deploy FaaS-hyper

Use hyper compose command `hyper compose up -d -p faas` to deploy FaaS-hyper and FaaS API Gateway.

Once deployed, you can use `http://x.x.x.x:8080/` to open the FaaS Gateway UI.

### Deploy a tester function

Use hyper pull command `hyper pull functions/nodeinfo` to pull the image from Docker Hub.

Click `CREATE NEW FUNCTION` on FaaS Gateway UI to deploy the function `functions/nodeinfo`:

![](https://camo.githubusercontent.com/72f71cb0b0f6cae1c84f5a40ad57b7a9e389d0b7/68747470733a2f2f7062732e7477696d672e636f6d2f6d656469612f44466b5575483158734141744e4a362e6a70673a6d656469756d)

### Invoke a tester function

`curl http://x.x.x.x:8080/function/nodeinfo`

### Remove the FaaS-hyper service

`hyper compose rm -p faas`
