### Commands 

### Install Soda CLI
`soda cli` helps to manage database connections\
URL [Soda](https://gobuffalo.io/documentation/database/soda/)

Install the soda cli
` go install github.com/gobuffalo/pop/v6/soda@latest`

### Database Configuration
Pop configurations managed by `database.yml` file located at the root of your project.

#### Sample config for a project base on postgresql
```yaml
development:
  dialect: postgres
  database: myapp_development
  user: postgres
  password: postgres
  host: 127.0.0.1
  pool: 5

test:
  url: {{envOr "TEST_DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/myapp_test"}}

production:
  url: {{envOr "DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/myapp_production"}}
```
### Using fizz
Fizz is a common DSL for migrating databases\
Generate Migrations
`soda generate fizz name_of_migrations`

### Create a User Table
`soda generate fizz CreateUserTable`
- this will generate an `up` and `down` migration files.

### up and down migrations
To create new table, write your query on the `up` file.
The `down` is used to reverse migrations created using `up` command.

### Execute -up
`sodo migrate`\
Example output
```bigquery
soda migrate
pop v6.1.1

[POP] 2024/11/15 13:51:26 info - > create_user_table
[POP] 2024/11/15 13:51:26 info - Successfully applied 1 migrations.
[POP] 2024/11/15 13:51:26 info - 0.0348 seconds
[POP] 2024/11/15 13:51:26 info - dumped schema for bookings
```

### Using down migration
Reverse user table `creation` - drop users table `sql("drop table users")`\
`sql migrate down`
```bigquery
    soda migrate down
    pop v6.1.1
    
    [POP] 2024/11/15 13:58:05 warn - ignoring file schema.sql because it does not match the migration file pattern
    [POP] 2024/11/15 13:58:05 info - < create_user_table
    [POP] 2024/11/15 13:58:05 info - 0.0132 seconds
    [POP] 2024/11/15 13:58:06 info - dumped schema for bookings
```

### Set up foregin Keys
To generate foreign key for a table, 
`soda generate fizz CreateFKForExxampleTable`\
_syntax_
```bigquery
    add_foreign_key("reservations", "room_id", {"rooms": ["id"]}, {
    "on_delete": "cascade",
    "on_update": "cascade",
})
```
The run `soda migrate` to apply\

#### Drop foreign key
On the `down` file,
`drop_foreign_key("reservations", "reservations_rooms_id_fk", {})`

### Create Table Indeces
To create unique indeces on a table field.
` soda generate fizz CreateUniqueIndexForUsersTable`\
**add_index("table_name", "column", "indextype")** - in you want default index, use {} otherwise:
_syntax_
```bigquery
    add_index("users","email", {"unique": true})
```
#### Drop Index
`drop_index("users", "users_email_idx")`\
`soda migrate down`

Create Indeces on two columns\
`soda generate fizz CreateIndecesOnRoomRestrictions`
On your `up` migrations\
```bigquery
add_index("room_restrictions", ["start_date","end_date"], {})
add_index("room_restrictions","room_id", {})
add_index("room_restrictions","reservation_id", {})

```

Add foreignKey and indeces together
`soda generate fizz CreateFKAndIndecesToReservationsTable`

### Drop and Recreate Dbs
To drop and reset the tables, run the command:\
`soda reset`
