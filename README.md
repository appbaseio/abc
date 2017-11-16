
[![latest version](https://img.shields.io/badge/version-0.4.1-yellow.svg)](https://github.com/appbaseio/abc/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/appbaseio/abc)](https://goreportcard.com/report/github.com/appbaseio/abc)
[![Travis branch](https://img.shields.io/travis/appbaseio/abc/dev.svg)]()
[![license](https://img.shields.io/github/license/appbaseio/abc.svg)]()


![abc banner image](https://user-images.githubusercontent.com/4047597/29240054-14e0e19a-7f7b-11e7-898b-ba6bad756b1d.png)

# ABC

ABC is a command-line too to interact with appbase.io. 
It can also serve as a swiss army knife to import data from any popular data source (Postgres, SQL, Mongo) to ElasticSearch. 
This feature works with minimum configuration and is totally automatic. 
In certain sources like Postgres and Mongo, you can even keep the database and ElasticSearch cluster in sync such that any change from source gets added in destination as well.


1. [Intro](#intro)
2. [Key Benefits](#key-benefits)
3. [Getting Started](#getting-started)
4. [Features](#features)
	1. [Appbase features](#appbase-features)
	2. [Importer features](#importer-features)
5. [Development setup](#development)
	1. [Local Setup](#local-setup)
	2. [Docker Setup](#docker-setup)
	3. [Build Variants](#build-variants)
6. [ABC Resources](#abc-resources)
	1. [Contributing to ABC](#contributing-to-abc)
	2. [Licensing](#licensing)

<a name="intro"></a>
## 1. Intro

ABC consists of two parts. 

1. Appbase module
2. Import module (closed source)

To get the list of all commands supported by ABC, use -

```sh
abc --help
```


<a name="key-benefits"></a>
## 2. Key Benefits

ABC comes with a lots of benefits over any other traditional solution to the same problem. Some of the key points are as follows -

- Whether your data resides in Postgres or a JSON file or MongoDB or in all three places, abc can index the data into Elasticsearch. Besides these, it also supports CSV, MySQL, SQLServer, and Elasticsearch itself to an Elasticsearch index.
- It can keep the Elasticsearch index synced in realtime with the data source. (Note: Currently only supported for MongoDB and Postgres)
- `abc import` is a single line CLI command that allows doing all of the above. It doesn’t require any external dependencies, takes zero lines of code configuration, and runs as an isolated process with a minimal resource footprint.
- abc also supports configureable user defined transformations for advanced uses to map data types, columns or transform the data itself before it gets indexed into Elasticsearch.



<a name="getting-started"></a>
## 3. Getting Started

ABC can be downloaded as an executable as well as through a Docker image. 

#### Using Executable

Download `abc`'s executable [from releases](https://github.com/appbaseio/abc/releases/latest) for your platform and preferrably put it in a PATH directory.
The access it as -

```sh
> abc
```

You should see a list of commands that `abc` supports.
Try logging in for example.

#### Using Docker

To use the Docker image, pull it as 

```sh
docker pull appbaseio/abc
```

Then create the volume to store config files across containers.

```sh
docker volume create --name abc
```

Finally you should be able to use `abc`

```sh
docker run -i --rm -v abc:/root appbaseio/abc
```

This command may look too long to you. We can create an alias to make things better. 

```sh
# create alias
alias abc='docker run -i --rm -v abc:/root appbaseio/abc'
# run a command
abc login google
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
COMMANDS
  login     login into appbase.io
  user      get user details
  apps      display user apps
  app       display app details
  create    create app
  delete    delete app
  logout    logout session
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
For more details on adaptors, see [Import docs](docs/appbase/import.md).



<a name="development"></a>
## 5. Development

ABC can be built locally via the traditional `go build` or by building a Docker image.

<a name="local-setup"></a>
### 5.1 Local Setup

You can install ABC by building it locally and then moving the executable to anywhere you like. 

To build it, you will require **Go 1.8** or above installed on your system. 

```sh
go get github.com/appbaseio/abc # alternatively, clone the repo in the `$GOPATH/src/github.com/appbaseio/abc` dir
cd $GOPATH/src/github.com/appbaseio/abc
go build -tags 'oss' ./cmd/abc/...
./abc --help  # voila, you just built abc from source!
```

Note - You might be wondering what is the tag `oss` doing there. That's covered in the section [Build Variants](#build-variants).

<a name="docker-setup"></a>
### 5.2 Docker Setup

```sh
git clone https://github.com/appbaseio/abc
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
# setting alias for easy usage
alias abc='docker run -i --rm -v abc:/root abc'
# using alias now :)
abc user
abc apps
```

<a name="build-variants"></a>
### 5.3 Build Variants

The ABC project you see in this repository is not the complete project. Appbase.io works on a proprietary version of ABC using this project as the base.
Hence we use the tag 'oss' to specify that this is an open source build. 
If you are curious, we use the tag '!oss' to make our private builds. 


#### How to know build variant from the executable? 

If you are not sure which build of `abc` you are using, you can run `abc version` and take note of the value under the VERSION header. 

For open source build, you will see

```
VERSION
  ... (oss)
```

For the proprietary builds, you will see 

```
VERSION
  ... (!oss)
```



<a name="abc-resources"></a>
## 6. ABC Resources

Checkout the [docs folder](docs/) for details on some ABC commands and topics.

<a name="contributing-to-abc"></a>

### 6.1 Contributing to ABC

Want to help out with ABC? Great! There are instructions to get you started [here](CONTRIBUTING.md).

<a name="licensing"></a>

### 6.2 Licensing

ABC's oss variant is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for full license text. ABC's !oss (read non-oss) variant which includes the `abc import` command and bundled in the binary is free to use while in beta.

