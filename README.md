# GoGL

GoGL is a 2D graphics library built on top of [Simple Directmedia Layer](https://www.libsdl.org/). It is CPU-based, meaning it does not perform well with intensive graphics.

This project was built as a learning experience. Therefore, it is not recommended for any serious applications.

## Capabilities

This library does:
- Draw 2D shapes
- Render text
- Handle keyboard and mouse inputs
- Provide fine-grained control of primitives and their interactions

This library does not:
- Support animations*
- Provide "game engine"-like functionality

*Animations and other complex graphics can be built on top of gogl. See [go-2048-battle](http://github.com/z-riley/go-2048-battle) as an example.

## Dependencies

### SDL2 - DirectMedia Layer

#### apt (Ubuntu):
```sh
sudo apt-get update && sudo apt-get install -y libsdl2-dev libsdl2-image-dev
```

#### brew (MacOS):
```sh
brew install sdl2_image
```

For other platforms, see [wiki.libsdl.org/SDL2/Installation](https://wiki.libsdl.org/SDL2/Installation/)
