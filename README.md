# Desert Island NFT Generator tool

## Desert Island NFT Collection
[Make sure to stay up to date with news about our upcoming NFT collection!](https://linktr.ee/desertislandone)

## About
This repository provides various utilities to help you build your NFT collection :joy: :

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
* Add `layers_priority.json` file to the root of the project
  * Add JSON configuration for layers priority:
  ```json
      {
        "layers_priority": [
          {
            "name": "Background",
            "priority": 0
          },
          {
            "name": "Torso",
            "priority": 1
          }
        ]
      }
    ```
    * Make sure that priority is set in correct order, it will be used to put layers in correct order
* Make sure to setup environment variables using `.env` file
* Run command `make run` to start collection generation

## TODO
- [x] Add image generation
- [ ] Add DNA generation algorithm (to not generate duplicates)
- [ ] Add ERC-721 metadata generation
- [ ] Add IPFS support

## Third-party libraries used
* [Logs](https://github.com/rs/zerolog)
* [Environment configuration](https://github.com/caarlos0/env)
* [Dotenv](https://github.com/joho/godotenv)
* [Image resizing](https://github.com/disintegration/imaging)
* [Errors](https://github.com/juju/errors)
* [UUID](https://github.com/google/uuid)

Many thanks to the authors of these libraries!

## Special thanks!

Big thanks to the guy who wrote the post on [habr.ru](https://habr.com/ru/post/595723/) and his [github](https://github.com/golang-enthusiast/nft)!
This library is based on the original authors' code, but was 99% re-written in my own manner :joy: