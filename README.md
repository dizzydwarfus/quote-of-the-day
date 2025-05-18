# quote-of-the-day

Simple api to randomly generate quote of the day (100 different quotes)

## To-dos

- [x] create data array of fixed quotes
- [x] create simple GET HTTP endpoint for randomly retrieving quote from fixed array
- [x] use htmx to serve the quote on the page

## Usage

This project uses air to run the server. You can install it using the following command:

```bash
go install github.com/cosmtrek/air@latest
```

Then, you can run the server using the following command:

```bash
air
```

The server is proxied on `localhost:3000`, and serving on `localhost:8080`.
