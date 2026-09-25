# Changelog

## [3.0.0](https://github.com/VitruvianSoftware/aspect-workflows-template/compare/v2.3.0...v3.0.0) (2026-09-25)


### ⚠ BREAKING CHANGES

* template now renders via 'aspect render-preset', not 'scaffold new'. Our delivery/CI is rewired in a later phase.

### Features

* **app:** application composition core (initializer P0) ([#44](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/44)) ([b6b6ba5](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/b6b6ba5b43f6da9f44717b95f8dca418c5adde39))
* **app:** embed the initializer in every starter (P0.5, ADR-026) ([#46](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/46)) ([819d9d4](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/819d9d49b7a98952699c8de690f0916e042d0469))
* migrate template engine to upstream native AXL renderer ([#38](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/38)) ([68cd928](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/68cd928ea95dd50ba716977b017eb46e937ab8db))


### Bug Fixes

* **app:** honour the host's gazelle build-file name and import prefix (P1.1) ([#47](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/47)) ([33b03ab](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/33b03ab25deac0110963ac6383d96ec5e847efe0))
* **catalog:** resolve every dangling owner and system reference this repo emits ([#42](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/42)) ([fcb3f87](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/fcb3f87a1a7da8050356facb767b5b898625ac9e))
* **rust:** mirror rules_rs's zlib .deb via Launchpad (Ubuntu archive 404) ([#49](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/49)) ([97a20f7](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/97a20f738d821ab14924bca1a155cb4a1bb1bf35))
* **setup_repos:** read presets from deliver.yaml, run on macOS, stop stacking keys ([#43](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/43)) ([540c614](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/540c6144584eec837c9b4ad1f0177ba9d8136019))
* **template:** emit distinct Template names for ruby and scala ([#41](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/41)) ([605bfe4](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/605bfe43d82ab05520a5d55c010764de7a2c2290))
* **test:** quote $(pwd) to prevent word splitting (SC2046) ([#50](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/50)) ([50f18ab](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/50f18ab7a84280d5ca5c27a0cbb8bd9f490b6b49))

## [2.3.0](https://github.com/VitruvianSoftware/aspect-workflows-template/compare/v2.2.1...v2.3.0) (2026-06-14)


### Features

* **copybara:** add one-way PR-import sync shape (offer both modes) ([#35](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/35)) ([6f7c624](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/6f7c624890391b71d97850db722677fa20ca262e))


### Bug Fixes

* **copybara:** use 0o-prefixed octal literals in the Go tool tests ([#37](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/37)) ([aea00de](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/aea00de158cc321b87cb6954351268fe35e913dc))

## [2.2.1](https://github.com/VitruvianSoftware/aspect-workflows-template/compare/v2.2.0...v2.2.1) (2026-06-13)


### Bug Fixes

* **template:** make every preset render a fixed point of `bazel run //:tidy` ([#33](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/33)) ([c2a647a](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/c2a647a68f7c92686945f7a1148b71d6a410a6b4))

## [2.2.0](https://github.com/VitruvianSoftware/aspect-workflows-template/compare/v2.1.0...v2.2.0) (2026-06-12)


### Features

* **bazel:** upgrade Bazel 9.0.0 -&gt; 9.1.1 ([#31](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/31)) ([e7a687a](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/e7a687aab59ec31425e8366eafec5411b7804ace))

## [2.1.0](https://github.com/VitruvianSoftware/aspect-workflows-template/compare/v2.0.0...v2.1.0) (2026-06-07)


### Features

* add .vscode ([#554](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/554)) ([7d54145](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/7d541450af7e35bbbbb8a807f581118cd5d82fa6))
* Add Backstage template generation and scaffolding support ([df92462](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/df92462b7d050487964d228c29b60af8f7c46810))
* add clangd support via compile_commands.json ([#495](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/495)) ([5ed68e4](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/5ed68e4e454fa35cd7da893acfbc03cebf94cc10))
* add devcontainer ([5b72316](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/5b7231618fa2a81942b365afc3b2f62247146cbc))
* add Kotlin ([#613](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/613)) ([ef47ffc](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/ef47ffc8bf46dc2ce06414cd84e1299ff62658d4))
* add Pulumi infrastructure for repository management ([#9](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/9)) ([d3f7351](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/d3f7351b621bf38eb1772526d5ddbaaedbe3219d))
* add Ruby ([#815](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/815)) ([6273da9](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/6273da919e909f664858507f101bf73d9f789dcb))
* add Scala support ([#831](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/831)) ([2cc4540](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/2cc4540fbbf83cf761ba222648b5ef721e1ea0a7))
* add support for AXL ([#825](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/825)) ([afef839](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/afef8397cda7b371e2defd7c9a9d0777685f8d8f))
* adopt bazelrc-preset.bzl ([#530](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/530)) ([5fd5721](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/5fd57212d8501c4d39dde1b21505a2959f656efc))
* BUILD generation for kotlin ([#623](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/623)) ([5cd0cd3](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/5cd0cd35960793561c4dc1bcf1917ef12b48a45c))
* dormant in-repo RBE (remote build) setup helper ([#14](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/14)) ([136ea74](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/136ea74c957e90a8d9d173ecb930ac58ae241075))
* **iac:** add swift + backstage-swift starter repos ([04e697f](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/04e697f395eca1af6b5645f1183b3a98cbb9179e))
* **iac:** mark starter repos as GitHub template repositories ([f7e55b8](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/f7e55b800725a0bb4ee1838333fd91577017db94))
* install scaffold for Go codegen ([#587](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/587)) ([63c91ec](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/63c91eccdeb4e5c631cb156225564f228da3aac1))
* lint using AXL ([#682](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/682)) ([f25ab50](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/f25ab50a9645c62dda12bbe306709a96ecf57336))
* **module:** add rules_swift toolchain gated on .Computed.swift ([d79dc6d](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/d79dc6de005925378b028d34a410304acc553e59))
* opt-in build-cache menu (local / shared bazel-remote / buildbuddy) ([#27](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/27)) ([ffa26ef](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/ffa26ef2a948dc647a416a3ebe09ae5447ff9bb2))
* opt-in license-enforcement + Copybara-sync template features ([#13](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/13)) ([d0f2ebf](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/d0f2ebf383aa723b0db54463f19b852b7a67d3a3))
* populate README for each repo ([#628](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/628)) ([c21d773](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/c21d773f1b807aa4263ca546b9a1b7864c43e5f2))
* **pulumi:** also protect platform-v2.0 on the template repo ([df6e6cc](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/df6e6cc00fb98b46496a881ebff0b7928c91a74e))
* **pulumi:** opt-in Pulumi-in-CI/CD automation (preview/apply + central preview) ([#19](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/19)) ([3955ee3](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/3955ee3b0af0770adb2ff572d1a6baaa9873990c))
* **py:** run uv lock ([f76e8f2](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/f76e8f26f0fbadfc45b78ca024d194fae10fec50))
* re-enable orion gazelle extension ([#685](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/685)) ([a788431](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/a788431d0185eb134fbd4c49ed3ecd87486c0731))
* **repo-config:** GitHub repo settings as code (Pulumi) + Bazel-wrapped Pulumi ([#17](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/17)) ([22451fc](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/22451fccd7207c445ce9b806653d752f021a0448))
* ruby/proto/gRPC ([#868](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/868)) ([7105044](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/71050446a4fd179b341b1616ae8d3c8073c943c6))
* **ruby:** wire rubocop ([#821](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/821)) ([e279d74](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/e279d7467f0422bca02d89644201b56b04db23c8))
* **scaffold:** add Swift language + swift/backstage-swift presets ([99b2bb9](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/99b2bb9386875a77ff2d7e25bccafcf4429af9af))
* **swift:** hermetic Swift toolchain via rules_swift swift extension ([b033a72](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/b033a724757eec17f4992382a5653735a7048324))
* **swift:** load swift_binary from //swift:defs.bzl in backstage-swift story ([6b67018](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/6b670187c394030aceedd926c2f064588f3e5372))
* **tidy:** //:tidy one-command BUILD hygiene + Tidy Check gate ([#24](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/24)) ([ced3abb](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/ced3abba5d5fec4e85d09d6e5ac2489bf3b921ce))
* **tools:** swift-format for Swift via rules_lint ([d086a10](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/d086a10706708ccd62ddc9f60a43e777cd436b37))
* wire ty for Python typechecking ([#781](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/781)) ([652ec0f](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/652ec0ff337cb32e64c1961f5b69a55ef8575354))


### Bug Fixes

* add backstage feature gate for template.yaml and skeleton files ([845aa91](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/845aa919f74d698f6882554746c5fddcd81bbf40))
* add backstage-ruby and backstage-scala presets to scaffold.yaml, add fail-fast: false ([c9cff7d](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/c9cff7dd0d290df4b227c7e258bdcc3db701d7de))
* add ktlint ([#670](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/670)) ([2fc91d6](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/2fc91d68893c670a112d356dc86c289363340f7a))
* add ruby and scala to catalog-info.yaml template branches ([b0c4b36](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/b0c4b36c43935614c91154996c817b5cebd2c905))
* always place the java*.bazelrc file ([31dc58e](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/31dc58e4b8baa6390704b4e84ae841994bac01c1))
* always set python bootstrap script ([15082d8](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/15082d85ae05aaaa82e084fa64dcedc2e323fb7d))
* avoid rebuilding the formatter when possible ([#760](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/760)) ([a5defb1](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/a5defb14e523f757517d52167e291e34c55f5f98))
* axl_dep name ([fad0cf3](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/fad0cf38002ea04b6714de617a10efee9b2ca534))
* bad syntax in bash comment ([c7fd597](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/c7fd597772e64263dab01d1b8116905805c1903f))
* **bazelrc:** don't try to specify all linters on one line ([#799](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/799)) ([b8b5c44](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/b8b5c44c3bc3a7c2a5c9a531c59f56eb8c37d094)), closes [#774](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/774)
* **bazelversion:** version number must be on line 1 ([#746](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/746)) ([faba304](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/faba304bda0943207963d81b57f362118f09c510))
* **ci:** don't try to format the pre-commit hook ([0937601](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/0937601b84f281c6afe73b8fc35f81d0eb3f6b6a))
* **ci:** remove corrupt control character in ci.yaml ([99dd102](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/99dd102ba820e61c1b3d40785409219fe180ea32))
* **ci:** remove corrupt control character in deliver.yaml ([5fa5c50](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/5fa5c5005cdfcdf4f576e1f54e3ee2f0fd7e73de))
* **devcontainer:** pin base to 24.04 and make tool downloads arch-aware ([#12](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/12)) ([9f46861](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/9f4686121e4c553c284f89f198f50573d5d9e4a2))
* **devcontainer:** troubleshoot loading issues ([81215b1](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/81215b125874ce8b934fa9fb6c4f9b016f6940bb))
* **devcontainer:** troubleshoot loading issues ([81215b1](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/81215b125874ce8b934fa9fb6c4f9b016f6940bb))
* don't create eslint config when lint not selected ([e3b143b](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/e3b143bae62ad7f4bebe0678d49cebbf7fec94db))
* encode ftp.gnu.org workaround ([5dccd14](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/5dccd148916c07a7c29add94c9419ce27c29c83c))
* **init task:** the preset flag isn't wired ([#838](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/838)) ([a934d0b](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/a934d0b87242e1ac69df99a54d404c4128b17d0e))
* js_binary needs to see type:module ([#876](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/876)) ([7ab262e](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/7ab262ea731523cf40f465c09be96dc34e7f24e5))
* **js:** add npm on the PATH ([#851](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/851)) ([2f87389](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/2f873891256f55d36ce738b1ca48942901790932))
* macos not darwin ([766e166](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/766e166a707926cc9281a23228ea59e4abff165d))
* multitool update gha ([#716](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/716)) ([b312cf1](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/b312cf1dbc3ac4734cd1e28a29dbdbec525ea827))
* multitool workflow in project is slightly different ([8daf8bb](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/8daf8bbdb5215fb2a147b876ab6bbf1bbe457160))
* python bootstrap with bash, not python ([2453c33](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/2453c333151409faa8af9fb1ec74c16004699735))
* **rbe:** default RBE exec image to rbe-ubuntu20-04 (glibc &gt;= 2.28) ([#15](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/15)) ([dd85115](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/dd85115d5939c88a08aa5f586bdd3cbb6422550c))
* **rbe:** default RBE executor image to rbe-ubuntu22-04 (modern git) ([#25](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/25)) ([882557e](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/882557e14f4041db8792b23a643f95364b055972))
* recommend js_library for linking ([#759](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/759)) ([b787b2d](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/b787b2de5e3bcce9311531ca882304e29ca599b5))
* remove corrupt control character in scaffold.yaml ([ff33ea6](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/ff33ea69ed0b3c31f11c0d57394560ee23de2c8d))
* remove remaining corrupt control characters from BUILD and pre-commit ([520c3ea](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/520c3eadb09b7fb9de7d9deb95d6eb8563f8fe18))
* remove sh_binary wrapper for prettier, use prettier_binary directly ([6f7bb8c](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/6f7bb8cc146b10787ca49eb863d175ca240cc508))
* **repo-config:** close CI-automation gaps found by vitruvian-core dogfood ([#20](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/20)) ([850fe4e](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/850fe4e17b93c2e5c3785bc5ffe572b6a03432b2))
* **repo-config:** harden Pulumi module from vitruvian-core port ([#18](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/18)) ([fed1adf](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/fed1adf6874e9efbed589ba383dbe2c76bf02762))
* restore rust clippy.toml in BUILD and fix unbalanced template end ([671f0c5](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/671f0c5c70a98a1a84e67fef354364051c63f366))
* **template:** bump darwin-aarch64 LLVM 15.0.7 -&gt; 17.0.6 for macOS linking ([fdea797](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/fdea797f518a15cb82c00b2cee1e7a20a191617b))
* unused load ([0dc5916](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/0dc5916666551f28c7a09c3bd77807346f23bb11))
* update past 'hit' ([32cc797](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/32cc797c2e13a541c60a2827327b7b9932638c07))
* use cargo from upstream_wrapper ([#625](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/625)) ([96834eb](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/96834ebd6e6e5619c36232212c22ec7196464db8))
* use released rules_lint ([74508c4](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/74508c4fc17166fd8de473c067386a54ec49c64a))
* vscode recommended extensions ([e137b82](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/e137b8267f0950a74dab10a33e9d41f6398f530f))
* wire clang-tidy ([#496](https://github.com/VitruvianSoftware/aspect-workflows-template/issues/496)) ([ed36073](https://github.com/VitruvianSoftware/aspect-workflows-template/commit/ed360737e9ea9d5bb41be215b4c8d5035672b631))
