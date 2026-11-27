![habit-tracker](https://user-images.githubusercontent.com/15822994/145180371-3d0fe074-ffb2-468f-a2c3-0fbdbcd75c83.png)

![open issues](https://img.shields.io/github/issues/bashbunni/habit-tracker)
![forks](https://img.shields.io/github/forks/bashbunni/habit-tracker)
![stars](https://img.shields.io/github/stars/bashbunni/habit-tracker)
![contributions](https://img.shields.io/badge/contributions-open-blueviolet)

## About

My project is a web and desktop application that allows you to track and visualize your consistency with your habits over time with a calendar heatmap.
This allows you to track any habits you’d like and see how you’re progressing towards your goals.
I hope that this tool can help those who struggle with being consistent or disciplined.

![habit page](https://user-images.githubusercontent.com/15822994/145179044-4c153162-075f-4842-9068-691660c6c2bd.png)
![log activity menu](https://user-images.githubusercontent.com/15822994/145179136-37aca325-ca19-42e0-83d9-2ea3e4498729.png)

## How to Run

### Dependencies

#### Frontend

- all frontend dependencies are included in the `package.json` file
- you can install these dependencies by navigating to the `frontend` folder then running `npm i`

#### Backend

- [Wails](https://wails.app/gettingstarted/) and its dependencies
  - Please follow their 'Getting Started' guide to install all project dependencies
- Go 1.17 or newer
- Dependencies listed in `go.mod`
  - install these dependencies by running `go mod download` in the backend/root directory (same level as main.go)
- Since auth isn't set up yet, you will need a `.env` file with the following values:

```
DB_USERNAME=insertusername
DB_PASSWORD=insertpassword
HOST=inserthost
PORT=insertport
DATABASE=insertdatabase
```

**_note_: You must use a MySQL db for this project**\
For example, if you're running a local MySQL database your .env file might look something like this:

```
DB_USERNAME=root
DB_PASSWORD=root
HOST=localhost
PORT=3306
DATABASE=HABIT_TRACKER
```

Create the tables with the schema in `schema.sql`:

```
mysql -u root -p HABIT_TRACKER < schema.sql
```

If you already have a database from an earlier version, run the migration statements at the bottom of `schema.sql` once instead (back up first). They merge duplicate entries per habit and day and add the unique key the app now relies on.

### Running the App

- In a terminal window, in the root folder (same folder as `main.go`) run `wails serve` to start the backend.
- In a new terminal window, navigate to the frontend folder, then run `npm run serve` to start the frontend.
  This runs the development server for the desktop application on `http://localhost:3000`

## Running the Tests

### Backend

The repository tests in `habit_repository_test.go` and `date_repository_test.go` use [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock), so they need neither a MySQL database nor a `.env` file. From the root folder run:

```
go test ./...
```

To see which statements are covered, run `go test -cover ./...`.

The tests live in package `main`, which imports Wails, so the Wails system dependencies have to be installed to compile them (the same ones you need for `wails serve`). On macOS that is the Xcode command line tools. On Debian/Ubuntu it is `libgtk-3-dev` and `libwebkit2gtk-4.0-dev`.

What the tests cover:

- every method of `HabitRepository` and `DateRepository`, including the error paths and invalid JSON input
- user input is passed to MySQL as query parameters, never pasted into the SQL string
- `AddCount` adds to an existing day in a single upsert statement, and rejects entries without a habit, with a count of 0 or with a date that is not `YYYY-MM-DD`
- "today" in `GetTodaysCountForHabit` is the local calendar day, not the UTC day
- an empty result is returned as `[]`, not `null`, so the frontend always gets a list

When you add a repository method, add a test next to the others and keep the expected SQL in the constants at the top of the test file in sync.

### Frontend

From the `frontend` folder run `npm test`. This starts Jest in watch mode; use `CI=true npm test` for a single run.

Note: `src/setupTests.js` and `src/App.test.js` import `@testing-library/jest-dom` and `@testing-library/react`, which are not in `package.json` yet. Until they are added, install them once with `npm i -D @testing-library/react@12 @testing-library/jest-dom@5`. `App.test.js` is still the Create React App placeholder and is expected to fail.

## How to Contribute

1. Create an issue\*
2. Fork the repository
3. Create a solution that solves _exactly_ the features that are mentioned in the issue
4. Create a pull request with a link to the issue following the format:

```
## Completed
- [x] add timer functionality
- [x] create timer component
- [x] integrate timer React component with timer functionality

## Goal
Build out the pomodoro timer as outlined in feature request issue #10
```

- the list of tasks completed will likely match your commit messages if you're writing meaningful commit messages, if not, it will add clarity to the pull request.

\*issues should include a clear description. For bugs this means including screenshots of related error messages and code snippets that you suspect are involved along with your OS, Go, Wails, Node versions.
For feature requests, please include a description of the feature, with a user story outlining the problem-space it aims to solve.

### Additional Information
Here's an overview of the project structure:
<img width="1168" alt="habit-tracker (2)" src="https://user-images.githubusercontent.com/15822994/145290719-315ed951-8e6f-489d-ba60-b7a790e13fd1.png">


## Next Steps

### Add Pomodoro functionality

This includes

- a functional timer component
- incrementing the count for a habit that has pomodoro enabled

![pomodoro page](https://user-images.githubusercontent.com/15822994/145179376-d87d1354-7497-4cef-90d0-ac7690438457.png)

### Add Authentication

I would like for users to be able to keep track of their own habits and separate data by user for production.

### Release Production Build

Habit Tracker is currently in release for development only, meaning that a stable production release still needs to be developed.
I'm planning on migrating the project to Wails v2 which has a lot of bug fixes available.
Once I've finished the migration I'll be working on building the executable desktop programs for Mac, Linux, and Windows.
That will be the first major release
