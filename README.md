# cdcdb: CHANGE DATA CAPTURE DATABASE

`cdcdb`: example demo program the connects to a database and receive some message events(`insert`, `update`, `delete`, ...)

## Database:

- [PostgreSQL 13.2](https://www.postgresql.org/docs/13/release-13-2.html) : [Write-Ahead Logging (WAL)](https://www.postgresql.org/docs/13/wal-intro.html)
  - [x] `insert`
  - [ ] `update`
  - [ ] `delete`
- [MongoDB 4.4](https://docs.mongodb.com/manual/release-notes/4.4/): [Change Streams](https://docs.mongodb.com/manual/changeStreams/)
  - [ ] `insert`
  - [ ] `update`
  - [ ] `replace`
  - [ ] `delete`
- [MySQL 8.0](https://dev.mysql.com/doc/relnotes/mysql/8.0/en/): [Binary Log](https://dev.mysql.com/doc/internals/en/binary-log-overview.html)
  - [ ] `insert`
  - [ ] `update`
  - [ ] `delete`

## How to run

### `Postgres`

- Into the executable directory

  ```
  cd pg
  ```

- Setup replication Postgres DB inside docker evironment:

  ```
  make setup
  ```

- Create some tables to listen
- Start listening:

  ```
  make dev
  ```
