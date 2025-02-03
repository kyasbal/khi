# Development guide

## Run your first build

Follow [the "Run from source code" section](../../README.md) on README.

## Setup environment for development

### Setup Git hook

Run the following shell command to setup Git hook. It runs format or lint codes before commiting changes.

```shell
$ make setup-hooks
```

### Setup VSCode config

Save the following code as `.vscode/launch.json`.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Start KHI Backend",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "./cmd/kubernetes-history-inspector/",
            "cwd": "${workspaceFolder}",
            "args": [
                "--host",
                "127.0.0.1",
                "--port",
                "8080",
                "--frontend-asset-folder",
                "./dist",
            ],
            "dlvLoadConfig": {
                "followPointers": true,
                "maxVariableRecurse": 1,
                "maxStringLen": 100000,
                "maxArrayValues": 64,
                "maxStructFields": -1
            },
        }
    ], 
}
```

You can run KHI with VSCode and features like break points are available with it.

### Run frontend server for development

To develop frontend code, we usually start Angular dev server on port 4200 with the following code.

```shell
$ make watch-web
```

This will build the frontend code with [the configuration to access APIs on `localhost:8080`](../../web/src/environments/environment.dev.ts).
You can use KHI with accessing `localhost:4200` instead of `localhost:8080`. Angular dev server automatically build and serve the new build when you change the frontend code.

### Run test

Run the following code to verify frontend and backend codes.

```shell
$ make test
```

## Auto generated codes

### Generated codes from backend codes

Several frontend codes are automativally generated from backend codes.

* `/web/src/app/generated.sass`
* `/web/src/app/generated.ts`

These files are generated with [`scripts/frontend-codegen/main.go` Golang codes](../../scripts/frontend-codegen/main.go). It reads several Golang constant arrays and generate frontend codes with templates.
