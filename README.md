# nginxfmt

Takes [Nginx](https://nginx.org/) configuration files and run them through
[go-crossplane](https://github.com/aluttik/go-crossplane).

This program is heavily inspired by `gofmt` and re-uses a lot of code and
style from that program.

## Usage

```
nginxfmt [flags] [path ...]
```

Flags are:

- `-l`
  Do not print reformatted configuration to standard output.
  If a file's formatting is different from `nginxfmt`'s, print its name
  to standard output.
- `-w`
  Do not print reformatted configuration to standard output.
  If a file's formatting is different from `nginxfmt`'s, overwrite it with
  `nginxfmt`'s version.

The `path`s can be a list of file names and/or directories. If a directory is
given, it is assumed all files in the directory are Nginx configuration files
for processing.
