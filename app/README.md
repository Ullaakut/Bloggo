# Bloggo app frontend

This folder contains the sources of the frontend, as well as the Dockerfile and nginx configuration to deploy it.

## Webapp

The style is built using `sass`, and the HTML is just static HTML. We might use React in the future when we need the interface to become more complex.

To compile the style, you can either use the docker image to avoid installing dependencies on your computer, or you can use the sass compiler.

### Sass installation

#### Standalone

You can install Sass on Windows, Mac, or Linux by downloading the package for your operating system from [GitHub](https://github.com/sass/dart-sass/releases/tag/1.13.4) and [adding it to your PATH](https://katiek2.github.io/path-doc/). That's all—there are no external dependencies and nothing else you need to install.

#### npm install

If you use Node.js, you can also install Sass using npm by running

`npm install -g sass`

However, please note that this will install the pure JavaScript implementation of Sass, which runs somewhat slower than the other options listed here. But it has the same interface, so it'll be easy to swap in another implementation later if you need a bit more speed!

#### Install on Windows

If you use the Chocolatey package manager for Windows, you can install Dart Sass by running

`choco install sass`

#### Install on Mac OS X

If you use the Homebrew package manager for Mac OS X, you can install Dart Sass by running

`brew install sass/sass/sass`

## Deployment

The Dockerfile simply compiles the `sass` and serves the folder as static files.

## Nginx configuration

Deploying the app with nginx allows to enable gzip compression, URL rewrite, etc.
