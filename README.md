<p align="center">
  <img src="https://i.issei.space/23pcsoHy.png" />
</p>

# Nordlayer helper
This is my pet project to learn Golang and try to understand the GRPC protocol (with reverse engineering effort). The idea behind this project was to create a simple taskbar menu with the most crucial information from `nordlayer status`. Unfortunately, on linux Nord doesn't provide a GUI version (like on Mac and Windows).

The GRPC API is also not documented publicly, so all the work was done using package capture and trying random things.

## Requirements
On Fedora with GNOME desktop, you need to install [AppIndicator and KStatusNotifierItem Support](https://extensions.gnome.org/extension/615/appindicator-support/) extenstion

## What is working?
* Displaying current status of connection
* Listing gateways

I am also planning to include an option to change gateways and a notification system about being offline, External IP address changes.

## Building
To build this application, you need to have `make gcc libgtk-3-dev libayatana-appindicator3-dev protoc`
### Build instructions
```
make build
```
This command will build `pb.go` files from `.proto` and build the application.

## Disclaimer
This repository is not associated in any way with Nord Security/Nordlayer