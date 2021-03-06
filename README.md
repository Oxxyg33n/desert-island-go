# Desert Island NFT Generator tool

## Desert Island NFT Collection
[Make sure to stay up to date with news about our upcoming NFT collection!](https://linktr.ee/desertislandone)

## About
This repository provides various utilities to help you build your NFT collection :rocket: :

    Generate images from source layers in specified order
    Generate ERC-721 traits
    Upload images & metadata to IPFS

## Steps to run
* Add your layer images to `input/` directory
  * Example file structure:
  ```
    input/Background/
                Common/0.png
                Common/1.png
                Rare/1.png
                Epic/2.png
    input/Torso/
                Common/0.png
                Rare/1.png
                Epic/2.png
    ```
* Add `collection_configuration.json` file to the root of the project
  * JSON configuration must be in the following format:
  ```json
      {
        "rareness": [
          {
            "name": "Common",
            "chance": 85
          },
          {
            "name": "Rare",
            "chance" :10
          },
          {
            "name": "Epic",
            "chance": 5
          }
        ],
        "layers": {
          "skip_multiple": false,
          "groups": [
            {
              "name":"Background",
              "priority": 0,
              "can_skip": false
            },
            {
              "name":"Base Torso",
              "priority": 1,
              "can_skip": true,
              "skip_chance": 20.5
            }
          ]
        }
      }
    ```
    * Make sure that priority is set in correct order, it will be used to put layers in correct order
* Make sure to setup environment variables using `.env.sample` file and rename it to `.env` afterwards
* Run command `make run` to start collection generation

## TODO
- [x] Add image generation
- [x] Add DNA generation algorithm (to not generate duplicates)
  - [ ] Improve DNA algorithm to be more smart
- [x] Add additional configuration to be able to skip adding some layers to the end image
- [x] Add ERC-721 metadata generation
- [x] Add IPFS support
- [ ] Refactor traits chance generation algorithm
- [ ] Add CLI command to generate collection_configuration.json
- [ ] Add comments to the code
- [ ] Concurrency support for faster generation

## Third-party libraries used
* [Golangci-lint](https://github.com/golangci/golangci-lint)
  * [Installation](https://golangci-lint.run/usage/install/#local-installation)
* [Logs](https://github.com/rs/zerolog)
* [Environment configuration](https://github.com/caarlos0/env)
* [Dotenv](https://github.com/joho/godotenv)
* [Image resizing](https://github.com/disintegration/imaging)
* [Errors](https://github.com/juju/errors)
* [UUID](https://github.com/google/uuid)

Many thanks to the authors of these libraries!

## Special thanks

Big thanks to the guy who wrote the post on [habr.ru](https://habr.com/ru/post/595723/) and his [github](https://github.com/golang-enthusiast/nft)!
This library is based on the original authors' code, but was 99% re-written in my own manner