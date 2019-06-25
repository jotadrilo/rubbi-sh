# Rubbi-sh

This is a Go and Shell tool that provides a sandbox directory for rubbish.

If you usually run tons of commands that generate garbage, wherever you are, you need an easy way to access a temporary folder where to run these tasks.

```
> jotadrilo @ ~ $ rbsh
> jotadrilo @ /tmp/rubbish/20190614 $
```

## TL;DR

### Linux

```
mkdir -p rubbi-sh && \
curl -sSL https://github.com/jotadrilo/rubbi-sh/releases/download/0.1.0/rubbi-sh_0.1.0_linux_x86_64.tar.gz | tar xzf - -C rubbi-sh && \
sudo bash -c 'cd rubbi-sh; ./install.sh'
```

### MacOS

Homebrew users:

```
brew tap jotadrilo/tap
brew install jotadrilo/tap/rubbi-sh
```

Alternative:

```
mkdir -p rubbi-sh && \
curl -sSL https://github.com/jotadrilo/rubbi-sh/releases/download/0.1.0/rubbi-sh_0.1.0_darwin_x86_64.tar.gz | tar xzf - -C rubbi-sh && \
sudo bash -c 'cd rubbi-sh; ./install.sh'
```

## Examples

### rbsh

This is the core functionality. You will get a sandbox for your rubbish and will get `cd`ed to it. This sandbox will be re-created every day and will keep old ones until you reboot your device (they are stored at `/tmp` by default).

```
> jotadrilo @ ~ $ rbsh
> jotadrilo @ /tmp/rubbish/20190614 $
```

> **NOTE**: Alias: `rubcd`, `rubbish`, `r`

> **NOTE**: During the first run, it will create a new configuration file at `$HOME/.rubbish`. Example:
> ```
> > jotadrilo @ ~ $ cat ~/.rubbish/config.json
> {
>   "folders": [
>     {
>       "name": "20190614",
>       "path": "/tmp/rubbish/20190614"
>     }
>   ],
>   "latest": {
>     "name": "20190614",
>     "path": "/tmp/rubbish/20190614"
>   },
>   "root": "/tmp/rubbish"
> }
> ```

### rubshow

This helper will show the list of rubbish folders.

```
> jotadrilo @ ~ $ rubshow
[0] 20190613	/tmp/rubbish/20190613
[1] 20190614	/tmp/rubbish/20190614
```

### rubadd

This helper will add a new custom rubbish folder. It will not change the working directory.

```
> jotadrilo @ ~ $ rubadd foo
/tmp/rubbish/foo
> jotadrilo @ ~ $ rubshow
[0] 20190613	/tmp/rubbish/20190613
[1] 20190614	/tmp/rubbish/20190614
[2] foo     	/tmp/rubbish/foo
```

### rubdel

This helper will prompt the list of folders and will delete to the choosen one.

```
> jotadrilo @ ~ $ rubdel
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Folder:
    /tmp/rubbish/20190625
  ▸ /tmp/rubbish/20190626
✔ /tmp/rubbish/20190626
> jotadrilo @ ~ $
```

> **NOTE:** If the choosen folder was the latest one, the last folder in the list will become the new latest folder.

### rubsel

This helper will prompt the list of folders and will `cd` to the choosen one.

```
> jotadrilo @ ~ $ rubsel
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Folder:
  ▸ /tmp/rubbish/20190625
    /tmp/rubbish/20190626
✔ /tmp/rubbish/20190625
> jotadrilo @ /tmp/rubbish/20190625 $
```

## Development

This projects uses `go mod` and `bazel`.

```
go mod vendor
bazel run //:gazelle
bazel build //:rubbi-sh-osx
```
