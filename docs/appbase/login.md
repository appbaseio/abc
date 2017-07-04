# login

login command allows you to login into `abc`. 

Here are a few examples -- 

```sh
abc login --help
# view login help

abc login
# returns current user information

abc login google
# login with google

abc login github
# login with github
```

Once you start login into google or github, you will be asked to authenticate `abc` with your google/github account if it isn't already done.
Then you will be redirected to a page with your Auth token. 
You will have to paste the token in the command-line window to be authenticated.
