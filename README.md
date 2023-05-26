# WebRepo

Repo manager with web technology. It allows to create and manage repositories through REST API. Additionally includes web-panel for easy management.

Developed to be used under Alpine-linux, but it should work on other Linux distributions or Windows.

## Architecture

* There are two parts: client and server.
* Client is written using [Materialize](https://github.com/materializecss/materialize), [Material icons](https://fonts.google.com/icons?selected=Material+Icons) and Cash.
* Server is written using Golang and [Echo](https://echo.labstack.com).

## API

* `/api/settings/get_all`
* `/api/settings/set_auth`
* `/api/settings/set_scoop`
* `/api/scoop/repo/list`
* `/api/scoop/repo/init`
* `/api/scoop/repo/add`
* `/api/scoop/repo/del`
* `/api/scoop/package/list`
* `/api/scoop/package/add`
* `/api/scoop/package/del`

## Notes

* Please, do _not_ use it in production. It is not enough tested yet.
