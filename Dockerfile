# Supposedly, there is a way to actually make the GUI show up in a docker container.
# https://medium.com/geekculture/run-a-gui-software-inside-a-docker-container-dce61771f9
# ^- article that talks about it (method only works on linux, sadly)
# FROM ubuntu:22.04
FROM node:20-bullseye
RUN apt-get update && apt-get install -y ca-certificates

# Copying over data
WORKDIR /synk
COPY ./ ./
# Copying golang install
COPY --from=golang:1.25.6 /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
# Installing wails + deps
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest
ENV PATH="/go/bin:${PATH}"
RUN yes | apt install build-essential libgtk-3-dev libwebkit2gtk-4.0-dev pkg-config
RUN npm install -g npm@11.10.0
CMD [ "wails", "dev" ]