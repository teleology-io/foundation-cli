# Foundation CLI
The Foundation CLI is your gateway to effortless interaction with the Foundation API. Designed with simplicity and efficiency in mind, this library abstracts away the complexity of direct API calls, providing a clean and intuitive interface for developers. 

## Installation

```
curl -s https://raw.githubusercontent.com/teleology-io/foundation-cli/master/download.sh | bash
```

## Authentication
The Foundation CLI can be authenticated in two ways:

1. The `--api-key` argument:
```bash
foundation --api-key='your_key_here' <command>
```

2. Or exposed in the current process:

```bash
export FOUNDATION_API_KEY='your_key_here'

foundation <command>
```

The second approach is the suggested method as it reduces the amount of typing and increases the readability of the commands. The following examples are assuming the environment variable approach.


### Environment Usage:
To display the environment for your project to stdout run:

```bash
foundation environment
```

To inject the environment into the currently running process run:

```bash
source <(foundation environment)
```

### Configuration Usage:
To display the configuration for your project to stdout run:

```bash
foundation configuration
```

### Variable Usage:
To make a request for the currently served variable:

```bash
foundation variable --name <name_of_variable>
```

To get the value served for a unique identifer run:

```bash
foundation variable --name <name_of_variable> --uid '<unique-id>'
```