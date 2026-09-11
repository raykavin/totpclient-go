# totpclient-go

A simple command-line TOTP (Time-based One-Time Password) client written in Go with no external dependencies.

The application generates and displays a 6-digit TOTP code using a Base32 secret provided through an environment variable.

## Requirements

* Go 1.26 or later

## Usage

Clone the repository:

```bash
git clone https://github.com/raykavin/totpclient-go.git
cd totpclient-go
```

Set your TOTP secret:

```bash
export TOTP_SECRET="YOUR_TOTP_SECRET"
```

Run the application:

```bash
go run .
```

The generated TOTP code will be displayed in the terminal and updated continuously.

## Build

To build the binary:

```bash
go build -o totpclient
```

Then run:

```bash
./totpclient
```
