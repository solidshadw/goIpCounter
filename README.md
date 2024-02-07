# goIpCounter

This Go script checks a list of IP addresses against a list of subnets to determine whether each IP address is found in any of the subnets.

You can also include subnets in the `ips.txt` file, but it will just match and print out the subnets in the `subnets.txt` file.

## Usage

1. Clone this repository:
    ```bash
    git clone <repository_url>
    ```

2. Navigate to the directory:
    ```bash
    cd ip-subnet-checker
    ```

3. Run the script with the following command:

    ```bash
    go run main.go -ci <ip_file> -s <subnet_file> [-f | -nf]
    ```

    - `<ip_file>`: Path to the file containing IP addresses/subnets.
    - `<subnet_file>`: Path to the file containing subnets.
    - `-f`: Flag to specify IPs found in subnets.
    - `-nf`: Flag to specify IPs not found in subnets.

    Note: Only one of the flags `-f` or `-nf` should be specified at a time.

## Example

Suppose you have the following files:

- `ips.txt:` Which contains ips/subnets to be checked
```
192.168.0.5
192.168.1.0/24
192.16.1.56
192.168.1.253
192.168.1.55
192.0.1.5
192.168.1.10
192.168.24.0/24
```
- `subnets.txt`: 
```
192.168.1.0/24
192.168.34.0/24
```

To check which IPs are found in the subnets, run the script as follows:

```bash
go run ./main.go -ci ips.txt -s subnets.txt -f
```
Output:
```bash
192.168.1.0/24
192.168.1.253
192.168.1.55
192.168.1.10
```
This indicates that `192.168.1.253 \ 192.168.1.55 \ 192.168.1.10 and 192.168.1.0/24` ARE found in the specified subnets file.