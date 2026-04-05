# Proton Pass Exports CLI - `protonpe`

`protonpe` is a CLI tool that lets you interact with Proton Pass PGP-encrypted exports to retrieve secrets completely offline.

Quickly access passwords, TOTPs, notes, and more using the local backup of your vaults without the need to decrypt everything into disk and go through a JSON with thousands of lines.  

## Documentation

Please check the project's [wiki](https://github.com/diogoftm/protonpe/wiki).

## Quick start

Start by installing the tool from the repository as follows:
```bash
go install github.com/diogoftm/protonpe@latest
```

The `protonpe` command must then be available. If not, on Linux, make sure `$GOPATH/bin` is part of the user's `$PATH`.

For example, to retrieve general information about the available vaults on a file named `data.pgp`, simply run:
```bash
protonpe info -f data.pgp
```

To simplify the usage of the tool, the `PROTONPE_FILE` environment variable can be set with the absolute path to the file exported from Proton Pass removing the need to indicate it in every command.

---

> This software is distributed under the MIT [license](license.txt). It is an independent project and is not affiliated with, endorsed by, or supported by Proton AG. Nonetheless, please always prefer European technology made proudly by European companies!
