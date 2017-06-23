[![Build Status](https://travis-ci.org/compose/transporter.svg?branch=master)](https://travis-ci.org/compose/transporter) [![Go Report Card](https://goreportcard.com/badge/github.com/compose/transporter)](https://goreportcard.com/report/github.com/compose/transporter) [![codecov](https://codecov.io/gh/compose/transporter/branch/master/graph/badge.svg)](https://codecov.io/gh/compose/transporter) [![Docker Repository on Quay](https://quay.io/repository/compose/transporter/status "Docker Repository on Quay")](https://quay.io/repository/compose/transporter) [![Gitter](https://img.shields.io/gitter/room/nwjs/nw.js.svg)](https://gitter.im/compose-transporter/Lobby)


# ABC

1. [Intro](#intro)
2. [Installation](#installation)
	1. [Basic Installation](#basic-installation)
	2. [Using Docker](#using-docker)
3. [Build Variants](#build-variants)
4. [Features](#features)
	1. [Appbase features](#appbase-features)
	2. [Importer features](#importer-features)


<a name="intro"></a>
## 1. Intro

ABC is a command-line client for appbase.io with nifty features to do data sync from on store to another.

It consists of two parts. 

1. Appbase module
2. Import module

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
```

You can look over help for each of these commands using the `--help` switch. 

```sh
abc login --help
```


<a name="importer-features"></a>
### 4.2 Importer features

Transporter allows the user to configure a number of data adaptors as sources or sinks. These can be databases, files or other resources. Data is read from the sources, converted into a message format, and then send down to the sink where the message is converted into a writable format for its destination. The user can also create data transformations in JavaScript which can sit between the source and sink and manipulate or filter the message flow.

Adaptors may be able to track changes as they happen in source data. This "tail" capability allows a Transporter to stay running and keep the sinks in sync.

#### BETA Feature

As of release `v0.4.0`, transporter contains support for being able to resume operations
after being stopped. The feature is disabled by default and can be enabled with the following:

```
source = mongodb({"uri": "mongo://localhost:27017/source_db"})
sink = mongodb({"uri": "mongo://localhost:27017/sink_db"})
t.Config({"log_dir":"/data/transporter"})
  .Source("source", source)
  .Save("sink", sink)
```

When using the above pipeline, all messages will be appended to a commit log and 
successful processing of a message is handled via consumer/sink offset tracking.

Below is a list of each adaptor and its support of the feature:

```
+---------------+-------------+----------------+
|    adaptor    | read resume | write tracking |
+---------------+-------------+----------------+
| elasticsearch |             |       X        | 
|     file      |             |       X        | 
|    mongodb    |      X      |       X        |
|     mssql     |             |       N/A      | 
|  postgresql   |             |       X        | 
|   rabbitmq    |      X      |                | 
|   rethinkdb   |             |       X        | 
+---------------+-------------+----------------+
```

#### Adaptors

Each adaptor has its own README page with details on configuration and capabilities.

* [elasticsearch](./adaptor/elasticsearch)
* [file](./adaptor/file)
* [mongodb](./adaptor/mongodb)
* mssql
* [postgresql](./adaptor/postgres)
* [rabbitmq](./adaptor/rabbitmq)
* [rethinkdb](./adaptor/rethinkdb)


#### Native Functions

Each native function can be used as part of a `Transform` step in the pipeline.

* [goja](./function/gojajs)
* [omit](./function/omit)
* [otto](./function/ottojs)
* [pick](./function/pick)
* [pretty](./function/pretty)
* [rename](./function/rename)
* [skip](./function/skip)

#### Commands

The importer module has the following commands.

```
run       run pipeline loaded from a file
test      display the compiled nodes without starting a pipeline
about     show information about available adaptors
init      initialize a config and pipeline file based from provided adaptors
xlog      manage the commit log
offset    manage the offset for sinks
```

Details have been covered in the Wiki page : [Importer Commands](https://github.com/appbaseio-confidential/abc/wiki/Importer-Commands). 


## Building guides

[macOS](https://github.com/appbaseio-confidential/abc/blob/master/READMEMACOS.md)
[Windows](https://github.com/appbaseio-confidential/abc/blob/master/READMEWINDOWS.md)
[Vagrant](https://github.com/appbaseio-confidential/abc/blob/master/READMEVAGRANT.md)


## ABC Resources

* [ABC Wiki](https://github.com/appbaseio-confidential/abc/wiki)


## Contributing to ABC

Want to help out with ABC? Great! There are instructions to get you
started [here](CONTRIBUTING.md).


## Licensing

ABC is licensed under the New BSD License. See [LICENSE](LICENSE) for full license text.

