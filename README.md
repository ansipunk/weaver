# Weaver

A Minecraft server manager.

## Usage

```sh
go install 'git.sr.ht/~ansipunk/weaver/cmd/weaver'
weaver
```

## Roadmap

- [ ] Server configuration: `pkg/cfg`
  - [x] Loader
  - [ ] Version constraints
  - [x] Game version
  - [ ] Mods
    - [x] Basic list
    - [ ] Version constraints
    - [ ] Dependencies
- [ ] Modrinth API: `pkg/modrinth`
  - [x] Get latest version
  - [ ] Get specific version
- [ ] Fabric API: `pkg/fabric`
  - [ ] Get latest core
  - [ ] Get specific core
- [ ] Lockfile: `pkg/lockfile`
  - [ ] Installed packages
  - [ ] Package hashes
  - [ ] `weaver.toml` hash
- [ ] FS: `pkg/fs`
  - [x] Compare file hashes
  - [x] Create `mods/` directory
  - [ ] Clear mods directory
- [ ] Command-line interface: `cmd/weaver`
  - [x] `install`
  - [ ] `add`
  - [ ] `update`

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
