# login

login command allows you to login into `abc`. 

Here are a few examples -- 

```sh
# view login help
> abc login --help
```

```sh
# returns current user information
> abc login
Logged in as some@email.com
```

```sh
# login with google
> abc login google
```

```sh
# login with github
> abc login github
```

Once you start login into google or github, you will be asked to authenticate `abc` with your google/github account if it isn't already done.
Then you will be redirected to a page with your Auth token. 
You will have to paste the token in the command-line window to be authenticated.
