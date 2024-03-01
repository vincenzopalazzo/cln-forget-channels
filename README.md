<div align="center">
  <h1>cln4go</h1>

  <img src="https://preview.redd.it/tcmyd3n69ng41.jpg?width=1999&format=pjpg&auto=webp&s=b79cf22d3e2adcaf52a2d22bcb0568e42eff8bc2" />

  <p>
    <strong> Go library for cln with flexible interface </strong>
  </p>

  <h4>
    <a href="https://github.com/vincenzopalazzo/cln4go">cln4go Homepage</a>
  </h4>
</div>

## Table of Content

- Introduction
- How to use
- How to contribute
- Build with
- License

## Introduction

A minimal template for projects based on [cln4go](https://github.com/vincenzopalazzo/cln4go). This should get you started if you want to 
develop a core lightning plugin using Go.

## How to use

You can use the make file in the following way make build OS=<> ARCH=<> for example `make build OS=darwin ARCH=amd64` for Mac M1

To produce a binary for the OS and ARCH.

However make build should be used only in release mode, in development mode you must use make that produce a debug binary, and also make a `go fmt` for you, In addition, it makes a static analysis of the code with the following tool https://github.com/golangci/golangci-lint

So you should install it.

This binary can be used with core lightning using the

```
--plugin = path_to_this_executable_binary
```

## Built with

- [cln4go](https://github.com/vincenzopalazzo/cln4go)

## License

<div align="center">
  <img src="https://opensource.org/files/osi_keyhole_300X300_90ppi_0.png" width="150" height="150"/>
</div>

Template to write a plugin for core lightning in dart lang.

Copyright (C) 2022 Vincenzo Palazzo vincenzopalazzodev@gmail.com
``

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.