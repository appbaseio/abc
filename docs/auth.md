# Authentication in ABC

A user can be authenticated in abc using the `login` command or by setting an environment variable.


### login command

Please see [login docs](appbase/login.md) for more details.

```sh
# login using google
abc login google
# login using github
abc login github
```


### Environment variable

User can set an environment variable `ABC_TOKEN` with their token to have them authenticated in abc.
To get your token, please visit the following url - 

[Get token - Google Auth](https://accapi.appbase.io/login/google?next=https://accapi.appbase.io/user/token)

[Get token - GitHub Auth](https://accapi.appbase.io/login/github?next=https://accapi.appbase.io/user/token)

Once `ABC_TOKEN` is set, `abc` can be used without having the need to run `abc login` command.

```sh
export ABC_TOKEN=myFullTokenTextHere
# you should be authenticated now
abc user
```
