# Rubbi-sh

This is a Go and Shell tool that provides a sandbox directory for rubbish.

If you usually run tons of commands, generating garbage everywhere, you need a place to work without generating dirty left overs.

## Getting Started

  * [Install rubbi-sh](install.sh)
    * - Download the [latest release](https://github.com/jotadrilo/rubbi-sh/releases)
    * - Uncompress the archive
    * - Run `install.sh`
  * [Get started with rubbi-sh](#examples)

## Examples

### rubcd

This is the core functinoality. You will get a sandbox for your rubbish and will get `cd`ed to it. This sandbox will be re-created every day and will keep old ones until you reboot your device (they are stored at `/tmp` by default).

```
> jotadrilo @ ~ $ rubcd
> jotadrilo @ /tmp/rubbish/20190614 $
```

> **NOTE**: During the first run, it will create a new configuration file at `$HOME/.rubbish`. Example:
> ```
> $ cat ~/.rubbish/config.json
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

### rubbish

This is an alias for `rubcd`

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
$ jotadrilo @ ~ $ rubadd foo
/tmp/rubbish/foo
> jotadrilo @ ~ $ rubshow
[0] 20190613	/tmp/rubbish/20190613
[1] 20190614	/tmp/rubbish/20190614
[2] foo     	/tmp/rubbish/foo
```

### rubdel

This helper will delete an existing rubbish folder by number.

```
> jotadrilo @ ~ $ rubshow
[0] 20190613	/tmp/rubbish/20190613
[1] 20190614	/tmp/rubbish/20190614
[2] foo     	/tmp/rubbish/foo
$ jotadrilo @ ~ $ rubdel 2
> jotadrilo @ ~ $ rubshow
[0] 20190613	/tmp/rubbish/20190613
[1] 20190614	/tmp/rubbish/20190614
```

### rubsel

This helper will prompt the list of folders and will ask for a folder number to `cd`.

```
> jotadrilo @ ~ $ rubsel
[0] 20190613	/tmp/rubbish/20190613
[1] 20190614	/tmp/rubbish/20190614

Folder to use: 1
> jotadrilo @ /tmp/rubbish/20190614 $
```