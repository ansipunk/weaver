# Weaver

A Minecraft server manager.

## Usage

```sh
go install 'git.sr.ht/~ansipunk/weaver/cmd/weaver'

weaver add starlight collective
weaver install
weaver remove collective
```

## Roadmap

- [ ] Server configuration: `pkg/cfg`
  - [x] Loader
  - [x] Game version
  - [ ] Mods
    - [x] Basic list
    - [ ] Version constraints
    - [x] Dependencies
- [ ] Fabric API: `pkg/fabric`
  - [ ] Get latest core
  - [ ] Get specific core
- [ ] Lockfile: `pkg/lockfile`
  - [ ] Installed packages
  - [ ] Package hashes
  - [ ] `weaver.toml` hash
- [x] Modrinth API: `pkg/modrinth`
  - [x] Get latest version
  - [x] Get specific version
- [x] FS: `pkg/fs`
  - [x] Compare file hashes
  - [x] Create `mods/` directory
  - [x] Clear mods directory
- [x] Command-line interface: `cmd/weaver`
  - [x] Threading
  - [x] Logging
  - [ ] Initialize project
  - [x] Install mods
  - [x] Add mods
  - [x] Remove mods

## File example

```toml
# weaver.toml
loader = "fabric"
game_version = "1.20"
mods = [
    "fabric-api",
    "lithium",
]
```
