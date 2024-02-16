# ZooQL
**Query Zookeeper via SQL.**

## Overview

ZooQL consists of two main components: a backend written in Golang and a frontend built with Next.js. The backend process connects to a Zookeeper cluster and reads all the paths from a specified base path, storing them in an in-memory SQLite database. It also creates a watch on the Zookeeper path and updates the SQLite database based on the events, ensuring the SQLite database is always up-to-date with the current Zookeeper data. The backend also exposes a POST API for executing queries on the SQLite database.

The frontend is a simple SQL editor that calls the backend API and displays the results in a table format with an option to export to CSV.

## Backend

Backend connects to a Zookeeper cluster (provided via the `--zookeepers` flag, defaulting to `localhost:2181`) and reads all the paths from the base path (provided via the `--base-path` flag, defaulting to `/`).

After starting, it reads all the paths recursively under the base path and stores them in an in-memory SQLite database. The table in SQLite is called `znodes` with two fields: `path` and `data`.

```sql
CREATE TABLE znodes (
    path TEXT PRIMARY KEY,
    data TEXT
)
```

It then creates a [watch](https://zookeeper.apache.org/doc/current/zookeeperProgrammers.html#ch_zkWatches) on the Zookeeper path and updates the SQLite database based on the events, ensuring the SQLite database is always up-to-date with the current Zookeeper data.

The backend also exposes a POST API at `/backend/query` that accepts a JSON body with a field called `query`.
For example:

```json
{ "query": "select count(*) from znodes;" }
```

**The backend accepts any query that is valid in SQLite.**

The query result is structured as follows:

```go
type queryResult struct {
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
	Elapsed string     `json:"elapsed"`
}
```

## Frontend

Located in the `/frontend` directory, the frontend is a Next.js application. It provides a basic SQL editor that calls the `/backend/query` API and displays the result as a table with an option to export the data to a CSV file.

## Getting Started

To get started with ZooQL, clone the repository and navigate to the respective directories to start the backend and frontend services.

**Starting the backend:**

```bash
cd backend
go run ./... --zookeepers=localhost:2183 --base-path=/
```

**Starting the frontend:**
```bash
cd frontend
npm run dev
```

## Contributing

Contributions to ZooQL are welcome!

## License

ZooQL is open-source software licensed under [MIT](LICENSE).
