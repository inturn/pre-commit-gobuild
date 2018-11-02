# pre-commit-gobuild 

The project contains pre-commit https://pre-commit.com/ hooks for building and running unit tests towards a go project. 
The hooks are written on golang as the example of how it may be done. Of course, the same 
functionality can be made with a bush script but with the golang script we can make the process
of building more flexible and adjusted for some personal needs.

This is the example of the .pre-commit-config.yaml settings file in order to use the hooks in the project:

```
fail_fast: false
repos:
-   repo: git://github.com/guntenbein/pre-commit-gobuild
    rev: HEAD
    hooks:
    -   id: go-build
        stages: [push]
    -   id: go-test
        stages: [push]
```
        
   
You need to execute the following in order to clean the pre-commit cache 
and install the hooks for your repo (should be run from the root of the repo):

```
pre-commit clean
pre-commit install -f --hook-type pre-commit
pre-commit install -f --hook-type pre-push
```