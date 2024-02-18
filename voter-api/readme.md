## Voter API

`make build` to build

`make run` to start the server (listens on 0.0.0.0:1080)

`make test` to run tests, note that server needs to be running for e2e tests to run

```
➜  voter-api git:(main) make
Usage make <TARGET>

  Targets:
           build                        Build the todo executable
           run                          Run the todo program from code
           run-bin                      Run the todo executable
           load-db                      Add sample data via curl
           get-by-id                    Get a todo by id pass id=<id> on command line
           get-all                      Get all todos
           update-2                     Update record 2, pass a new title in using title=<title> on command line
           delete-all                   Delete all todos
           delete-by-id                 Delete a todo by id pass id=<id> on command line
           get-v2                       Get all todos by done status pass done=<true|false> on command line
           get-v2-all                   Get all todos using version 2
```

