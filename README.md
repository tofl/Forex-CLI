# Forex CLI

## Installation
Clone this repository and install it with the Go compiler (make sure it is installed on your machine).

```bash
$ go install
```

Sign up for a free account on the [Alpha Vantage API website](https://www.alphavantage.co/support/#api-key). You will be given a key. Then run :
```bash
$ forex key [your key]
```

## Usage
```bash
$ forex [sum] [from] [to]
```

For example, the following command will return the value of `12USD` in `EUR` :
```bash
$ forex 12 USD EUR
```