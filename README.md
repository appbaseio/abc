
# ABC

1. [Intro](#intro)
2. [Installation](#installation)
	1. [Basic Installation](#basic-installation)
	2. [Using Docker](#using-docker)
3. [Build Variants](#build-variants)
4. [Features](#features)
	1. [Appbase features](#appbase-features)
	2. [Importer features](#importer-features)
5. [ABC Resources](#abc-resources)
	1. [Contributing to ABC](#contributing-to-abc)
	2. [Licensing](#licensing)

<a name="intro"></a>
## 1. Intro

ABC is a command-line client for appbase.io with nifty features. The paid version of this allows import data into Appbase as well. 

It consists of two parts. 

1. Appbase module
2. Import module (paid version)

To get the list of all commands supported by ABC, use -

```sh
abc --help
```

<a name="installation"></a>
## 2. Installation

ABC can be installed and used via the traditional `go build` or using a Docker image.

<a name="basic-installation"></a>
### 2.1 Basic installation

You can install ABC by building it locally and then moving the executable to anywhere you like. 

To build it, you require **Go 1.8** insalled on your system. 

```sh
go get github.com/appbaseio-confidential/abc
cd $GOPATH/src/github.com/appbaseio-confidential/abc
go build -tags 'oss' ./cmd/abc/...
./abc --help
```

Note - You might be wondering what is the tag `oss` doing there. That's covered in the section [Build Variants](#build-variants).

<a name="using-docker"></a>
### 2.2 Using Docker

```sh
git clone https://github.com/appbaseio-confidential/abc
cd abc
docker build --build-arg ABC_BUILD=oss -t abc .
docker volume create --name abc
```

Volume is used to store abc config files across containers.
Now `abc` can be ran through Docker like in the following example which starts google login.  

```sh
docker run -i --rm -v abc:/root abc login google
```

Some more examples

```sh
docker run -i --rm -v abc:/root abc user
docker run -i --rm -v abc:/root abc apps
```

<a name="build-variants"></a>
## 3. Build Variants

The ABC project you see in this repository is not the complete project. Appbase.io works on a proprietary version of ABC using this project as the base.
Hence we use the tag 'oss' to specify that this is an open source build. 
If you are curious, we use the tag '!oss' to make our private builds. 


#### How to know build variant from the executable? 

If you are not sure which build of `abc` you are using, you can run `abc --help` and take note of the value under the version header. 

For open source build, you will see

```
VERSION
  oss
```

For the proprietary builds, you will see 

```
VERSION
  proprietary
```

<a name="features"></a>
## 4. Features

ABC's features can be broadly categorized into 2 components. 

1. Appbase features
2. Importer features

<a name="appbase-features"></a>
### 4.1 Appbase features

Appbase features allows you to control your appbase.io account using ABC. You can see them under the *Appbase* heading in the list of commands.

```sh
APPBASE
  login     login into appbase.io
  user      get user details
  apps      display user apps
  app       display app details
  create    create app
  delete    delete app
  import    import data from various sources into appbase app
```

You can look over help for each of these commands using the `--help` switch. 
Alternatively we have detailed docs for them at [docs/appbase folder](docs/appbase).


```sh
abc login --help
```

#### Example

```sh
# display all commands
abc
# login into system
abc login google
# get user details
abc user
# get list of apps
abc apps
# get details of an app
abc app MyAppName
# delete that app
abc delete MyAppName
# create it again
abc create MyAppName
# view its metrics. It will be 0 as it is a new app
# here we are using AppID. We can use AppName too.
abc app -m 2489
```

<a name="importer-features"></a>
### 4.2 Importer features

ABC allows the user to configure a number of data adaptors as sources or sinks. These can be databases, files or other resources. Data is read from the sources, converted into a message format, and then send down to the sink where the message is converted into a writable format for its destination. The user can also create data transformations in JavaScript which can sit between the source and sink and manipulate or filter the message flow.

Adaptors may be able to track changes as they happen in source data. This "tail" capability allows a ABC to stay running and keep the sinks in sync.
For more details on adaptors, see **ABC pro website**.

<a name="abc-resources"></a>
## 5. ABC Resources

Checkout the [docs folder](docs/) for details on some ABC commands and topics.

<a name="contributing-to-abc"></a>
### 5.1 Contributing to ABC

Want to help out with ABC? Great! There are instructions to get you started [here](CONTRIBUTING.md).

<a name="licensing"></a>
### 5.2 Licensing

ABC is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for full license text.

