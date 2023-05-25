# WebRepo

Repo manager with web technology. It allows to create and manage repositories through REST API. Additionally includes web-panel for easy management.

Developed to be used under Alpine-linux, but it should work on other Linux distributions or Windows.

## Architecture

* There are two parts: client and server.
* Client is written using [Materialize](https://github.com/materializecss/materialize), [Material icons](https://fonts.google.com/icons?selected=Material+Icons) and Cash (alternative to jquery).
* Server is written using Golang and [Echo](https://echo.labstack.com).

## API

* `/api/plugins?action=list` - Get json file with list of plugins. NOTE: Plugin list is updated only on startup.
* `/api/plugins?action=get_status&id=your_id` - Get status of job with `your_id` id.
* `/api/plugins?action=jobs_clear` - Clear completed jobs. Sends `bad requested` if any job is not completed.
* `/api/plugins/your_plugin?action=add&param=your_param` - Start `your_plugin` with `your_param` parameter.
* `/api/internal?action=set_auth&user=your_user&pass_hash=your_pass_hash` - Set defined user and pass_hash.

## Notes

* Please, do _not_ use it in production. It is not enough tested yet.
