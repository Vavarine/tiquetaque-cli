# TiqueTaque CLI

**ttq-cli** is a command-line interface (CLI) tool written in Go to punch employee time using the TiqueTaque API.  
The CLI securely saves tokens and user info (keyring or file fallback), allowing multiple executions.

## Features

- Login via email and SMS code
- Punches using saved token
- Secure storage of token, employeeID, and employee name
- Optional date and time flags for manual punching
- Modular and testable structure (Client, App, CMD)


## Installation
To install ttq-cli, use the following command:

```bash
  curl -fsSL https://raw.githubusercontent.com/Vavarine/tiquetaque-cli/main/install.sh | bash
```

## Usage
After installation, you can use the `ttq` command in your terminal. 

- To log in: `ttq login --email <email> --code <sms_code>` or `ttq l --email <email> --code <sms_code>`
- To punch: `ttq punch` or `ttq p`
- To view history: `ttq history` or `ttq h`

Log in is required only once. The token and user info are securely stored for future use.

## Contributing
Contributions are welcome! Please fork the repository and create a pull request with your changes.