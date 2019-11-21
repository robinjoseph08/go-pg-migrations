
<a name="v2.0.1"></a>
## [v2.0.1](https://github.com/robinjoseph08/go-pg-migrations/compare/v2.0.0...v2.0.1) (2019-11-20)

### Bug Fixes

* **create:** update migration template to refer to /v2 ([#18](https://github.com/robinjoseph08/go-pg-migrations/issues/18))


<a name="v2.0.0"></a>
## [v2.0.0](https://github.com/robinjoseph08/go-pg-migrations/compare/v1.0.1...v2.0.0) (2019-10-26)

### Features

* **modules:** Add go.mod for Go modules ([#16](https://github.com/robinjoseph08/go-pg-migrations/issues/16))


<a name="v1.0.1"></a>
## [v1.0.1](https://github.com/robinjoseph08/go-pg-migrations/compare/v1.0.0...v1.0.1) (2019-10-26)

### Bug Fixes

* **lock:** Changed migrations.go to use the use_zero flag ([#15](https://github.com/robinjoseph08/go-pg-migrations/issues/15))


<a name="v1.0.0"></a>
## [v1.0.0](https://github.com/robinjoseph08/go-pg-migrations/compare/v0.1.2...v1.0.0) (2019-08-05)

### Features

* **files:** updated timestamp in filename to be UTC ([#13](https://github.com/robinjoseph08/go-pg-migrations/issues/13))


<a name="v0.1.2"></a>
## [v0.1.2](https://github.com/robinjoseph08/go-pg-migrations/compare/v0.1.1...v0.1.2) (2018-12-22)

### Code Refactoring

* **migrate:** acquire lock in a single statement ([#12](https://github.com/robinjoseph08/go-pg-migrations/issues/12))


<a name="v0.1.1"></a>
## [v0.1.1](https://github.com/robinjoseph08/go-pg-migrations/compare/v0.1.0...v0.1.1) (2018-08-23)

### Bug Fixes

* **setup:** Use `*pg.DB.CreateTable` for a consistent interface ([#10](https://github.com/robinjoseph08/go-pg-migrations/issues/10))

### Documentation

* **coveralls:** Send coverage info to Coveralls ([#8](https://github.com/robinjoseph08/go-pg-migrations/issues/8))

### Features

* **errors:** Add migration name for migrate/rollback errors ([#11](https://github.com/robinjoseph08/go-pg-migrations/issues/11))


<a name="v0.1.0"></a>
## v0.1.0 (2018-08-18)

### Documentation

* **changelog:** Add chglog support ([#7](https://github.com/robinjoseph08/go-pg-migrations/issues/7))
* **help:** Add help command and flesh out README.md ([#6](https://github.com/robinjoseph08/go-pg-migrations/issues/6))
* **readme:** Add a README.md base

### Features

* **base:** Add base for migrations package ([#1](https://github.com/robinjoseph08/go-pg-migrations/issues/1))
* **create:** Add create command ([#2](https://github.com/robinjoseph08/go-pg-migrations/issues/2))
* **migrate:** Add migrate command ([#4](https://github.com/robinjoseph08/go-pg-migrations/issues/4))
* **rollback:** Add rollback command ([#5](https://github.com/robinjoseph08/go-pg-migrations/issues/5))
* **setup:** Add migration tables and functions to set them up ([#3](https://github.com/robinjoseph08/go-pg-migrations/issues/3))

