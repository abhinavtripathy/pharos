[![Actions Status](https://github.com/abhinavtripathy/pharos/workflows/Node%20CI/badge.svg)](https://github.com/abhinavtripathy/pharos/actions)
[![License](http://img.shields.io/badge/License-MIT-brightgreen.svg)](./LICENSE)
![GitHub stars](https://img.shields.io/github/stars/abhinavtripathy/pharos.svg)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/abhinavtripathy/pharos.svg?color=red)

# pharos

> Pharos *Noun* A light house to direct beginners


### Motivation

Pharos is built for high school and middle school students who are transitioning from a programming tool like scratch to an actual programming language. This language is meant for easing the transition. 

### Tech Stack

Since this is a Domain Specific Language, we wanted base language to be fast, popular so people can contribute to it and give us robust performance. We chose to build our language to built on top of Go. 


## Installation 

For front-end:

```
npm install 

npm start
```

or 

```
yarn install

yarn start
```

For back-end:

Install rust from https://www.rust-lang.org/tools/install

We use a library called Rocket for building our server. Rocket needs the beta verstion of Rust (called Rust Nightly) to run. This version of rust is unstable but Rocket is one of the best libraries we have and for the our initial puroposes this will suffice. A stable version of Rocket is expected to be released soon, we will probably update to that version or rebuild the server with a more stable version for the final version.  

After installing Rust you need to go the directory backed. Check the version of rust with the following code:

```
rustc --version
```

if it is not nightly enter the following code:

```
rustup override set nightly
```

You can also revert it back to stable version with the following code:

```
rustup override set stable
```

To setup the server run:

```
cargo build
cargo run
```



