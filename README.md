# warewulf v4

![Warewulf](warewulf-logo.png)

#### Quick Links:

* [Website](https://warewulf.org)
* [Documentation](https://warewulf.org/docs)
* [Download / Releases](http://github.com/warewulf/warewulf/releases)
* [Slack](http://warewulf.slack.com/)
* [Support](https://ciq.com/products/warewulf)

## About Warewulf

### History

For over two decades, Warewulf has powered HPC systems around the
world. From simple “under the desk” clusters to large institutional
systems at HPC centers as well as enterprises who rely on performance
critical computing.

Through the evolution of Warewulf, we have seen various iterations
provisioning models starting from CDROM / ISO images to Etherboot
(predecessor to PXE), then PXE, and more recently iPXE, but even
during these different bootloaders, Warewulf in it’s heart, has always
been first and foremost a stateless provisioning system (e.g. the
operating system node image is not written to any persistent storage
and rather it boots from the network directly into a runtime system).

Warewulf v3 has been in production for over 6 years now as it has
stabilized into a very solid and full featured solution. But over the
last few years, there have been many innovations in Enterprise
technologies which can (and should) be leveraged as part of
Warewulf. Additionally, some of the lessons learned from Warewulf v3
architecture should be rolled into an updated architecture for
provisioning management.

### Warewulf v4

Leveraging this legacy of provisioning and cluster management brings
us to where we are today. The next generation of Warewulf. Warewulf v4
is a complete rewrite in GoLang, taking in the legacy of what we've
come to expect with Warewulf, bringing it into the present, and
looking out into the future.

Warewulf v4 combines ultra scalability, flexibility, and simplicity
with being light weight, non-intrusive, and a great tool for
scientists and seasoned system administrators alike. Warewulf empowers
you to scalably and easily manage thousands of compute resources.

### Architecture

One of the design tenants of Warewulf is how to scalably administrate
many thousands of compute nodes. Generally speaking, operating system
state introduce a surface for potential discrepancies and version
creep between nodes and thus Warewulf has always gone with the "single
system image" approach to clustered operating system management.  This
means that you can have a single node "image".

At its core, Warewulf v4 focuses on what has made Warewulf so widely
loved: simplicity, ultra scalable, lightweight, and an easy to manage
solution built for both scientists and seasons system administrators
to be able to design a highly functional yet easy to maintain cluster
no matter how big or small or customized it needs to be.

## iPXE

Warewulf uses iPXE for network boot. Typically iPXE is provided by the
operating system; but the iPXE binaries can be built with
`scripts/build-ipxe.sh`. This script accepts command-line arguments that are
passed to the underlying `make` process. e.g.,

```bash

echo "#!ipxe
echo Tagging with vlan 1000
vcreate --tag 1000 net0 autoboot || shell" >vlan-1000.ipxe

sh scripts/build-ipxe.sh EMBED=$(readlink -f vlan-1000.ipxe)
```

By default, `build-ipxe.sh` will attempt to write iPXE builds to
`/usr/local/share/ipxe/`. This path can be specified using the `DESTDIR`
environment variable. Other supported environment variables include
`IPXE_BRANCH` and `TARGETS`.

```bash

IPXE_BRANCH=master TARGETS=bin-arm64-efi/snponly.efi DESTDIR=. sh scripts/build-ipxe.sh
```

Update `warewulf.conf` to use locally-built iPXE.

```
tftp:
  enabled: true
  systemd name: tftp
  ipxe:
    00:09: /usr/local/share/ipxe/bin-x86_64-efi-snponly.efi
    00:00: /usr/local/share/ipxe/bin-x86_64-pcbios-undionly.kpxe
    00:0B: /usr/local/share/ipxe/bin-arm64-efi-snponly.efi
    00:07: /usr/local/share/ipxe/bin-x86_64-efi-snponly.efi
```
