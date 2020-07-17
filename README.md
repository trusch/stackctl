stackctl
========

Lets you run compose files as rootless podman pods!

## Scope

This project aims at proving a tool to seamlessly spin-up and manage [rootless podman](https://github.com/containers/podman) pods following a compose file.
It provides some additional extensions which try to make the developer workflow with the application stack as easy as possible.

## Usage

Lets assume the following `compose.yaml` file that lives in a folder called `test-project`:

```yaml
version: "3.8"
services:
  alpine1:
    image: alpine
    command: ["tail", "-f", "/dev/null"]
  alpine2:
    image: alpine
    command: ["tail", "-f", "/dev/null"]
```

You can now easily spin that up by using `stackctl up`:

```bash
➜  test-project git:(master) ✗ stackctl up
➜  test-project git:(master) ✗ stackctl status
       ID      |  NAME   |              IMAGE              |      STATUS
---------------+---------+---------------------------------+-------------------
  2d2d21f7a4f5 | alpine1 | docker.io/library/alpine:latest | Up 3 seconds ago
  bb71867756fb | alpine2 | docker.io/library/alpine:latest | Up 2 seconds ago
```

You also have more fine grained control over the process if you want:

```bash
➜  test-project git:(master) ✗ stackctl create
INFO[0000] creating pod
INFO[0000] creating container for alpine1
INFO[0000] creating container for alpine2
➜  test-project git:(master) ✗ stackctl status
       ID      |  NAME   |              IMAGE              | STATUS
---------------+---------+---------------------------------+----------
  612c7e3e1b18 | alpine1 | docker.io/library/alpine:latest | Created
  6361bbfad2ae | alpine2 | docker.io/library/alpine:latest | Created
➜  test-project git:(master) ✗ stackctl start alpine1
INFO[0000] starting container for alpine1
➜  test-project git:(master) ✗ stackctl status
       ID      |  NAME   |              IMAGE              |      STATUS
---------------+---------+---------------------------------+-------------------
  612c7e3e1b18 | alpine1 | docker.io/library/alpine:latest | Up 2 seconds ago
  6361bbfad2ae | alpine2 | docker.io/library/alpine:latest | Created
➜  test-project git:(master) ✗ stackctl start alpine2
INFO[0000] starting container for alpine2
➜  test-project git:(master) ✗ stackctl status
       ID      |  NAME   |              IMAGE              |      STATUS
---------------+---------+---------------------------------+--------------------
  612c7e3e1b18 | alpine1 | docker.io/library/alpine:latest | Up 14 seconds ago
  6361bbfad2ae | alpine2 | docker.io/library/alpine:latest | Up 8 seconds ago
➜  test-project git:(master) ✗ stackctl stop
INFO[0000] stopping pod
➜  test-project git:(master) ✗ stackctl remove
INFO[0000] removing pod
```

If you want to restart a service, but not rerendering the image, you can do so using `stackctl restart alpine1`

If you want to really rerender the image and start a new container from this use `stackctl recreate alpine1`. This is usefull when you locally rebuild an image and want to upgrade one of your services with the newest local image.

## Additional nice things

### `recreate` with tricks

You can override the image that is used when recreating a service without touching your compose file. Just supply it as argument to the `stackctl recreate` call:

```bash
➜  test-project git:(master) ✗ stackctl recreate alpine1 --image alpine:edge
INFO[0000] start service alpine1
INFO[0001] remove service alpine1
INFO[0001] create service alpine1
Trying to pull docker.io/library/alpine:edge...
Getting image source signatures
Copying blob 5d2415897100 done
Copying config 3c791e92a8 done
Writing manifest to image destination
Storing signatures
93a9ee4f683ed714ac5ddae50652df186feed762acc7a7f7f5af07ab730690be
INFO[0006] start service alpine1
➜  test-project git:(master) ✗ stackctl status
       ID      |  NAME   |              IMAGE              |      STATUS
---------------+---------+---------------------------------+--------------------
  93a9ee4f683e | alpine1 | docker.io/library/alpine:edge   | Up 2 seconds ago
  03e27a4a6aa5 | alpine2 | docker.io/library/alpine:latest | Up 10 minutes ago

```

If you have images build from your CI for pull requests of specific components you can also use the special config directive `x-pr-template` to reference the PRs while recreating.

```bash
➜  test-project git:(master) ✗ cat compose.yaml
version: "3.8"
x-pr-template: "your-registry.io/pr-templates/{{ .Service }}:{{ .PR }}"
services:
  alpine1:
    image: alpine
    command: ["tail", "-f", "/dev/null"]
  alpine2:
    image: alpine
    command: ["tail", "-f", "/dev/null"]
➜  test-project git:(master) ✗ stackctl recreate alpine1 --with-pr 123
INFO[0000] stop service alpine1
INFO[0001] remove service alpine1
INFO[0001] create service alpine1
Trying to pull your-registry.io/pr-templates/alpine1:123...
[...]
```

You can also supply a service local `x-pr-template` directive that takes precedence over the global template in case the PR image name is not simply constructable by using the service name and the PR number (like if you have two services which have different names, but the same image repository).

### Centralize your port forwarding

Since all services get created in one network namespace, and therefore all port declarations can't be conflicting, we can also declare them on the toplevel to have them easily discoverable:

```yaml
version: "3.8"
x-ports:
  "127.0.0.1:3001": "3001/tcp"
services:
  alpine1:
    image: alpine
    command: ["tail", "-f", "/dev/null"]
  alpine2:
    image: alpine
    command: ["tail", "-f", "/dev/null"]

```

## Completeness And Compability Disclaimer

This doesn't implement all of the [compose-spec](https://github.com/compose/compose-spec). In fact it just implements the most important and most used parts of the spec needed to define your services. Since the goal of this project is to provide a developer tool that uses rootless podman pods, its also simply not possible to implement certain things like networking, resource contraints, privileges etc. Once rootless containers are matured, so that these are all doable things, I would be more than happy to extend in this regard.

There is also no (or just very limited) command line compability to `docker-compose`.

