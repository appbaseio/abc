# abc-import

Private submodule for abc adaptors.

### Steps

#### Prerequisites

First clone the open source abc repo.

```sh
go get github.com/appbaseio/abc
cd $GOPATH/src/github.com/appbaseio/abc
```

Then setup private module folder (poor man's submodule).

```sh
mkdir private
git clone git@github.com:appbaseio-confidential/abc-import.git private
```

Mac specific instructions to get dependencies related to Elastic.v7. You can find instructions for your platform over [here](https://github.com/appbaseio/abc/blob/dev/.travis.yml).

```sh
go get github.com/olivere/elastic/v7
```

Then build the project as follows.

```sh
go build -tags '!oss' ./cmd/abc/...   # private
go build -tags 'oss' ./cmd/abc/...    # oss
```

### Release

1. Update the `abc` version in `cmd/abc/appbase_version.go`
2. Update the `abc` version in `build.sh` file.
3. Ensure a docker build by running the command: `docker build -t appbaseio/abc:${version} -f Dockerfile .`.
4. Run the `build.sh` file to produce compressed platform binaries. They are generated under `build` directory.
5. Create a release on GitHub (provide comprehensive release notes and upload the compressed platform binaries).
6. Push the tagged docker image to Docker Hub: `docker push appbaseio/abc:${version}`.
